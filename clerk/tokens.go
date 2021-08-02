package clerk

import (
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"

	"gopkg.in/square/go-jose.v2/jwt"
)

type Claims struct {
	jwt.Claims
	UserClaims
}

type UserClaims struct {
	Name                string `json:"name"`
	Picture             string `json:"picture"`
	UpdatedAt           int64  `json:"updated_at"`
	GivenName           string `json:"given_name,omitempty"`
	FamilyName          string `json:"family_name,omitempty"`
	PreferredUsername   string `json:"preferred_username,omitempty"`
	Gender              string `json:"gender,omitempty"`
	Birthdate           string `json:"birthdate,omitempty"`
	Email               string `json:"email,omitempty"`
	EmailVerified       bool   `json:"email_verified,omitempty"`
	PhoneNumber         string `json:"phone_number,omitempty"`
	PhoneNumberVerified bool   `json:"phone_number_verified,omitempty"`
}

// DecodeToken decodes the session jwt token without verifying it.
func (c *client) DecodeToken(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	claims := Claims{}
	if err = parsedToken.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// VerifyToken verifies the session jwt token.
func (c *client) VerifyToken(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	if len(parsedToken.Headers) == 0 {
		return nil, fmt.Errorf("missing jwt headers")
	}

	kid := parsedToken.Headers[0].KeyID
	if kid == "" {
		return nil, fmt.Errorf("missing jwt kid header claim")
	}

	jwk, err := c.getJWK(kid)
	if err != nil {
		return nil, err
	}

	if parsedToken.Headers[0].Algorithm != jwk.Algorithm {
		return nil, fmt.Errorf("invalid signing algorithm %s", jwk.Algorithm)
	}

	claims := Claims{}
	if err = parsedToken.Claims(jwk.Key, &claims); err != nil {
		return nil, err
	}

	expectedClaims := jwt.Expected{
		Audience: jwt.Audience{"clerk"},
		Time:     time.Now(),
	}

	if err = claims.Claims.ValidateWithLeeway(expectedClaims, 0); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (c *client) getJWK(kid string) (*jose.JSONWebKey, error) {
	if c.jwksCache.isInvalid() {
		jwks, err := c.jwks.ListAll()
		if err != nil {
			return nil, err
		}

		c.jwksCache.set(jwks)
	}

	return c.jwksCache.get(kid)
}
