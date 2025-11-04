package roleset

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleSetClient_Create(t *testing.T) {
	roleSetID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	name := "Admin Role Set"
	key := "admin-roles"
	description := "Roles for administrative users"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","key":"%s","description":"%s", "roles": ["org:member"], "type": "initial", "default_role_key": "org:member"}`, name, key, description)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role_set","id":"%s","name":"%s","key":"%s","description":"%s","roles":[],"default_role":{"object":"role_set_item","id":"role_123","name":"Member","key":"org:member","description":"Default member role","created_at":1234567890,"updated_at":1234567890},"type":"initial","created_at":1234567890,"updated_at":1234567890}`, roleSetID, name, key, description)),
			Method: http.MethodPost,
			Path:   "/v1/role_sets",
		},
	}
	client := NewClient(config)
	roleSet, err := client.Create(context.Background(), &CreateParams{
		Name:           clerk.String(name),
		Key:            clerk.String(key),
		Description:    clerk.String(description),
		Roles:          &[]string{"org:member"},
		Type:           clerk.String("initial"),
		DefaultRoleKey: clerk.String("org:member"),
	})
	require.NoError(t, err)
	assert.Equal(t, roleSetID, roleSet.ID)
	assert.Equal(t, name, roleSet.Name)
	assert.Equal(t, key, roleSet.Key)
	assert.Equal(t, description, *roleSet.Description)
	assert.NotNil(t, roleSet.DefaultRole)
	assert.Equal(t, "org:member", roleSet.DefaultRole.Key)
	assert.Equal(t, "Member", roleSet.DefaultRole.Name)
}

func TestRoleSetClient_Get(t *testing.T) {
	roleSetKey := "admin-roles"
	roleSetID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role_set","id":"%s","name":"Admin Role Set","key":"%s","description":"Roles for administrative users","roles":[],"default_role":{"object":"role_set_item","id":"role_123","name":"Member","key":"org:member","description":"Default member role","created_at":1234567890,"updated_at":1234567890},"type":"initial","created_at":1234567890,"updated_at":1234567890}`, roleSetID, roleSetKey)),
			Method: http.MethodGet,
			Path:   "/v1/role_sets/" + url.PathEscape(roleSetKey),
		},
	}
	client := NewClient(config)
	roleSet, err := client.Get(context.Background(), roleSetKey)
	require.NoError(t, err)
	assert.Equal(t, roleSetID, roleSet.ID)
	assert.Equal(t, roleSetKey, roleSet.Key)
}

func TestRoleSetClient_Update(t *testing.T) {
	roleSetKey := "admin-roles"
	roleSetID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	newName := "Updated Admin Role Set"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s", "type": "initial", "default_role_key": "org:admin"}`, newName)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role_set","id":"%s","name":"%s","key":"%s","description":"Roles for administrative users","roles":[],"default_role":{"object":"role_set_item","id":"role_456","name":"Admin","key":"org:admin","description":"Default admin role","created_at":1234567890,"updated_at":1234567890},"type":"initial","created_at":1234567890,"updated_at":1234567891}`, roleSetID, newName, roleSetKey)),
			Method: http.MethodPatch,
			Path:   "/v1/role_sets/" + url.PathEscape(roleSetKey),
		},
	}
	client := NewClient(config)
	roleSet, err := client.Update(context.Background(), roleSetKey, &UpdateParams{
		Name:           clerk.String(newName),
		Type:           clerk.String("initial"),
		DefaultRoleKey: clerk.String("org:admin"),
	})
	require.NoError(t, err)
	assert.Equal(t, roleSetID, roleSet.ID)
	assert.Equal(t, newName, roleSet.Name)
}

func TestRoleSetClient_Delete(t *testing.T) {
	roleSetKey := "admin-roles"
	roleSetID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role_set","id":"%s","deleted":true}`, roleSetID)),
			Method: http.MethodDelete,
			Path:   "/v1/role_sets/" + url.PathEscape(roleSetKey),
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), roleSetKey)
	require.NoError(t, err)
	assert.Equal(t, roleSetID, deletedResource.ID)
	assert.True(t, deletedResource.Deleted)
}

func TestRoleSetClient_List(t *testing.T) {
	roleSet1ID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	roleSet2ID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKf"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
				"data": [
					{"object":"role_set","id":"%s","name":"Admin Role Set","key":"admin-roles","description":"Admin roles","type":"initial","roles":[],"default_role":{"object":"role_set_item","id":"role_123","name":"Member","key":"org:member","description":"Default member role","created_at":1234567890,"updated_at":1234567890},"created_at":1234567890,"updated_at":1234567890},
					{"object":"role_set","id":"%s","name":"User Role Set","key":"user-roles","description":"User roles","type":"custom","roles":[],"default_role":{"object":"role_set_item","id":"role_456","name":"User","key":"org:user","description":"Default user role","created_at":1234567890,"updated_at":1234567890},"created_at":1234567890,"updated_at":1234567890}
				],
				"total_count": 2
			}`, roleSet1ID, roleSet2ID)),
			Method: http.MethodGet,
			Path:   "/v1/role_sets",
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	assert.Len(t, list.RoleSets, 2)
	assert.Equal(t, int64(2), list.TotalCount)
	assert.Equal(t, roleSet1ID, list.RoleSets[0].ID)
	assert.Equal(t, roleSet2ID, list.RoleSets[1].ID)
}

func TestRoleSetClient_AddRoles(t *testing.T) {
	roleSetKey := "admin-roles"
	roleSetID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	roleKeys := []string{"role:admin", "role:manager"}

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"role_keys":["role:admin","role:manager"], "default_role_key": "role:admin"}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role_set","id":"%s","name":"Admin Role Set","type":"initial","key":"%s","description":"Admin roles","roles":[{"object":"role_set_item","id":"role_1","name":"Admin","key":"role:admin","description":"Admin role","created_at":1234567890,"updated_at":1234567890}],"default_role":{"object":"role_set_item","id":"role_123","name":"Member","key":"org:member","description":"Default member role","created_at":1234567890,"updated_at":1234567890},"created_at":1234567890,"updated_at":1234567891}`, roleSetID, roleSetKey)),
			Method: http.MethodPost,
			Path:   "/v1/role_sets/" + url.PathEscape(roleSetKey) + "/roles",
		},
	}
	client := NewClient(config)
	roleSet, err := client.AddRoles(context.Background(), roleSetKey, &AddRolesParams{
		RoleKeys:       roleKeys,
		DefaultRoleKey: clerk.String("role:admin"),
	})
	require.NoError(t, err)
	assert.Equal(t, roleSetID, roleSet.ID)
	assert.Len(t, roleSet.Roles, 1)
	assert.Equal(t, "role:admin", roleSet.Roles[0].Key)
}

func TestRoleSetClient_RemoveRoles(t *testing.T) {
	roleSetKey := "admin-roles"
	roleSetID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	roleKeys := []string{"role:admin"}

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"role_keys":["role:admin"]}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role_set","id":"%s","name":"Admin Role Set","type":"initial","key":"%s","description":"Admin roles","roles":[],"default_role":{"object":"role_set_item","id":"role_123","name":"Member","key":"org:member","description":"Default member role","created_at":1234567890,"updated_at":1234567890},"created_at":1234567890,"updated_at":1234567891}`, roleSetID, roleSetKey)),
			Method: http.MethodDelete,
			Path:   "/v1/role_sets/" + url.PathEscape(roleSetKey) + "/roles",
		},
	}
	client := NewClient(config)
	roleSet, err := client.RemoveRoles(context.Background(), roleSetKey, &RemoveRolesParams{
		RoleKeys: roleKeys,
	})
	require.NoError(t, err)
	assert.Equal(t, roleSetID, roleSet.ID)
	assert.Len(t, roleSet.Roles, 0)
}

func TestListParams_ToQuery(t *testing.T) {
	t.Parallel()

	t.Run("with all parameters", func(t *testing.T) {
		params := &ListParams{
			ListParams: clerk.ListParams{
				Limit:  clerk.Int64(10),
				Offset: clerk.Int64(20),
			},
			Query:   clerk.String("admin"),
			OrderBy: clerk.String("created_at"),
		}

		query := params.ToQuery()
		assert.Equal(t, "10", query.Get("limit"))
		assert.Equal(t, "20", query.Get("offset"))
		assert.Equal(t, "admin", query.Get("query"))
		assert.Equal(t, "created_at", query.Get("order_by"))
	})

	t.Run("with only base list parameters", func(t *testing.T) {
		params := &ListParams{
			ListParams: clerk.ListParams{
				Limit: clerk.Int64(5),
			},
		}

		query := params.ToQuery()
		assert.Equal(t, "5", query.Get("limit"))
		assert.Empty(t, query.Get("query"))
		assert.Empty(t, query.Get("order_by"))
	})

	t.Run("with only new parameters", func(t *testing.T) {
		params := &ListParams{
			Query:   clerk.String("test query"),
			OrderBy: clerk.String("name"),
		}

		query := params.ToQuery()
		assert.Equal(t, "test query", query.Get("query"))
		assert.Equal(t, "name", query.Get("order_by"))
		assert.Empty(t, query.Get("limit"))
		assert.Empty(t, query.Get("offset"))
	})

	t.Run("with nil parameters", func(t *testing.T) {
		params := &ListParams{}

		query := params.ToQuery()
		assert.Empty(t, query.Get("query"))
		assert.Empty(t, query.Get("order_by"))
		assert.Empty(t, query.Get("limit"))
		assert.Empty(t, query.Get("offset"))
	})
}

func TestRoleSetClient_ListWithQueryParams(t *testing.T) {
	t.Parallel()
	roleSet1ID := "role_set_2b6E7b8FdHPjQKsrrakHLUPOzKe"

	expectedQuery := url.Values{}
	expectedQuery.Set("query", "admin")
	expectedQuery.Set("order_by", "created_at")
	expectedQuery.Set("limit", "10")

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
				"data": [
					{"object":"role_set","id":"%s","name":"Admin Role Set","key":"admin-roles","description":"Admin roles","type":"initial","roles":[],"default_role":{"object":"role_set_item","id":"role_123","name":"Member","key":"org:member","description":"Default member role","created_at":1234567890,"updated_at":1234567890},"created_at":1234567890,"updated_at":1234567890}
				],
				"total_count": 1
			}`, roleSet1ID)),
			Method: http.MethodGet,
			Path:   "/v1/role_sets",
			Query:  &expectedQuery,
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{
		ListParams: clerk.ListParams{
			Limit: clerk.Int64(10),
		},
		Query:   clerk.String("admin"),
		OrderBy: clerk.String("created_at"),
	})
	require.NoError(t, err)
	assert.Len(t, list.RoleSets, 1)
	assert.Equal(t, int64(1), list.TotalCount)
	assert.Equal(t, roleSet1ID, list.RoleSets[0].ID)
}
