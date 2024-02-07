package jwt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestVerify_InvalidToken(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{},
		},
	}))

	ctx := context.Background()
	_, err := Verify(ctx, &VerifyParams{
		Token: "this-is-not-a-token",
	})
	require.Error(t, err)
}

func TestVerify_Cache(t *testing.T) {
	ctx := context.Background()
	totalRequests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/v1/jwks" {
			totalRequests++
		}
		_, err := w.Write([]byte(`{
	"keys": [{
		"use": "sig",
		"kty": "RSA",
		"kid": "ins_123",
		"alg": "RS256",
		"n": "9m1LJW0dgEuK8SnN1Oy4LY8vaWABVS-hBTMA--_4LN1PZlMS5B2RPL85WkXYlHb0KXOSVrFKZLwYP-a9l3MFlW2YrPVAIvYfqPyqY5fmSEf-2qfrwosIhB2NSHyNRBQQ8-BX1RO9rIXIqYDKxGqktqMvYJmEGClmijbmFyQb2hpHD5PDbAB_DZvpZTEzWcQBL2ytHehILkYfg-ZZRyt7O8h5Gdy1v_TUlg8iMvchHlAkrIAmXNQigZmX_lne91tW8t4KMNJRfmUyLVCLbPnwxlmXXcice-0tmFw0OkCOteNWBeRNctJ3AIreGMzaJOJ2HeSUmJoX8iRKLLT3fsURLw",
		"e": "AQAB"
	}]
}`))
		require.NoError(t, err)
	}))
	defer ts.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: ts.Client(),
		URL:        clerk.String(ts.URL),
	}))

	token := "eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yOWR6bUdmQ3JydzdSMDRaVFFZRDNKSTB5dkYiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJodHRwczovL2Rhc2hib2FyZC5wcm9kLmxjbGNsZXJrLmNvbSIsImV4cCI6MTcwNzMwMDMyMiwiaWF0IjoxNzA3MzAwMjYyLCJpc3MiOiJodHRwczovL2NsZXJrLnByb2QubGNsY2xlcmsuY29tIiwibmJmIjoxNzA3MzAwMjUyLCJvcmdzIjp7Im9yZ18ySUlwcVIxenFNeHJQQkhSazNzTDJOSnJUQkQiOiJvcmc6YWRtaW4iLCJvcmdfMllHMlNwd0IzWEJoNUo0ZXF5elFVb0dXMjVhIjoib3JnOmFkbWluIiwib3JnXzJhZzJ6bmgxWGFjTXI0dGRXYjZRbEZSQ2RuaiI6Im9yZzphZG1pbiIsIm9yZ18yYWlldHlXa3VFSEhaRmRSUTFvVjYzMnZWaFciOiJvcmc6YWRtaW4ifSwic2lkIjoic2Vzc18yYm84b2gyRnIyeTNueVoyRVZQYktBd2ZvaU0iLCJzdWIiOiJ1c2VyXzI5ZTBXTnp6M245V1Q5S001WlpJYTBVVjNDNyJ9.6GtQafMBYY3Ij3pKHOyBYKt76LoLeBC71QUY_ho3k5nb0FBSvV0upKFLPBvIXNuF7hH0FK2QqDcAmrhbzAI-2qF_Ynve8Xl4VZCRpbTuZI7uL-tVjCvMffEIH-BHtrZ-QcXhEmNFQNIPyZTu21242he7U6o4S8st_aLmukWQzj_4qir7o5_fmVhm7YkLa0gYG5SLjkr2czwem1VGFHEVEOrHjun-g6eMnDNMMMysIOkZFxeqiCnqpc4u1V7Z7jfoK0r_-Unp8mGGln5KWYMCQyp1l1SkGwugtxeWfSbE4eklKRmItGOdVftvTyG16kDGpzsb22AQGtg65Iygni4PHg"
	// Providing a custom key will not trigger a request to fetch the
	// key set.
	_, _ = Verify(ctx, &VerifyParams{
		Token: token,
		JWK:   &clerk.JSONWebKey{},
	})
	require.Equal(t, 0, totalRequests)

	// Verify without providing a key. The method will trigger a request
	// to fetch the key set.
	_, _ = Verify(ctx, &VerifyParams{
		Token: token,
	})
	require.Equal(t, 1, totalRequests)
	// Verifying again won't trigger a request because the key set is
	// cached.
	_, _ = Verify(ctx, &VerifyParams{
		Token: token,
	})
	require.Equal(t, 1, totalRequests)
}
