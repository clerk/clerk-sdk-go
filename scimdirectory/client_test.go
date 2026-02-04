package scimdirectory

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

	enterpriseConnectionID := "entc_123"
	response := map[string]interface{}{
		"object":                   "scim_directory",
		"id":                       "scim_test123",
		"name":                     "Test SCIM Directory",
		"enterprise_connection_id": enterpriseConnectionID,
		"endpoint_url":             "https://scim.example.com",
		"provider":                 "okta",
		"enabled":                  true,
		"api_key":                  "sk_test_xxxxx",
		"created_at":               1640995200,
		"updated_at":               1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"enterprise_connection_id":"entc_123","name":"Test SCIM Directory","provider":"okta"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/scim_directories",
		},
	}

	client := NewClient(config)
	params := &CreateParams{
		EnterpriseConnectionID: clerk.String(enterpriseConnectionID),
		Name:                   clerk.String("Test SCIM Directory"),
		Provider:               clerk.String("okta"),
	}

	scimDirectory, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, "scim_test123", scimDirectory.ID)
	assert.Equal(t, "Test SCIM Directory", scimDirectory.Name)
	assert.Equal(t, "okta", scimDirectory.Provider)
	assert.True(t, scimDirectory.Enabled)
	assert.Equal(t, "sk_test_xxxxx", *scimDirectory.APIKey)
	assert.Equal(t, enterpriseConnectionID, *scimDirectory.EnterpriseConnectionID)
}

func TestGet(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":                   "scim_directory",
		"id":                       "scim_test123",
		"name":                     "Test SCIM Directory",
		"enterprise_connection_id": "entc_123",
		"endpoint_url":             "https://scim.example.com",
		"provider":                 "okta",
		"enabled":                  true,
		"created_at":               1640995200,
		"updated_at":               1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/scim_directories/scim_test123",
		},
	}

	client := NewClient(config)
	scimDirectory, err := client.Get(context.Background(), "scim_test123")
	require.NoError(t, err)
	assert.Equal(t, "scim_test123", scimDirectory.ID)
	assert.Equal(t, "Test SCIM Directory", scimDirectory.Name)
	assert.Equal(t, "okta", scimDirectory.Provider)
	assert.True(t, scimDirectory.Enabled)
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":                   "scim_directory",
		"id":                       "scim_test123",
		"name":                     "Updated SCIM Directory",
		"enterprise_connection_id": "entc_123",
		"endpoint_url":             "https://scim.example.com",
		"provider":                 "okta",
		"enabled":                  false,
		"created_at":               1640995200,
		"updated_at":               1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"enabled":false,"name":"Updated SCIM Directory"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPatch,
			Path:   "/v1/scim_directories/scim_test123",
		},
	}

	client := NewClient(config)
	params := &UpdateParams{
		Name:    clerk.String("Updated SCIM Directory"),
		Enabled: clerk.Bool(false),
	}

	scimDirectory, err := client.Update(context.Background(), "scim_test123", params)
	require.NoError(t, err)
	assert.Equal(t, "Updated SCIM Directory", scimDirectory.Name)
	assert.False(t, scimDirectory.Enabled)
}

func TestList(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"object":                   "scim_directory",
				"id":                       "scim_test123",
				"name":                     "Test SCIM Directory",
				"enterprise_connection_id": "entc_123",
				"endpoint_url":             "https://scim.example.com",
				"provider":                 "okta",
				"enabled":                  true,
				"created_at":               1640995200,
				"updated_at":               1640995200,
			},
			{
				"object":                   "scim_directory",
				"id":                       "scim_test456",
				"name":                     "Another SCIM Directory",
				"enterprise_connection_id": "entc_456",
				"endpoint_url":             "https://scim.example.com",
				"provider":                 "okta",
				"enabled":                  true,
				"created_at":               1640995200,
				"updated_at":               1640995200,
			},
		},
		"total_count": int64(2),
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/scim_directories",
			Query: &url.Values{
				"limit":  []string{"10"},
				"offset": []string{"0"},
			},
		},
	}

	client := NewClient(config)
	params := &ListParams{}
	params.Limit = clerk.Int64(10)
	params.Offset = clerk.Int64(0)

	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, int64(2), list.TotalCount)
	assert.Len(t, list.SCIMDirectories, 2)
	assert.Equal(t, "scim_test123", list.SCIMDirectories[0].ID)
	assert.Equal(t, "scim_test456", list.SCIMDirectories[1].ID)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	response := map[string]any{
		"object":  "scim_directory",
		"id":      "scim_test123",
		"deleted": true,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodDelete,
			Path:   "/v1/scim_directories/scim_test123",
		},
	}

	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), "scim_test123")
	require.NoError(t, err)
	assert.Equal(t, "scim_test123", deletedResource.ID)
	assert.True(t, deletedResource.Deleted)
}

func TestRotateAPIKey(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":                   "scim_directory",
		"id":                       "scim_test123",
		"name":                     "Test SCIM Directory",
		"enterprise_connection_id": "entc_123",
		"endpoint_url":             "https://scim.example.com",
		"provider":                 "okta",
		"enabled":                  true,
		"api_key":                  "sk_new_rotated_key",
		"created_at":               1640995200,
		"updated_at":               1640995200,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/scim_directories/scim_test123/rotate_api_key",
		},
	}

	client := NewClient(config)
	scimDirectory, err := client.RotateAPIKey(context.Background(), "scim_test123")
	require.NoError(t, err)
	assert.Equal(t, "scim_test123", scimDirectory.ID)
	assert.Equal(t, "Test SCIM Directory", scimDirectory.Name)
	assert.Equal(t, "sk_new_rotated_key", *scimDirectory.APIKey)
}
