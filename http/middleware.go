// Package http provides HTTP utilities and handler middleware.
package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	jose "github.com/go-jose/go-jose/v3/jwt"
)

// RequireHeaderAuthorization will respond with HTTP 403 Forbidden if
// the Authorization header does not contain a valid session token.
func RequireHeaderAuthorization(opts ...AuthorizationOption) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WithHeaderAuthorization(opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok || claims == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

// WithHeaderAuthorization checks the Authorization request header
// for a valid Clerk authorization JWT. The token is parsed and verified
// and the active session claims are written to the http.Request context.
// The middleware uses Bearer authentication, so the Authorization header
// is expected to have the following format:
// Authorization: Bearer <token>
func WithHeaderAuthorization(opts ...AuthorizationOption) func(http.Handler) http.Handler {
	var paramsErr error
	params := &AuthorizationParams{}
	for _, opt := range opts {
		paramsErr = opt(params)
		if paramsErr != nil {
			break
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if paramsErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			authorization := strings.TrimSpace(r.Header.Get("Authorization"))
			if authorization == "" {
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(authorization, "Bearer ")
			_, err := jwt.Decode(r.Context(), &jwt.DecodeParams{Token: token})
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if params.JWK == nil {
				params.JWK, err = getJWK(r.Context(), params.JWKSClient, token)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			params.Token = token
			claims, err := jwt.Verify(r.Context(), &params.VerifyParams)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Token was verified. Add the session claims to the request context.
			newCtx := clerk.ContextWithSessionClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// Retrieve the JSON web key for the provided token from the JWKS set.
// Tries a cached value first, but if there's no value or the entry
// has expired, it will fetch the JWK set from the API and cache the
// value.
func getJWK(ctx context.Context, jwksClient *jwks.Client, token string) (*clerk.JSONWebKey, error) {
	kid, err := getKID(token)
	if err != nil {
		return nil, err
	}

	jwk := getCache().Get(kid)
	if jwk == nil {
		jwk, err = forceGetJWK(ctx, jwksClient, kid)
		if err != nil {
			return nil, err
		}
	}
	getCache().Set(kid, jwk, time.Now().UTC().Add(time.Hour))
	return jwk, nil
}

// Fetch the JSON Web Key Set from the Clerk API and return the JSON
// Web Key corresponding to the provided KeyID.
// A default client will be initialized if the provided jwks.Client
// is nil.
func forceGetJWK(ctx context.Context, jwksClient *jwks.Client, kid string) (*clerk.JSONWebKey, error) {
	if jwksClient == nil {
		jwksClient = &jwks.Client{
			Backend: clerk.GetBackend(),
		}
	}
	jwks, err := jwksClient.Get(ctx, &jwks.GetParams{})
	if err != nil {
		return nil, err
	}
	if jwks == nil || len(jwks.Keys) == 0 {
		return nil, fmt.Errorf("no jwks found")
	}
	for _, k := range jwks.Keys {
		if k.KeyID == kid {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("no jwk key found for kid %s", kid)
}

// Extract the KeyID claim from the token.
func getKID(token string) (string, error) {
	parsedToken, err := jose.ParseSigned(token)
	if err != nil {
		return "", err
	}
	if len(parsedToken.Headers) == 0 {
		return "", fmt.Errorf("missing jwt headers")
	}
	kid := parsedToken.Headers[0].KeyID
	if kid == "" {
		return "", fmt.Errorf("missing jwt kid header claim")
	}
	return kid, nil
}

type AuthorizationParams struct {
	jwt.VerifyParams
	// JWKSClient is the jwks.Client that will be used to fetch the
	// JSON Web Key Set. A default client will be used if none is
	// provided.
	JWKSClient *jwks.Client
}

// AuthorizationOption is a functional parameter for configuring
// authorization options.
type AuthorizationOption func(*AuthorizationParams) error

// AuthorizedParty allows to provide a handler that accepts the
// 'azp' claim.
// The handler can be used to perform validations on the azp claim
// and should return false to indicate that something is wrong.
func AuthorizedParty(handler func(string) bool) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.AuthorizedPartyHandler = handler
		return nil
	}
}

// AuthorizedPartyMatches registers a handler that checks that the
// 'azp' claim's value is included in the provided parties.
func AuthorizedPartyMatches(parties ...string) func(string) bool {
	authorizedParties := make(map[string]struct{})
	for _, p := range parties {
		authorizedParties[p] = struct{}{}
	}

	return func(azp string) bool {
		if azp == "" || len(authorizedParties) == 0 {
			return true
		}
		_, ok := authorizedParties[azp]
		return ok
	}
}

// CustomClaims allows to pass a type (e.g. struct), which will be populated with the token claims based on json tags.
// You must pass a pointer for this option to work.
func CustomClaims(claims any) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.CustomClaims = claims
		return nil
	}
}

// Leeway allows to set a custom leeway when comparing time values
// for JWT verification.
// The leeway gives some extra time to the token. That is, if the
// token is expired, it will still be accepted for 'leeway' amount
// of time.
// This option accomodates for clock skew.
func Leeway(leeway time.Duration) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.Leeway = leeway
		return nil
	}
}

// ProxyURL can be used to set the URL that proxies the Clerk Frontend
// API. Useful for proxy based setups.
// See https://clerk.com/docs/advanced-usage/using-proxies
func ProxyURL(proxyURL string) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.ProxyURL = clerk.String(proxyURL)
		return nil
	}
}

// Satellite can be used to signify that the authorization happens
// on a satellite domain.
// See https://clerk.com/docs/advanced-usage/satellite-domains
func Satellite(isSatellite bool) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.IsSatellite = isSatellite
		return nil
	}
}

// JSONWebKey allows to provide a custom JSON Web Key (JWK) based on
// which the authorization JWT will be verified.
// When verifying the authorization JWT without a custom key, the JWK
// will be fetched from the Clerk API and cached for one hour, then
// the JWK will be fetched again from the Clerk API.
// Passing a custom JSON Web Key means that no request to fetch JSON
// web keys will be made. It's the caller's responsibility to refresh
// the JWK when keys are rolled.
func JSONWebKey(key string) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		// From the Clerk docs: "Note that the JWT Verification key is not in
		// PEM format, the header and footer are missing, in order to be shorter
		// and single-line for easier setup."
		if !strings.HasPrefix(key, "-----BEGIN") {
			key = "-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----"
		}
		jwk, err := clerk.JSONWebKeyFromPEM(key)
		if err != nil {
			return err
		}
		params.JWK = jwk
		return nil
	}
}

// JWKSClient allows to provide a custom jwks.Client that will be
// used when fetching the JSON Web Key Set with which the JWT
// will be verified.
// The JSONWebKey option takes precedence. If a web key is already
// provided through the JSONWebKey option, the JWKS client will
// not be used at all.
func JWKSClient(client *jwks.Client) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.JWKSClient = client
		return nil
	}
}

// A cache to store JSON Web Keys.
type jwkCache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
}

// Each entry in the JWK cache has a value and an expiration date.
type cacheEntry struct {
	value     *clerk.JSONWebKey
	expiresAt time.Time
}

// Get fetches the JSON Web Key for the provided key, unless the
// entry has expired.
func (c *jwkCache) Get(key string) *clerk.JSONWebKey {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || entry == nil {
		return nil
	}
	if entry.expiresAt.Before(time.Now().UTC()) {
		return nil
	}
	return entry.value
}

// Set stores the JSON Web Key in the provided key and sets the
// expiration date.
func (c *jwkCache) Set(key string, value *clerk.JSONWebKey, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

var cacheInit sync.Once

// A "singleton" JWK cache for the package.
var cache *jwkCache

// getCache returns the library's default cache singleton.
// Please note that the returned Cache is a package-level variable.
// Using the package with more than one Clerk API secret keys might
// require to use different Clients with their own Cache.
func getCache() *jwkCache {
	cacheInit.Do(func() {
		cache = &jwkCache{
			entries: map[string]*cacheEntry{},
		}
	})
	return cache
}
