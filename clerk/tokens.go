package clerk

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"

	"gopkg.in/square/go-jose.v2/jwt"
)

type TokensService service

type jwk struct {
	// Sig (for signature) or Enc (for encryption)
	PublicKeyUse string `json:"use"`

	// Algorithm family (RSA, ECDSA etc.)
	KeyType string `json:"kty"`

	// RSA256
	Algorithm string `json:"alg"`

	// Clerk instance ID
	KeyID string `json:"kid"`

	Modulus  string `json:"n"`
	Exponent string `json:"e"`
}

// Verify the short-lived jwt token.
func (s *TokensService) Verify(token string) error {
	issuer, err := getIssuer(token)
	if err != nil {
		return err
	}

	//TODO: refresh periodically in the background
	jwks, err := s.fetchJWKs(issuer)
	if err != nil {
		panic(err)
		return err
	}

	pubKey, err := createRSAPublicKey(jwks[0])
	if err != nil {
		return err
	}

	err = verifyToken(token, jwks[0].Algorithm, pubKey)
	if err != nil {
		return err
	}

	return nil
}

// Decode the short-lived jwt token.
func (s *TokensService) Decode(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	claims := make(map[string]interface{})
	if err = parsedToken.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func getIssuer(token string) (string, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return "", err
	}

	claims := make(map[string]interface{})
	if err = parsedToken.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return "", err
	}

	issuer, ok := claims["iss"]
	if !ok {
		return "", fmt.Errorf("issuer not present in claims")
	}

	return issuer.(string), nil
}

func verifyToken(token, signingAlg string, pubKey *rsa.PublicKey) error {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return err
	}

	if len(parsedToken.Headers) == 0 {
		return fmt.Errorf("no JWT headers found")
	}

	if parsedToken.Headers[0].Algorithm != signingAlg {
		return fmt.Errorf("unexpected signing algorithm")
	}

	claims := make(map[string]interface{})
	if err = parsedToken.Claims(pubKey, &claims); err != nil {
		return err
	}

	for key, value := range claims {
		fmt.Println(key, " -> ", value)
	}

	return nil
}

func createRSAPublicKey(jwk *jwk) (*rsa.PublicKey, error) {
	if jwk == nil {
		return nil, fmt.Errorf("jwk cannot be nil")
	}

	if jwk.Exponent == "" || jwk.Modulus == "" {
		return nil, fmt.Errorf("jwk exponent or modulus cannot be empty")
	}

	exponent, err := base64.RawURLEncoding.DecodeString(jwk.Exponent)
	if err != nil {
		return nil, err
	}

	modulus, err := base64.RawURLEncoding.DecodeString(jwk.Modulus)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(modulus),
		E: int(big.NewInt(0).SetBytes(exponent).Uint64()),
	}, nil
}

func (s *TokensService) fetchJWKs(issuer string) ([]*jwk, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/.well-known/jwks.json", issuer), nil)
	if err != nil {
		return nil, err
	}

	jwks := struct {
		Keys []*jwk `json:"keys"`
	}{}

	_, err = s.client.Do(req, &jwks)
	if err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		panic("invalid length")
	}

	return jwks.Keys, nil
}
