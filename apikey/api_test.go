package apikey

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageCreate(t *testing.T) {
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

	// Set up the backend with our mock transport
	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(`{"name":"test_key","subject":"user_123","description":"Test API key"}`),
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodPost,
				Path:   "/v1/api_keys",
			},
		},
	}

	// Set the backend globally
	clerk.SetBackend(clerk.NewBackend(backend))

	params := &CreateParams{
		Name:        clerk.String("test_key"),
		Subject:     clerk.String("user_123"),
		Description: clerk.String("Test API key"),
	}

	apiKey, err := Create(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", apiKey.ID)
	assert.Equal(t, "test_key", apiKey.Name)
	assert.Equal(t, "user_123", apiKey.Subject)
	assert.Equal(t, "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", apiKey.Secret)
}

func TestPackageGet(t *testing.T) {
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

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodGet,
				Path:   "/v1/api_keys/ak_test123",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	apiKey, err := Get(context.Background(), "ak_test123")
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", apiKey.ID)
	assert.Equal(t, "test_key", apiKey.Name)
	assert.Equal(t, "user_123", apiKey.Subject)
	assert.Equal(t, "Test API key", *apiKey.Description)
	assert.False(t, apiKey.Revoked)
}

func TestPackageList(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{
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
		"total_count": int64(1),
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
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
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &ListParams{
		Subject: clerk.String("user_123"),
		Type:    clerk.String("api_key"),
	}

	apiKeys, err := List(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, int64(1), apiKeys.TotalCount)
	assert.Len(t, apiKeys.APIKeys, 1)
	assert.Equal(t, "ak_test123", apiKeys.APIKeys[0].ID)
}

func TestPackageGetSecret(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"secret": "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodGet,
				Path:   "/v1/api_keys/ak_test123/secret",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	secret, err := GetSecret(context.Background(), "ak_test123")
	require.NoError(t, err)
	assert.Equal(t, "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", secret.Secret)
}

func TestPackageUpdate(t *testing.T) {
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

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(`{"description":"Updated description"}`),
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodPatch,
				Path:   "/v1/api_keys/ak_test123",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &UpdateParams{
		Description: clerk.String("Updated description"),
	}

	apiKey, err := Update(context.Background(), "ak_test123", params)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", *apiKey.Description)
}

func TestPackageDelete(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":  "api_key",
		"id":      "ak_test123",
		"deleted": true,
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodDelete,
				Path:   "/v1/api_keys/ak_test123",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	deletedResource, err := Delete(context.Background(), "ak_test123")
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", deletedResource.ID)
	assert.True(t, deletedResource.Deleted)
	assert.Equal(t, "api_key", deletedResource.Object)
}

func TestPackageUpdateWithSecondsUntilExpiration(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":            "api_key",
		"id":                "ak_test123",
		"type":              "api_key",
		"subject":           "user_123",
		"name":              "test_key",
		"description":       "Updated with expiration",
		"claims":            map[string]interface{}{},
		"scopes":            []string{},
		"revoked":           false,
		"revocation_reason": nil,
		"expired":           false,
		"expiration":        1716883200,
		"created_by":        "user_123",
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(`{"description":"Updated with expiration","seconds_until_expiration":3600}`),
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodPatch,
				Path:   "/v1/api_keys/ak_test123",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &UpdateParams{
		Description:            clerk.String("Updated with expiration"),
		SecondsUntilExpiration: clerk.Int64(3600),
	}

	apiKey, err := Update(context.Background(), "ak_test123", params)
	require.NoError(t, err)
	assert.Equal(t, "Updated with expiration", *apiKey.Description)
	assert.Equal(t, int64(1716883200), *apiKey.Expiration)
}

func TestPackageUpdateSecondsUntilExpirationOnly(t *testing.T) {
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
		"expiration":        1716883200,
		"created_by":        "user_123",
		"last_used_at":      nil,
		"created_at":        1640995200,
		"updated_at":        1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(`{"seconds_until_expiration":3600}`),
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodPatch,
				Path:   "/v1/api_keys/ak_test123",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &UpdateParams{
		SecondsUntilExpiration: clerk.Int64(3600),
	}

	apiKey, err := Update(context.Background(), "ak_test123", params)
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", apiKey.ID)
	assert.Equal(t, int64(1716883200), *apiKey.Expiration)
}

func TestPackageRevoke(t *testing.T) {
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

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(`{"revocation_reason":"Security breach"}`),
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodPost,
				Path:   "/v1/api_keys/ak_test123/revoke",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &RevokeParams{
		RevocationReason: clerk.String("Security breach"),
	}

	apiKey, err := Revoke(context.Background(), "ak_test123", params)
	require.NoError(t, err)
	assert.True(t, apiKey.Revoked)
	assert.Equal(t, "Security breach", *apiKey.RevocationReason)
}

func TestPackageVerify(t *testing.T) {
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

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(`{"secret":"ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`),
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodPost,
				Path:   "/v1/api_keys/verify",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &VerifyParams{
		Secret: "ak_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	apiKey, err := Verify(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "ak_test123", apiKey.ID)
	assert.False(t, apiKey.Revoked)
}
