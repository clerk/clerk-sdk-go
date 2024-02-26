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
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := &AuthorizationParams{}
			for _, opt := range opts {
				err := opt(params)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			if params.Clock == nil {
				params.Clock = clerk.NewClock()
			}

			authorization := strings.TrimSpace(r.Header.Get("Authorization"))
			if authorization == "" {
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(authorization, "Bearer ")
			decoded, err := jwt.Decode(r.Context(), &jwt.DecodeParams{Token: token})
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if params.JWK == nil {
				params.JWK, err = getJWK(r.Context(), params.JWKSClient, decoded.KeyID, params.Clock)
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
func getJWK(ctx context.Context, jwksClient *jwks.Client, kid string, clock clerk.Clock) (*clerk.JSONWebKey, error) {
	if kid == "" {
		return nil, fmt.Errorf("missing jwt kid header claim")
	}

	jwk := getCache().Get(kid)
	if jwk == nil || !getCache().IsValid(kid, clock.Now().UTC()) {
		var err error
		jwk, err = jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
			KeyID:      kid,
			JWKSClient: jwksClient,
		})
		if err != nil {
			return nil, err
		}
	}
	getCache().Set(kid, jwk, clock.Now().UTC().Add(time.Hour))
	return jwk, nil
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
func AuthorizedPartyMatches(parties ...string) AuthorizationOption {
	authorizedParties := make(map[string]struct{})
	for _, p := range parties {
		authorizedParties[p] = struct{}{}
	}

	return func(params *AuthorizationParams) error {
		params.AuthorizedPartyHandler = func(azp string) bool {
			if azp == "" || len(authorizedParties) == 0 {
				return true
			}
			_, ok := authorizedParties[azp]
			return ok
		}
		return nil
	}
}

// Clock allows to pass a clock implementation that will be the
// authority for time related operations.
// You can use a custom clock for testing purposes, or to
// eliminate clock skew if your code runs on different servers.
func Clock(c clerk.Clock) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.Clock = c
		return nil
	}
}

// CustomClaimsConstructor allows to pass a constructor function
// which returns a pointer to a type (struct) to hold custom token
// claims.
// The instance of the custom claims type will be then made available
// through the clerk.SessionClaims struct.
//
//	// Define a type to describe the custom claims.
//	type MyCustomClaims struct {
//		ACustomClaim string `json:"a_custom_claim"`
//	}
//
//	// In your HTTP server mux, configure the middleware with
//	// the custom claims constructor.
//	WithHeaderAuthorization(CustomClaimsConstructor(func(_ context.Context) any {
//		return &MyCustomClaims{}
//	})
//
//	// In the HTTP handler, access the active session claims. The
//	// custom claims are available in the SessionClaims.Custom field.
//	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
//	customClaims, ok := sessionClaims.Custom.(*MyCustomClaims)
func CustomClaimsConstructor(constructor func(context.Context) any) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.CustomClaimsConstructor = constructor
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

// IsValid returns true if a non-expired entry exists in the cache
// for the provided key, false otherwise.
func (c *jwkCache) IsValid(key string, t time.Time) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return ok && entry != nil && entry.expiresAt.After(t)
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
