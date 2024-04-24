package testingtoken

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestTestingTokenClientCreate(t *testing.T) {
	t.Parallel()

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: nil,
			Out: json.RawMessage(
				`{"object":"testing_token","token":"1713877200-c_3G2MvPu9PnXcuhbPZNao0LOXqK9A7YrnBn0HmIWxt","expires_at":1713880800}`,
			),
			Method: http.MethodPost,
			Path:   "/v1/testing_tokens",
		},
	}
	client := NewClient(config)
	token, err := client.Create(context.Background())
	require.NoError(t, err)
	require.Equal(t, "testing_token", token.Object)
	require.Equal(t, "1713877200-c_3G2MvPu9PnXcuhbPZNao0LOXqK9A7YrnBn0HmIWxt", token.Token)
	require.Equal(t, int64(1713880800), token.ExpiresAt)
}
