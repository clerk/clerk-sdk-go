package apikey

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
		"object":            "api_key",
		"id":                "ak_test123",
		"type":              "api_key",
		"subject":           "user_123",
		"name":              "test_key",
		"description":       "Test API key",
		"claims":            map[string]interface{}{},
		"scopes":            []string{},
		"secret":            "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		"revoked":           false,
		"revocation_reason": nil,
		"expired":           false,
		"expiration":        nil,
		"created_by":        "user_123",
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"name":"test_key","subject":"user_123","description":"Test API key"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/api_keys",
		},
	}

	client := NewClient(config)
	params := &CreateParams{
		Name:        clerk.String("test_key"),
		Subject:     clerk.String("user_123"),
		Description: clerk.String("Test API key"),
	}

	apiKey, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", apiKey.ID)
	assert.Equal(t, "test_key", apiKey.Name)
	assert.Equal(t, "user_123", apiKey.Subject)
	assert.Equal(t, "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", apiKey.Secret)
}

func TestList(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"api_keys": []map[string]interface{}{
			{
				"object":            "api_key",
				"id":                "ak_test123",
				"type":              "api_key",
				"subject":           "user_123",
				"name":              "test_key",
				"description":       "Test API key",
				"claims":            map[string]interface{}{},
				"scopes":            []string{},
				"revoked":           false,
				"revocation_reason": nil,
				"expired":           false,
				"expiration":        nil,
				"created_by":        "user_123",
				"last_used_at":      nil,
				"created_at":        1640995200,
				"updated_at":        1640995200,
			},
		},
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/api_keys",
			Query: &url.Values{
				"subject": []string{"user_123"},
				"type":    []string{"api_key"},
			},
		},
	}

	client := NewClient(config)
	params := &ListParams{
		Subject: clerk.String("user_123"),
		Type:    clerk.String("api_key"),
	}

	apiKeys, err := client.List(context.Background(), params)
	require.NoError(t, err)
	assert.Len(t, apiKeys.APIKeys, 1)
	assert.Equal(t, "ak_test123", apiKeys.APIKeys[0].ID)
}

func TestGetSecret(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"secret": "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/api_keys/ak_test123/secret",
		},
	}

	client := NewClient(config)
	secret, err := client.GetSecret(context.Background(), "ak_test123")
	require.NoError(t, err)
	assert.Equal(t, "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", secret.Secret)
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "api_key",
		"id":                "ak_test123",
		"type":              "api_key",
		"subject":           "user_123",
		"name":              "test_key",
		"description":       "Updated description",
		"claims":            map[string]interface{}{},
		"scopes":            []string{},
		"revoked":           false,
		"revocation_reason": nil,
		"expired":           false,
		"expiration":        nil,
		"created_by":        "user_123",
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"description":"Updated description"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPatch,
			Path:   "/v1/api_keys/ak_test123",
		},
	}

	client := NewClient(config)
	params := &UpdateParams{
		Description: clerk.String("Updated description"),
	}

	apiKey, err := client.Update(context.Background(), "ak_test123", params)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", *apiKey.Description)
}

func TestRevoke(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "api_key",
		"id":                "ak_test123",
		"type":              "api_key",
		"subject":           "user_123",
		"name":              "test_key",
		"description":       "Test API key",
		"claims":            map[string]interface{}{},
		"scopes":            []string{},
		"revoked":           true,
		"revocation_reason": "Security breach",
		"expired":           false,
		"expiration":        nil,
		"created_by":        "user_123",
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
			Path:   "/v1/api_keys/ak_test123/revoke",
		},
	}

	client := NewClient(config)
	params := &RevokeParams{
		RevocationReason: clerk.String("Security breach"),
	}

	apiKey, err := client.Revoke(context.Background(), "ak_test123", params)
	require.NoError(t, err)
	assert.True(t, apiKey.Revoked)
	assert.Equal(t, "Security breach", *apiKey.RevocationReason)
}

func TestVerify(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "api_key",
		"id":                "ak_test123",
		"type":              "api_key",
		"subject":           "user_123",
		"name":              "test_key",
		"description":       "Test API key",
		"claims":            map[string]interface{}{},
		"scopes":            []string{},
		"revoked":           false,
		"revocation_reason": nil,
		"expired":           false,
		"expiration":        nil,
		"created_by":        "user_123",
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"secret":"ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/api_keys/verify",
		},
	}

	client := NewClient(config)
	params := &VerifyParams{
		Secret: "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	apiKey, err := client.Verify(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", apiKey.ID)
	assert.False(t, apiKey.Revoked)
}
