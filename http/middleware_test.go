package http

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestWithHeaderAuthorization_InvalidAuthorization(t *testing.T) {
	kid := "kid-" + t.Name()
	// Mock the Clerk API server. We expect requests to GET /jwks.
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			_, err := fmt.Fprintf(w,
				`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
				kid)
			require.NoError(t, err)
			return
		}
	}))
	defer clerkAPI.Close()

	// Mock the clerk backend
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: clerkAPI.Client(),
		URL:        &clerkAPI.URL,
	}))

	// This is the user's server, guarded by Clerk's middleware.
	ts := httptest.NewServer(WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := clerk.SessionClaimsFromContext(r.Context())
		require.False(t, ok)
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})))
	defer ts.Close()

	// Request without Authorization header
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Request with invalid Authorization header
	req.Header.Add("authorization", "Bearer whatever")
	res, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Request with unverifiable Bearer token
	tokenClaims := map[string]any{
		"sid": "sess_123",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, clerktest.WithKID(kid))
	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	require.NoError(t, err)
	res, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestWithHeaderAuthorization_JWTHeaderType(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	kid := "kid-" + t.Name()
	jwksBody := clerktest.ConvertToJWKS(t, priv.Public(), kid)

	// Mock the Clerk API server. We expect requests to GET /jwks.
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			_, err := w.Write(jwksBody)
			require.NoError(t, err)
			return
		}
	}))
	defer clerkAPI.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: clerkAPI.Client(),
		URL:        &clerkAPI.URL,
	}))

	tokenClaims := map[string]any{
		"sid": "sess_123",
		"sub": "user_123",
		"iss": "https://clerk.com",
	}

	t.Run("default accepts typ JWT", func(t *testing.T) {
		ts := httptest.NewServer(WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			require.True(t, ok)
			require.Equal(t, "user_123", claims.Subject)
			_, err := w.Write([]byte("{}"))
			require.NoError(t, err)
		})))
		defer ts.Close()

		// Generate a valid JWT with the default typ "JWT"
		tokenGood, _ := clerktest.GenerateJWT(t, tokenClaims,
			clerktest.WithType("JWT"), clerktest.WithRSAPrivateKey(priv), clerktest.WithKID(kid))
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenGood)
		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("default rejects wrong typ", func(t *testing.T) {
		ts := httptest.NewServer(WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("{}"))
			require.NoError(t, err)
		})))
		defer ts.Close()

		// Generate a valid JWT with the wrong typ "at+jwt"
		tokenWrongTyp, _ := clerktest.GenerateJWT(t, tokenClaims,
			clerktest.WithType("at+jwt"), clerktest.WithRSAPrivateKey(priv), clerktest.WithKID(kid))
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenWrongTyp)
		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("ignore succeeds for wrong typ", func(t *testing.T) {
		// Setup middleware to ignore the typ header
		ts := httptest.NewServer(WithHeaderAuthorization(IgnoreJWTHeaderType())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			require.True(t, ok)
			require.Equal(t, "user_123", claims.Subject)
			_, err := w.Write([]byte("{}"))
			require.NoError(t, err)
		})))
		defer ts.Close()

		// Generate a valid JWT with the wrong typ "at+jwt"
		tokenWrongTyp, _ := clerktest.GenerateJWT(t, tokenClaims,
			clerktest.WithType("at+jwt"), clerktest.WithRSAPrivateKey(priv), clerktest.WithKID(kid))
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenWrongTyp)
		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("custom expected typ matches", func(t *testing.T) {
		// Setup middleware to expect custom type "custom+abc"
		customType := "custom+abc"
		ts := httptest.NewServer(WithHeaderAuthorization(ExpectedJWTHeaderType(customType))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			require.True(t, ok)
			require.Equal(t, "user_123", claims.Subject)
			_, err := w.Write([]byte("{}"))
			require.NoError(t, err)
		})))
		defer ts.Close()

		// Generate a valid JWT with the typ "custom+abc"
		tokenAtJWT, _ := clerktest.GenerateJWT(t, tokenClaims,
			clerktest.WithType(customType), clerktest.WithRSAPrivateKey(priv), clerktest.WithKID(kid))
		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenAtJWT)
		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestRequireHeaderAuthorization_InvalidAuthorization(t *testing.T) {
	ts := httptest.NewServer(RequireHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})))
	defer ts.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	// Request without Authorization header
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, res.StatusCode)

	// Request with invalid Authorization header
	req.Header.Add("authorization", "Bearer whatever")
	res, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestWithHeaderAuthorization_Caching(t *testing.T) {
	kid := "kid-" + t.Name()
	clock := clerktest.NewClockAt(time.Now().UTC())

	// Mock the Clerk API server. We expect requests to GET /jwks.
	totalJWKSRequests := 0
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			// Count the number of requests to the JWKS endpoint
			totalJWKSRequests++
			_, err := fmt.Fprintf(w,
				`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
				kid)
			require.NoError(t, err)
			return
		}
	}))
	defer clerkAPI.Close()

	// Mock the clerk backend
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: clerkAPI.Client(),
		URL:        &clerkAPI.URL,
	}))

	// This is the user's server, guarded by Clerk's http middleware.
	ts := httptest.NewServer(WithHeaderAuthorization(Clock(clock))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})))
	defer ts.Close()

	// Generate a token with the claims below.
	tokenClaims := map[string]any{
		"sid": "sess_123",
		"sub": "user_123",
		"iss": "https://clerk.com",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, clerktest.WithKID(kid))
	// The first request needs to fetch the JSON web key set, because
	// the cache is empty.
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	require.NoError(t, err)
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, totalJWKSRequests)

	// The next request will use the cached value
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, totalJWKSRequests)

	// If we move past the cache's expiry date, the JWKS will be fetched again.
	clock.Advance(2 * time.Hour)
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 2, totalJWKSRequests)

	// The next time the JWKS will be cached again.
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 2, totalJWKSRequests)
}

// Regression test: cache entries must expire at a fixed TTL from their
// original fetch time, not on a sliding window that resets on every
// cache hit. Otherwise a revoked key could remain trusted indefinitely
// on a server that receives at least one request per TTL window.
func TestWithHeaderAuthorization_CacheFixedTTL(t *testing.T) {
	kid := "kid-" + t.Name()
	clock := clerktest.NewClockAt(time.Now().UTC())

	totalJWKSRequests := 0
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			totalJWKSRequests++
			_, err := w.Write([]byte(
				fmt.Sprintf(
					`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
					kid,
				),
			))
			require.NoError(t, err)
			return
		}
	}))
	defer clerkAPI.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: clerkAPI.Client(),
		URL:        &clerkAPI.URL,
	}))

	ts := httptest.NewServer(WithHeaderAuthorization(Clock(clock))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})))
	defer ts.Close()

	tokenClaims := map[string]any{
		"sid": "sess_123",
		"sub": "user_123",
		"iss": "https://clerk.com",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, clerktest.WithKID(kid))
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	require.NoError(t, err)

	// Initial fetch populates the cache at t=0. Expiry is t=60min.
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, totalJWKSRequests)

	// At t=30min, cache is still valid and must be used without refetching.
	// Critically, this hit must NOT extend the expiry (no sliding window).
	clock.Advance(30 * time.Minute)
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, totalJWKSRequests)

	// At t=65min, the original TTL has elapsed. If the previous cache hit
	// had extended the expiry to t=90min (sliding window), this request
	// would still use the stale cache and totalJWKSRequests would be 1.
	// With a fixed TTL, the cache is expired and the JWKS is refetched.
	clock.Advance(35 * time.Minute)
	_, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, 2, totalJWKSRequests)
}

func TestWithHeaderAuthorization_CustomFailureHandler(t *testing.T) {
	kid := "kid-" + t.Name()
	// Mock the Clerk API server. We expect requests to GET /jwks.
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			_, err := fmt.Fprintf(w,
				`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
				kid)
			require.NoError(t, err)
			return
		}
	}))
	defer clerkAPI.Close()

	// Define a custom failure handler which returns a custom HTTP
	// status code.
	customFailureHandler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}

	// Apply the custom failure handler to the WithHeaderAuthorization
	// middleware.
	middleware := WithHeaderAuthorization(
		AuthorizationFailureHandler(http.HandlerFunc(customFailureHandler)),
	)
	// This is the user's server, guarded by Clerk's http middleware.
	ts := httptest.NewServer(middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := clerk.SessionClaimsFromContext(r.Context())
		require.False(t, ok)
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})))
	defer ts.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: clerkAPI.Client(),
		URL:        &clerkAPI.URL,
	}))

	tokenClaims := map[string]any{
		"sid": "sess_123",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, clerktest.WithKID(kid))
	// Request with invalid Authorization header
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTeapot, res.StatusCode)
}

func TestAuthorizedPartyFunc(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		azp     string
		parties []string
		want    bool
	}{
		{
			azp:     "clerk.com",
			parties: []string{"clerk.com", "clerk.dev"},
			want:    true,
		},
		{
			azp:     "clerk.com",
			parties: []string{"clerk.dev"},
			want:    false,
		},
		{
			azp:     "",
			parties: []string{"clerk.com"},
			want:    true,
		},
		{
			azp:     "clerk.com",
			parties: []string{},
			want:    true,
		},
	} {
		options := &AuthorizationParams{}
		err := AuthorizedPartyMatches(tc.parties...)(options)
		require.NoError(t, err)
		require.Equal(t, tc.want, options.AuthorizedPartyHandler(tc.azp))
	}
}

func TestAuthorizationJWTExtractor(t *testing.T) {
	middleware := RequireHeaderAuthorization(AuthorizationJWTExtractor(func(r *http.Request) string {
		return r.Header.Get("X-Clerk-JWT-Test")
	}))
	ts := httptest.NewServer(middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})))
	defer ts.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	// Request without JWT
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, res.StatusCode)

	// Request with invalid JWT
	req.Header.Add("X-Clerk-JWT-Test", "whatever")
	res, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestNeedsSessionReverification(t *testing.T) {
	for _, tc := range []struct {
		name           string
		factorAges     [2]int64
		policy         clerk.SessionReverificationPolicy
		expectedStatus int
	}{
		{
			name:       "first factor - valid",
			factorAges: [2]int64{5, -1},
			policy: clerk.SessionReverificationPolicy{
				AfterMinutes: 10,
				Level:        clerk.SessionReverificationLevelFirstFactor,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "first factor - needs reverification",
			factorAges: [2]int64{15, -1},
			policy: clerk.SessionReverificationPolicy{
				AfterMinutes: 10,
				Level:        clerk.SessionReverificationLevelFirstFactor,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:       "second factor - valid",
			factorAges: [2]int64{20, 5},
			policy: clerk.SessionReverificationPolicy{
				AfterMinutes: 10,
				Level:        clerk.SessionReverificationLevelSecondFactor,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "second factor - needs reverification",
			factorAges: [2]int64{20, 15},
			policy: clerk.SessionReverificationPolicy{
				AfterMinutes: 10,
				Level:        clerk.SessionReverificationLevelSecondFactor,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:       "multi factor - valid",
			factorAges: [2]int64{5, -1},
			policy: clerk.SessionReverificationPolicy{
				AfterMinutes: 10,
				Level:        clerk.SessionReverificationLevelSecondFactor,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "multi factor - needs reverification",
			factorAges: [2]int64{5, 15},
			policy: clerk.SessionReverificationPolicy{
				AfterMinutes: 10,
				Level:        clerk.SessionReverificationLevelMultiFactor,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "using predefined policy",
			factorAges:     [2]int64{15, 15},
			policy:         clerk.SessionReverificationStrictMFA,
			expectedStatus: http.StatusForbidden,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			includeSessionClaims := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					claims := &clerk.SessionClaims{
						Claims: clerk.Claims{
							FactorVerificationAge: tc.factorAges,
						},
					}
					ctx := clerk.ContextWithSessionClaims(r.Context(), claims)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			}

			// This is the user's server, using the NeedsSessionReverification middleware.
			ts := httptest.NewServer(includeSessionClaims(NeedsSessionReverification(tc.policy)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("{}"))
				require.NoError(t, err)
			}))))
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			// Send the request
			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify the response status
			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				return
			}

			// Verify the error response has the expected structure and reason
			var errResp ClerkErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&errResp)
			require.NoError(t, err)
			require.Equal(t, "reverification-error", errResp.ClerkError.Reason)
		})
	}
}
