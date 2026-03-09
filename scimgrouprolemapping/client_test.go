package scimgrouprolemapping

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	response := []map[string]interface{}{
		{
			"object":                  "scim_group_role_mapping",
			"id":                      "scim_grp_role_123",
			"scim_directory_id":       "scim_dir_123",
			"scim_group_id":           "group_123",
			"scim_group_display_name": "Admins",
			"role": map[string]interface{}{
				"object":      "role",
				"id":          "role_admin",
				"name":        "Admin",
				"key":         "org:admin",
				"description": "Administrator role",
				"permissions": []string{},
			},
			"precedence": 1,
			"created_at": 1640995200000,
			"updated_at": 1640995200000,
		},
		{
			"object":                  "scim_group_role_mapping",
			"id":                      "scim_grp_role_456",
			"scim_directory_id":       "scim_dir_123",
			"scim_group_id":           "group_456",
			"scim_group_display_name": "Members",
			"role": map[string]interface{}{
				"object":      "role",
				"id":          "role_member",
				"name":        "Member",
				"key":         "org:member",
				"description": "Member role",
				"permissions": []string{},
			},
			"precedence": 2,
			"created_at": 1640995200000,
			"updated_at": 1640995200000,
		},
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/scim_directories/scim_dir_123/group_role_mappings",
		},
	}

	client := NewClient(config)
	mappings, err := client.List(context.Background(), "scim_dir_123")
	require.NoError(t, err)
	assert.Len(t, mappings, 2)
	assert.Equal(t, "scim_grp_role_123", mappings[0].ID)
	assert.Equal(t, "Admins", mappings[0].SCIMGroupDisplayName)
	assert.Equal(t, "role_admin", mappings[0].Role.ID)
	assert.Equal(t, "org:admin", mappings[0].Role.Key)
	assert.Equal(t, 1, mappings[0].Precedence)
	assert.Equal(t, "scim_grp_role_456", mappings[1].ID)
	assert.Equal(t, "Members", mappings[1].SCIMGroupDisplayName)
	assert.Equal(t, 2, mappings[1].Precedence)
}

func TestListGroups(t *testing.T) {
	t.Parallel()

	response := []map[string]interface{}{
		{
			"object":       "scim_group",
			"id":           "group_123",
			"display_name": "Admins",
			"updated_at":   1640995200000,
		},
		{
			"object":       "scim_group",
			"id":           "group_456",
			"display_name": "Members",
			"updated_at":   1640995200000,
		},
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/scim_directories/scim_dir_123/groups",
		},
	}

	client := NewClient(config)
	groups, err := client.ListGroups(context.Background(), "scim_dir_123")
	require.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "group_123", groups[0].ID)
	assert.Equal(t, "Admins", groups[0].DisplayName)
	assert.Equal(t, "group_456", groups[1].ID)
	assert.Equal(t, "Members", groups[1].DisplayName)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":                  "scim_group_role_mapping",
		"id":                      "scim_grp_role_123",
		"scim_directory_id":       "scim_dir_123",
		"scim_group_id":           "group_123",
		"scim_group_display_name": "Admins",
		"role": map[string]interface{}{
			"object":      "role",
			"id":          "role_admin",
			"name":        "Admin",
			"key":         "org:admin",
			"permissions": []string{},
		},
		"precedence": 1,
		"created_at": 1640995200000,
		"updated_at": 1640995200000,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"role_id":"role_admin","scim_group_id":"group_123"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/scim_directories/scim_dir_123/group_role_mappings",
		},
	}

	client := NewClient(config)
	params := &CreateParams{
		SCIMGroupID: "group_123",
		RoleID:      "role_admin",
	}

	mapping, err := client.Create(context.Background(), "scim_dir_123", params)
	require.NoError(t, err)
	assert.Equal(t, "scim_grp_role_123", mapping.ID)
	assert.Equal(t, "group_123", mapping.SCIMGroupID)
	assert.Equal(t, "Admins", mapping.SCIMGroupDisplayName)
	assert.Equal(t, "role_admin", mapping.Role.ID)
	assert.Equal(t, 1, mapping.Precedence)
}

func TestBulkUpdate(t *testing.T) {
	t.Parallel()

	response := []map[string]interface{}{
		{
			"object":                  "scim_group_role_mapping",
			"id":                      "scim_grp_role_456",
			"scim_directory_id":       "scim_dir_123",
			"scim_group_id":           "group_456",
			"scim_group_display_name": "Members",
			"role": map[string]interface{}{
				"object":      "role",
				"id":          "role_member",
				"name":        "Member",
				"key":         "org:member",
				"permissions": []string{},
			},
			"precedence": 1,
			"created_at": 1640995200000,
			"updated_at": 1640995200000,
		},
		{
			"object":                  "scim_group_role_mapping",
			"id":                      "scim_grp_role_123",
			"scim_directory_id":       "scim_dir_123",
			"scim_group_id":           "group_123",
			"scim_group_display_name": "Admins",
			"role": map[string]interface{}{
				"object":      "role",
				"id":          "role_admin",
				"name":        "Admin",
				"key":         "org:admin",
				"permissions": []string{},
			},
			"precedence": 2,
			"created_at": 1640995200000,
			"updated_at": 1640995200000,
		},
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"mappings":[{"id":"scim_grp_role_456"},{"id":"scim_grp_role_123"}]}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPatch,
			Path:   "/v1/scim_directories/scim_dir_123/group_role_mappings",
		},
	}

	client := NewClient(config)
	params := &BulkUpdateParams{
		Mappings: []MappingUpdate{
			{ID: "scim_grp_role_456"},
			{ID: "scim_grp_role_123"},
		},
	}

	mappings, err := client.BulkUpdate(context.Background(), "scim_dir_123", params)
	require.NoError(t, err)
	assert.Len(t, mappings, 2)
	// First mapping should now have precedence 1 (was reordered to first position)
	assert.Equal(t, "scim_grp_role_456", mappings[0].ID)
	assert.Equal(t, 1, mappings[0].Precedence)
	// Second mapping should now have precedence 2
	assert.Equal(t, "scim_grp_role_123", mappings[1].ID)
	assert.Equal(t, 2, mappings[1].Precedence)
}

func TestBulkUpdateWithRoleChange(t *testing.T) {
	t.Parallel()

	newRoleID := "role_new"
	response := []map[string]interface{}{
		{
			"object":                  "scim_group_role_mapping",
			"id":                      "scim_grp_role_123",
			"scim_directory_id":       "scim_dir_123",
			"scim_group_id":           "group_123",
			"scim_group_display_name": "Admins",
			"role": map[string]interface{}{
				"object":      "role",
				"id":          newRoleID,
				"name":        "New Role",
				"key":         "org:new",
				"permissions": []string{},
			},
			"precedence": 1,
			"created_at": 1640995200000,
			"updated_at": 1640995200000,
		},
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"mappings":[{"id":"scim_grp_role_123","role_id":"role_new"}]}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPatch,
			Path:   "/v1/scim_directories/scim_dir_123/group_role_mappings",
		},
	}

	client := NewClient(config)
	params := &BulkUpdateParams{
		Mappings: []MappingUpdate{
			{ID: "scim_grp_role_123", RoleID: &newRoleID},
		},
	}

	mappings, err := client.BulkUpdate(context.Background(), "scim_dir_123", params)
	require.NoError(t, err)
	assert.Len(t, mappings, 1)
	assert.Equal(t, newRoleID, mappings[0].Role.ID)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":  "scim_group_role_mapping",
		"id":      "scim_grp_role_123",
		"deleted": true,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodDelete,
			Path:   "/v1/scim_directories/scim_dir_123/group_role_mappings/scim_grp_role_123",
		},
	}

	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), "scim_dir_123", "scim_grp_role_123")
	require.NoError(t, err)
	assert.Equal(t, "scim_grp_role_123", deletedResource.ID)
	assert.True(t, deletedResource.Deleted)
}
