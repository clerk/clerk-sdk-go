// Package jwt provides operations for decoding and validating
// JSON Web Tokens.
package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/go-jose/go-jose/v3/jwt"
)

// AuthorizedPartyHandler is a type that can be used to perform checks
// on the 'azp' claim.
type AuthorizedPartyHandler func(string) bool

// CustomClaimsConstructor can initialize structs for holding custom
// JWT claims.
type CustomClaimsConstructor func(context.Context) any

type VerifyParams struct {
	// Token is the JWT that will be verified. Required.
	Token string
	// JWK is the custom JSON Web Key that will be used to verify the
	// Token with.
	// If the JWK parameter is provided, the Verify method won't
	// fetch the JSON Web Key Set and there's no need to provide
	// the JWKSClient parameter.
	JWK *clerk.JSONWebKey
	// JWKSClient is a jwks API client that will be used to fetch the
	// JSON Web Key Set for verifying the Token with.
	// If the JWK parameter is provided, the JWKSClient is not needed.
	// If no JWK or JWKSClient is provided, the Verify method will use
	// a JWKSClient with the default Backend.
	JWKSClient *jwks.Client
	// Clock can be used to keep track of time and will replace usage of
	// the [time] package. Pass a custom Clock to control the source of
	// time or facilitate testing chronologically sensitive flows.
	Clock clerk.Clock
	// CustomClaimsConstructor will be called when parsing the Token's
	// claims. It's useful for parsing custom claims into user-defined
	// types.
	// Make sure it returns a pointer to a type (struct) that describes
	// any custom claims schema with the correct JSON tags.
	//	type MyCustomClaims struct {}
	//	VerifyParams{
	//		CustomClaimsConstructor: func(_ context.Context) any {
	//			return &MyCustomClaims{}
	//		},
	//	}
	CustomClaimsConstructor CustomClaimsConstructor
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
	parsedToken, err := jwt.ParseSigned(params.Token)
	if err != nil {
		return nil, err
	}
	if len(parsedToken.Headers) == 0 {
		return nil, fmt.Errorf("missing JWT headers")
	}
	jwk := params.JWK
	if jwk == nil {
		jwk, err = GetJSONWebKey(ctx, &GetJSONWebKeyParams{
			KeyID:      parsedToken.Headers[0].KeyID,
			JWKSClient: params.JWKSClient,
		})
		if err != nil {
			return nil, err
		}
	}
	if jwk == nil {
		return nil, fmt.Errorf("missing json web key, need to set JWK in the params")
	}

	if parsedToken.Headers[0].Algorithm != jwk.Algorithm {
		return nil, fmt.Errorf("invalid signing algorithm %s", jwk.Algorithm)
	}

	claims := &clerk.SessionClaims{}
	allClaims := []any{claims}
	if params.CustomClaimsConstructor != nil {
		claims.Custom = params.CustomClaimsConstructor(ctx)
		allClaims = append(allClaims, claims.Custom)
	}
	err = parsedToken.Claims(jwk.Key, allClaims...)
	if err != nil {
		return nil, err
	}

	clock := params.Clock
	if clock == nil {
		clock = clerk.NewClock()
	}
	err = claims.ValidateWithLeeway(clock.Now().UTC(), params.Leeway)
	if err != nil {
		return nil, err
	}

	// Non-satellite domains must validate the issuer.
	if !params.IsSatellite && !isValidIssuer(claims.Issuer, params.ProxyURL) {
		return nil, fmt.Errorf("invalid issuer %s", claims.Issuer)
	}

	if params.AuthorizedPartyHandler != nil && !params.AuthorizedPartyHandler(claims.AuthorizedParty) {
		return nil, fmt.Errorf("invalid authorized party %s", claims.AuthorizedParty)
	}

	return claims, nil
}

func isValidIssuer(iss string, proxyURL *string) bool {
	if proxyURL != nil {
		return iss == *proxyURL
	}
	return strings.HasPrefix(iss, "https://clerk.") ||
		strings.Contains(iss, ".clerk.accounts")
}

type DecodeParams struct {
	Token string
}

// Decode decodes a JWT without verifying it.
// WARNING: The token is not validated, therefore the returned Claims
// should NOT be trusted.
func Decode(_ context.Context, params *DecodeParams) (*clerk.UnverifiedToken, error) {
	parsedToken, err := jwt.ParseSigned(params.Token)
	if err != nil {
		return nil, err
	}

	standardClaims := clerk.RegisteredClaims{}
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

	claims := &clerk.UnverifiedToken{
		RegisteredClaims: standardClaims,
		Extra:            extraClaims,
	}
	if len(parsedToken.Headers) > 0 {
		claims.KeyID = parsedToken.Headers[0].KeyID
	}
	return claims, nil
}

type GetJSONWebKeyParams struct {
	// KeyID is the token's 'kid' claim.
	KeyID string
	// JWKSClient can be used to call the jwks Get Clerk API operation.
	JWKSClient *jwks.Client
}

// GetJSONWebKey fetches the JSON Web Key Set from the Clerk API
// and returns the JSON Web Key corresponding to the provided KeyID.
// A default client will be initialized if the provided JWKSClient
// is nil.
func GetJSONWebKey(ctx context.Context, params *GetJSONWebKeyParams) (*clerk.JSONWebKey, error) {
	if params.KeyID == "" {
		return nil, fmt.Errorf("missing jwt kid header claim")
	}

	jwksClient := params.JWKSClient
	if jwksClient == nil {
		jwksClient = &jwks.Client{
			Backend: clerk.GetBackend(),
		}
	}
	jwks, err := jwksClient.Get(ctx, &jwks.GetParams{})
	if err != nil {
		return nil, err
	}
	if jwks == nil {
		return nil, fmt.Errorf("no jwks found")
	}

	for _, k := range jwks.Keys {
		if k != nil && k.KeyID == params.KeyID {
			return k, nil
		}
	}
	return nil, fmt.Errorf("missing json web key")
}
