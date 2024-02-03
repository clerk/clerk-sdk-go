package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

func TestWithHeaderAuthorization_InvalidAuthorization(t *testing.T) {
	ts := httptest.NewServer(WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := clerk.SessionClaimsFromContext(r.Context())
		require.False(t, ok)
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
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Request with invalid Authorization header
	req.Header.Add("authorization", "Bearer whatever")
	res, err = ts.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
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
