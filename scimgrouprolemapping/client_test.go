package scimgrouprolemapping

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

func TestList(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{
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
		},
		"total_count": 2,
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
	listResp, err := client.List(context.Background(), "scim_dir_123")
	require.NoError(t, err)
	assert.Len(t, listResp.Data, 2)
	assert.Equal(t, int64(2), listResp.TotalCount)
	assert.Equal(t, "scim_grp_role_123", listResp.Data[0].ID)
	assert.Equal(t, "Admins", listResp.Data[0].SCIMGroupDisplayName)
	assert.Equal(t, "role_admin", listResp.Data[0].Role.ID)
	assert.Equal(t, "org:admin", listResp.Data[0].Role.Key)
	assert.Equal(t, 1, listResp.Data[0].Precedence)
	assert.Equal(t, "scim_grp_role_456", listResp.Data[1].ID)
	assert.Equal(t, "Members", listResp.Data[1].SCIMGroupDisplayName)
	assert.Equal(t, 2, listResp.Data[1].Precedence)
}

func TestListGroups(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{
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
		},
		"cursor": map[string]interface{}{
			"starting_after": nil,
			"ending_before":  nil,
			"has_next_page":  false,
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
	groups, err := client.ListGroups(context.Background(), "scim_dir_123", nil)
	require.NoError(t, err)
	assert.Len(t, groups.Data, 2)
	assert.Equal(t, "group_123", groups.Data[0].ID)
	assert.Equal(t, "Admins", groups.Data[0].DisplayName)
	assert.Equal(t, "group_456", groups.Data[1].ID)
	assert.Equal(t, "Members", groups.Data[1].DisplayName)
	assert.False(t, groups.Cursor.HasNextPage)
}

func TestListGroupsParams_ToQuery(t *testing.T) {
	t.Parallel()

	startingAfter := "group_123"
	endingBefore := "group_456"
	limit := 25

	t.Run("nil receiver produces empty query", func(t *testing.T) {
		var p *ListGroupsParams
		assert.Equal(t, url.Values{}, p.ToQuery())
	})

	t.Run("empty params produces empty query", func(t *testing.T) {
		p := &ListGroupsParams{}
		assert.Equal(t, url.Values{}, p.ToQuery())
	})

	t.Run("starting_after only", func(t *testing.T) {
		p := &ListGroupsParams{StartingAfter: &startingAfter}
		q := p.ToQuery()
		assert.Equal(t, "group_123", q.Get("starting_after"))
		assert.Empty(t, q.Get("ending_before"))
		assert.Empty(t, q.Get("limit"))
	})

	t.Run("ending_before only", func(t *testing.T) {
		p := &ListGroupsParams{EndingBefore: &endingBefore}
		q := p.ToQuery()
		assert.Equal(t, "group_456", q.Get("ending_before"))
		assert.Empty(t, q.Get("starting_after"))
		assert.Empty(t, q.Get("limit"))
	})

	t.Run("limit only", func(t *testing.T) {
		p := &ListGroupsParams{Limit: &limit}
		q := p.ToQuery()
		assert.Equal(t, "25", q.Get("limit"))
		assert.Empty(t, q.Get("starting_after"))
		assert.Empty(t, q.Get("ending_before"))
	})

	t.Run("all params combined", func(t *testing.T) {
		p := &ListGroupsParams{
			StartingAfter: &startingAfter,
			EndingBefore:  &endingBefore,
			Limit:         &limit,
		}
		q := p.ToQuery()
		assert.Equal(t, "group_123", q.Get("starting_after"))
		assert.Equal(t, "group_456", q.Get("ending_before"))
		assert.Equal(t, "25", q.Get("limit"))
	})
}

func TestListGroupsWithPaginationParams(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"object":       "scim_group",
				"id":           "group_456",
				"display_name": "Members",
				"updated_at":   1640995200000,
			},
		},
		"cursor": map[string]interface{}{
			"starting_after": "group_456",
			"ending_before":  nil,
			"has_next_page":  true,
		},
	}

	responseJSON, _ := json.Marshal(response)

	expectedQuery := url.Values{
		"starting_after": []string{"group_123"},
		"limit":          []string{"1"},
	}

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodGet,
			Path:   "/v1/scim_directories/scim_dir_123/groups",
			Query:  &expectedQuery,
		},
	}

	client := NewClient(config)
	startingAfter := "group_123"
	limit := 1
	params := &ListGroupsParams{
		StartingAfter: &startingAfter,
		Limit:         &limit,
	}

	groups, err := client.ListGroups(context.Background(), "scim_dir_123", params)
	require.NoError(t, err)
	assert.Len(t, groups.Data, 1)
	assert.Equal(t, "group_456", groups.Data[0].ID)
	assert.True(t, groups.Cursor.HasNextPage)
	require.NotNil(t, groups.Cursor.StartingAfter)
	assert.Equal(t, "group_456", *groups.Cursor.StartingAfter)
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

	response := map[string]interface{}{
		"data": []map[string]interface{}{
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
		},
		"total_count": 2,
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

	listResp, err := client.BulkUpdate(context.Background(), "scim_dir_123", params)
	require.NoError(t, err)
	assert.Len(t, listResp.Data, 2)
	// First mapping should now have precedence 1 (was reordered to first position)
	assert.Equal(t, "scim_grp_role_456", listResp.Data[0].ID)
	assert.Equal(t, 1, listResp.Data[0].Precedence)
	// Second mapping should now have precedence 2
	assert.Equal(t, "scim_grp_role_123", listResp.Data[1].ID)
	assert.Equal(t, 2, listResp.Data[1].Precedence)
}

func TestBulkUpdateWithRoleChange(t *testing.T) {
	t.Parallel()

	newRoleID := "role_new"
	response := map[string]interface{}{
		"data": []map[string]interface{}{
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
		},
		"total_count": 1,
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

	listResp, err := client.BulkUpdate(context.Background(), "scim_dir_123", params)
	require.NoError(t, err)
	assert.Len(t, listResp.Data, 1)
	assert.Equal(t, newRoleID, listResp.Data[0].Role.ID)
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
