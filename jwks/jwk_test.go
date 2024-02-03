package jwks

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestJWKGet(t *testing.T) {
	key := map[string]any{
		"use": "sig",
		"kty": "RSA",
		"kid": "the-kid",
		"alg": "RS256",
		"n":   "the-key",
		"e":   "AQAB",
	}
	out := map[string]any{
		"keys": []map[string]any{key},
	}
	raw, err := json.Marshal(out)
	require.NoError(t, err)

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    raw,
				Path:   "/v1/jwks",
				Method: http.MethodGet,
			},
		},
	}))
	jwk, err := Get(context.Background(), &GetParams{})
	require.NoError(t, err)
	require.Equal(t, 1, len(jwk.Keys))
	require.NotNil(t, jwk.Keys[0].Key)
	require.Equal(t, key["use"], jwk.Keys[0].Use)
	require.Equal(t, key["alg"], jwk.Keys[0].Algorithm)
	require.Equal(t, key["kid"], jwk.Keys[0].KeyID)
}
