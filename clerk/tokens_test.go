package clerk

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2"

	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	dummyTokenClaims = map[string]interface{}{
		"iss":     "issuer",
		"sub":     "subject",
		"aud":     []string{"clerk"},
		"name":    "name",
		"picture": "picture",
	}

	dummyTokenClaimsExpected = TokenClaims{
		Claims: jwt.Claims{
			Issuer:   "issuer",
			Subject:  "subject",
			Audience: jwt.Audience{"clerk"},
			Expiry:   nil,
			IssuedAt: nil,
		},
		Extra: map[string]interface{}{
			"name":    "name",
			"picture": "picture",
		},
	}

	dummySessionClaims = SessionClaims{
		Claims: jwt.Claims{
			Issuer:   "https://clerk.issuer",
			Subject:  "subject",
			Audience: nil,
			Expiry:   nil,
			IssuedAt: nil,
		},
		SessionID:       "session_id",
		AuthorizedParty: "authorized_party",
	}
)

func TestClient_DecodeToken_EmptyToken(t *testing.T) {
	c, _ := NewClient("token")

	_, err := c.DecodeToken("")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_DecodeToken_Success(t *testing.T) {
	c, _ := NewClient("token")
	token, _ := testGenerateTokenJWT(t, dummyTokenClaims, "kid")

	got, _ := c.DecodeToken(token)

	if !reflect.DeepEqual(got, &dummyTokenClaimsExpected) {
		t.Errorf("Expected %+v, but got %+v", &dummyTokenClaimsExpected, got)
	}
}

func TestClient_VerifyToken_EmptyToken(t *testing.T) {
	c, _ := NewClient("token")

	_, err := c.VerifyToken("")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_MissingKID(t *testing.T) {
	c, _ := NewClient("token")
	token, _ := testGenerateTokenJWT(t, dummySessionClaims, "")

	_, err := c.VerifyToken(token)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_MismatchKID(t *testing.T) {
	c, _ := NewClient("token")
	token, pubKey := testGenerateTokenJWT(t, dummySessionClaims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "invalid-kid"))

	_, err := c.VerifyToken(token)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_MismatchAlgorithm(t *testing.T) {
	c, _ := NewClient("token")
	token, pubKey := testGenerateTokenJWT(t, dummySessionClaims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS512, "kid"))

	_, err := c.VerifyToken(token)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_InvalidKey(t *testing.T) {
	c, _ := NewClient("token")
	token, _ := testGenerateTokenJWT(t, dummySessionClaims, "kid")
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, privKey.Public(), jose.RS256, "kid"))

	_, err := c.VerifyToken(token)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_InvalidIssuer(t *testing.T) {
	c, _ := NewClient("token")

	claims := dummySessionClaims
	claims.Issuer = "issuer"

	token, pubKey := testGenerateTokenJWT(t, claims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	_, err := c.VerifyToken(token)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_ExpiredToken(t *testing.T) {
	c, _ := NewClient("token")

	expiredClaims := dummySessionClaims
	expiredClaims.Expiry = jwt.NewNumericDate(time.Now().Add(time.Second * -1))
	token, pubKey := testGenerateTokenJWT(t, expiredClaims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	_, err := c.VerifyToken(token)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_InvalidAuthorizedParty(t *testing.T) {
	c, _ := NewClient("token")

	claims := dummySessionClaims
	claims.AuthorizedParty = "fake-party"

	token, pubKey := testGenerateTokenJWT(t, claims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	_, err := c.VerifyToken(token, WithAuthorizedParty("authorized_party"))
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestClient_VerifyToken_Success(t *testing.T) {
	c, _ := NewClient("token")
	token, pubKey := testGenerateTokenJWT(t, dummySessionClaims, "kid")

	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	got, _ := c.VerifyToken(token)
	if !reflect.DeepEqual(got, &dummySessionClaims) {
		t.Errorf("Expected %+v, but got %+v", dummySessionClaims, got)
	}
}

func TestClient_VerifyToken_Success_ExpiredCache(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	token, pubKey := testGenerateTokenJWT(t, dummySessionClaims, "kid")

	mux.HandleFunc("/jwks", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		_ = json.NewEncoder(w).Encode(testBuildJWKS(t, pubKey, jose.RS256, "kid"))
	})

	client := c.(*client)
	client.jwksCache.expiresAt = time.Now().Add(time.Second * -5)

	got, _ := c.VerifyToken(token)
	if !reflect.DeepEqual(got, &dummySessionClaims) {
		t.Errorf("Expected %+v, but got %+v", dummySessionClaims, got)
	}
}

func TestClient_VerifyToken_Success_AuthorizedParty(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	token, pubKey := testGenerateTokenJWT(t, dummySessionClaims, "kid")

	mux.HandleFunc("/jwks", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		_ = json.NewEncoder(w).Encode(testBuildJWKS(t, pubKey, jose.RS256, "kid"))
	})

	client := c.(*client)
	client.jwksCache.expiresAt = time.Now().Add(time.Second * -5)

	got, _ := c.VerifyToken(token, WithAuthorizedParty("authorized_party"))
	if !reflect.DeepEqual(got, &dummySessionClaims) {
		t.Errorf("Expected %+v, but got %+v", dummySessionClaims, got)
	}
}

func testGenerateTokenJWT(t *testing.T, claims interface{}, kid string) (string, crypto.PublicKey) {
	t.Helper()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}

	signerOpts := &jose.SignerOptions{}
	signerOpts.WithType("JWT")
	if kid != "" {
		signerOpts.WithHeader("kid", kid)
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privKey}, signerOpts)
	if err != nil {
		t.Error(err)
	}

	builder := jwt.Signed(signer)
	builder = builder.Claims(claims)

	token, err := builder.CompactSerialize()
	if err != nil {
		t.Error(err)
	}

	return token, privKey.Public()
}

func testBuildJWKS(t *testing.T, pubKey crypto.PublicKey, alg jose.SignatureAlgorithm, kid string) *JWKS {
	t.Helper()

	return &JWKS{Keys: []jose.JSONWebKey{
		{
			Key:       pubKey,
			KeyID:     kid,
			Algorithm: string(alg),
			Use:       "sig",
		},
	}}
}
