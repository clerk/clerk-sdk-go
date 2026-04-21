package m2m_token

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "machine_to_machine_token",
		"id":                "mt_test123",
		"subject":           "mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
		"claims":            map[string]interface{}{"foo": "bar"},
		"scopes":            []string{"mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk"},
		"token":             "mt_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		"revoked":           false,
		"revocation_reason": nil,
		"expired":           false,
		"expiration":        1716883200,
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"claims":{"foo":"bar"},"seconds_until_expiration":3600,"token_format":"opaque"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/m2m_tokens",
		},
	}

	client := NewClient(config)
	claims := json.RawMessage(`{"foo":"bar"}`)
	params := &CreateParams{
		Claims:                 &claims,
		SecondsUntilExpiration: clerk.Int64(3600),
		TokenFormat:            clerk.String("opaque"),
	}

	token, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "mt_test123", token.ID)
	assert.Equal(t, "mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk", token.Subject)
	assert.Equal(t, "mt_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", token.Token)
}

func TestList(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"m2m_tokens": []map[string]interface{}{
			{
				"object":            "machine_to_machine_token",
				"id":                "mt_test123",
				"subject":           "mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
				"claims":            map[string]interface{}{"foo": "bar"},
				"scopes":            []string{"mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk"},
				"revoked":           false,
				"revocation_reason": nil,
				"expired":           false,
				"expiration":        1716883200,
				"last_used_at":      nil,
				"created_at":        1640995200,
				"updated_at":        1640995200,
			},
		},
		"total_count": 1,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/m2m_tokens",
			Query: &url.Values{
				"subject": []string{"mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk"},
				"revoked": []string{"false"},
			},
		},
	}

	client := NewClient(config)
	params := &ListParams{
		Subject: clerk.String("mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk"),
		Revoked: clerk.Bool(false),
	}

	tokens, err := client.List(context.Background(), params)
	require.NoError(t, err)
	assert.Len(t, tokens.M2MTokens, 1)
	assert.Equal(t, "mt_test123", tokens.M2MTokens[0].ID)
	assert.Equal(t, int64(1), tokens.TotalCount)
}

func TestRevoke(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "machine_to_machine_token",
		"id":                "mt_test123",
		"subject":           "mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
		"claims":            map[string]interface{}{"foo": "bar"},
		"scopes":            []string{"mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk"},
		"revoked":           true,
		"revocation_reason": "Security breach",
		"expired":           false,
		"expiration":        1716883200,
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"revocation_reason":"Security breach"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/m2m_tokens/mt_test123/revoke",
		},
	}

	client := NewClient(config)
	params := &RevokeParams{
		RevocationReason: clerk.String("Security breach"),
	}

	token, err := client.Revoke(context.Background(), "mt_test123", params)
	require.NoError(t, err)
	assert.True(t, token.Revoked)
	assert.Equal(t, "Security breach", *token.RevocationReason)
}

func TestVerify(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "machine_to_machine_token",
		"id":                "mt_test123",
		"subject":           "mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
		"claims":            map[string]interface{}{"foo": "bar"},
		"scopes":            []string{"mch_2xhFjEI5X2qWRvtV13BzSj8H6Dk"},
		"revoked":           false,
		"revocation_reason": nil,
		"expired":           false,
		"expiration":        1716883200,
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"token":"mt_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/m2m_tokens/verify",
		},
	}

	client := NewClient(config)
	params := &VerifyParams{
		Token: "mt_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	token, err := client.Verify(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "mt_test123", token.ID)
	assert.False(t, token.Revoked)
}
