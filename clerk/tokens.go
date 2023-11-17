package clerk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

var standardClaimsKeys = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}

type TokenClaims struct {
	jwt.Claims
	Extra map[string]interface{}
}

type SessionClaims struct {
	jwt.Claims
	SessionID                     string          `json:"sid"`
	AuthorizedParty               string          `json:"azp"`
	ActiveOrganizationID          string          `json:"org_id"`
	ActiveOrganizationSlug        string          `json:"org_slug"`
	ActiveOrganizationRole        string          `json:"org_role"`
	ActiveOrganizationPermissions []string        `json:"org_permissions"`
	Actor                         json.RawMessage `json:"act,omitempty"`
}

// DecodeToken decodes a jwt token without verifying it.
func (c *client) DecodeToken(token string) (*TokenClaims, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	standardClaims := jwt.Claims{}
	extraClaims := make(map[string]interface{})

	if err = parsedToken.UnsafeClaimsWithoutVerification(&standardClaims, &extraClaims); err != nil {
		return nil, err
	}

	// Delete any standard claims included in the extra claims
	for _, key := range standardClaimsKeys {
		delete(extraClaims, key)
	}

	return &TokenClaims{Claims: standardClaims, Extra: extraClaims}, nil
}

type verifyTokenOptions struct {
	authorizedParties map[string]struct{}
	leeway            time.Duration
	jwk               *jose.JSONWebKey
	customClaims      interface{}
	isSatellite       bool
	proxyURL          string
}

// VerifyToken verifies the session jwt token.
func (c *client) VerifyToken(token string, opts ...VerifyTokenOption) (*SessionClaims, error) {
	options := &verifyTokenOptions{}

	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}

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

	jwk := options.jwk
	if jwk == nil {
		jwk, err = c.getJWK(kid)
		if err != nil {
			return nil, err
		}
	}

	if parsedToken.Headers[0].Algorithm != jwk.Algorithm {
		return nil, fmt.Errorf("invalid signing algorithm %s", jwk.Algorithm)
	}

	claims := SessionClaims{}
	if err = verifyTokenParseClaims(parsedToken, jwk.Key, &claims, options); err != nil {
		return nil, err
	}

	if err = claims.Claims.ValidateWithLeeway(jwt.Expected{Time: time.Now()}, options.leeway); err != nil {
		return nil, err
	}

	issuer := newIssuer(claims.Issuer).
		WithSatelliteDomain(options.isSatellite).
		WithProxyURL(options.proxyURL)

	if !issuer.IsValid() {
		return nil, fmt.Errorf("invalid issuer %s", claims.Issuer)
	}

	if claims.AuthorizedParty != "" && len(options.authorizedParties) > 0 {
		if _, ok := options.authorizedParties[claims.AuthorizedParty]; !ok {
			return nil, fmt.Errorf("invalid authorized party %s", claims.AuthorizedParty)
		}
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

func verifyTokenParseClaims(parsedToken *jwt.JSONWebToken, key interface{}, sessionClaims *SessionClaims, options *verifyTokenOptions) error {
	if options.customClaims == nil {
		return parsedToken.Claims(key, sessionClaims)
	}
	return parsedToken.Claims(key, sessionClaims, options.customClaims)
}
