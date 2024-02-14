// Package jwt provides operations for decoding and validating
// JSON Web Tokens.
package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-jose/go-jose/v3/jwt"
)

type AuthorizedPartyHandler func(string) bool

type VerifyParams struct {
	// Token is the JWT that will be verified. Required.
	Token string
	// JWK the custom JSON Web Key that will be used to verify the
	// Token with. Required.
	JWK          *clerk.JSONWebKey
	CustomClaims any
	// Leeway is the duration which the JWT is considered valid after
	// it's expired. Useful for defending against server clock skews.
	Leeway time.Duration
	// IsSatellite signifies that the JWT is verified on a satellite domain.
	IsSatellite bool
	// ProxyURL is the URL of the server that proxies the Clerk Frontend API.
	ProxyURL *string
	// AuthorizedPartyHandler can be used to perform validations on the
	// 'azp' claim.
	AuthorizedPartyHandler AuthorizedPartyHandler
}

// Verify verifies a Clerk session JWT and returns the parsed
// clerk.SessionClaims.
func Verify(ctx context.Context, params *VerifyParams) (*clerk.SessionClaims, error) {
	jwk := params.JWK
	if jwk == nil {
		return nil, fmt.Errorf("missing json web key, need to set JWK in the params")
	}

	parsedToken, err := jwt.ParseSigned(params.Token)
	if err != nil {
		return nil, err
	}
	if parsedToken.Headers[0].Algorithm != jwk.Algorithm {
		return nil, fmt.Errorf("invalid signing algorithm %s", jwk.Algorithm)
	}

	claims := &clerk.SessionClaims{}
	allClaims := []any{claims}
	if params.CustomClaims != nil {
		allClaims = append(allClaims, params.CustomClaims)
	}
	err = parsedToken.Claims(jwk.Key, allClaims...)
	if err != nil {
		return nil, err
	}

	err = claims.Claims.ValidateWithLeeway(jwt.Expected{Time: time.Now().UTC()}, params.Leeway)
	if err != nil {
		return nil, err
	}

	iss := claims.Issuer
	if params.ProxyURL != nil && *params.ProxyURL != "" {
		iss = *params.ProxyURL
	}
	// Non-satellite domains must validate the issuer.
	if !params.IsSatellite && !isValidIssuer(iss) {
		return nil, fmt.Errorf("invalid issuer %s", iss)
	}

	if params.AuthorizedPartyHandler != nil && !params.AuthorizedPartyHandler(claims.AuthorizedParty) {
		return nil, fmt.Errorf("invalid authorized party %s", claims.AuthorizedParty)
	}

	return claims, nil
}

func isValidIssuer(iss string) bool {
	return strings.HasPrefix(iss, "https://clerk.") ||
		strings.Contains(iss, ".clerk.accounts")
}

type DecodeParams struct {
	Token string
}

// Decode decodes a JWT without verifying it.
// WARNING: The token is not validated, therefore the returned Claims
// should NOT be trusted.
func Decode(_ context.Context, params *DecodeParams) (*clerk.Claims, error) {
	parsedToken, err := jwt.ParseSigned(params.Token)
	if err != nil {
		return nil, err
	}

	standardClaims := jwt.Claims{}
	extraClaims := make(map[string]any)
	err = parsedToken.UnsafeClaimsWithoutVerification(&standardClaims, &extraClaims)
	if err != nil {
		return nil, err
	}

	// Delete any standard claims included in the extra claims.
	standardClaimsKeys := []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}
	for _, key := range standardClaimsKeys {
		delete(extraClaims, key)
	}

	return &clerk.Claims{
		Claims: standardClaims,
		Extra:  extraClaims,
	}, nil
}
