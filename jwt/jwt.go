// Package jwt provides operations for decoding and validating
// JSON Web Tokens.
package jwt

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/go-jose/go-jose/v3/jwt"
)

type VerifyParams struct {
	// Token is the JWT that will be verified.
	Token string
	// JWK is a custom JSON Web Key that can be provided to skip
	// fetching one.
	JWK          *clerk.JSONWebKey
	CustomClaims any
	// Leeway is the duration which the JWT is considered valid after
	// it's expired. Useful for defending against server clock skews.
	Leeway time.Duration
	// IsSatellite signifies that the JWT is verified on a satellite domain.
	IsSatellite bool
	// ProxyURL is the URL of the server that proxies the Clerk Frontend API.
	ProxyURL *string
	// List of values that should match the azp claim.
	// Use SetAuthorizedParties to set the value.
	authorizedParties map[string]struct{}
}

// SetAuthorizedParties accepts a list of authorized parties to be
// set on the params.
func (params *VerifyParams) SetAuthorizedParties(parties ...string) {
	azp := make(map[string]struct{})
	for _, p := range parties {
		azp[p] = struct{}{}
	}
	params.authorizedParties = azp
}

// Verify verifies a Clerk session JWT and returns the parsed
// clerk.SessionClaims.
func Verify(ctx context.Context, params *VerifyParams) (*clerk.SessionClaims, error) {
	parsedToken, err := jwt.ParseSigned(params.Token)
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

	jwk := params.JWK
	if jwk == nil {
		jwk, err = getJWK(ctx, kid)
		if err != nil {
			return nil, fmt.Errorf("get jwk: %w", err)
		}
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

	if claims.AuthorizedParty != "" && len(params.authorizedParties) > 0 {
		if _, ok := params.authorizedParties[claims.AuthorizedParty]; !ok {
			return nil, fmt.Errorf("invalid authorized party %s", claims.AuthorizedParty)
		}
	}

	return claims, nil
}

// Retrieve the JSON web key for the provided id from the set.
func getJWK(ctx context.Context, kid string) (*clerk.JSONWebKey, error) {
	jwks, err := getJWKSWithCache(ctx)
	if err != nil {
		return nil, err
	}
	for _, k := range jwks.Keys {
		if k.KeyID == kid {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("no jwk key found for kid %s", kid)
}

// Returns the JSON web key set. Tries a cached value first, but if
// there's no value or the entry has expired, it will fetch the set
// from the API and cache the value.
func getJWKSWithCache(ctx context.Context) (*clerk.JSONWebKeySet, error) {
	const cacheKey = "/v1/jwks"
	var jwks *clerk.JSONWebKeySet
	var err error

	// Try the cache first. Make sure we have a non-expired entry and
	// that the value is a valid JWKS.
	entry, ok := getCache().Get(cacheKey)
	if ok && !entry.HasExpired() {
		jwks, ok = entry.GetValue().(*clerk.JSONWebKeySet)
		if !ok || jwks == nil || len(jwks.Keys) == 0 {
			jwks, err = forceGetJWKS(ctx, cacheKey)
			if err != nil {
				return nil, err
			}
		}
	} else {
		jwks, err = forceGetJWKS(ctx, cacheKey)
		if err != nil {
			return nil, err
		}
	}
	return jwks, err
}

// Fetches the JSON web key set from the API and caches it.
func forceGetJWKS(ctx context.Context, cacheKey string) (*clerk.JSONWebKeySet, error) {
	jwks, err := jwks.Get(ctx, &jwks.GetParams{})
	if err != nil {
		return nil, err
	}
	getCache().Set(cacheKey, jwks, time.Now().UTC().Add(time.Hour))
	return jwks, nil
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

// Caching store.
type cache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
}

// Get returns the cache entry for the provided key, if one exists.
func (c *cache) Get(key string) (*cacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return entry, ok
}

// Set adds a new entry with the provided value in the cache under
// the provided key. An expiration date will be set for the entry.
func (c *cache) Set(key string, value any, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

// A cache entry has a value and an expiration date.
type cacheEntry struct {
	value     any
	expiresAt time.Time
}

// HasExpired returns true if the cache entry's expiration date
// has passed.
func (entry *cacheEntry) HasExpired() bool {
	if entry == nil {
		return true
	}
	return entry.expiresAt.Before(time.Now())
}

// GetValue returns the cache entry's value.
func (entry *cacheEntry) GetValue() any {
	if entry == nil {
		return nil
	}
	return entry.value
}

var cacheInit sync.Once

// A "singleton" cache for the package.
var defaultCache *cache

// Lazy initialize and return the default cache singleton.
func getCache() *cache {
	cacheInit.Do(func() {
		defaultCache = &cache{
			entries: map[string]*cacheEntry{},
		}
	})
	return defaultCache
}
