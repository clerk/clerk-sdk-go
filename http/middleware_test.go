package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestWithHeaderAuthorization_InvalidAuthorization(t *testing.T) {
	kid := "kid-" + t.Name()
	// Mock the Clerk API server. We expect requests to GET /jwks.
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
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
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	require.NoError(t, err)
	res, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
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
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
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

func TestWithHeaderAuthorization_CustomFailureHandler(t *testing.T) {
	kid := "kid-" + t.Name()
	// Mock the Clerk API server. We expect requests to GET /jwks.
	clerkAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
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
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
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

func TestAuthorizedJWTExtractor(t *testing.T) {
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
