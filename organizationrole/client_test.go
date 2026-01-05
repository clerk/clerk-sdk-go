package organizationrole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationRoleClient_Create(t *testing.T) {
	roleID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	name := "Billing Admin"
	key := "org:billing_admin"
	description := "Role for billing administrators"
	permissionID := "perm_123"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","key":"%s","description":"%s","permissions":["%s"],"include_in_initial_role_set":false}`, name, key, description, permissionID)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role","id":"%s","name":"%s","key":"%s","description":"%s","permissions":[{"object":"permission","id":"%s","name":"Manage Billing","key":"org:billing:manage","type":"user","created_at":1234567890,"updated_at":1234567890}],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567890}`, roleID, name, key, description, permissionID)),
			Method: http.MethodPost,
			Path:   "/v1/organization_roles",
		},
	}
	client := NewClient(config)
	role, err := client.Create(context.Background(), &CreateParams{
		Name:                    clerk.String(name),
		Key:                     clerk.String(key),
		Description:             clerk.String(description),
		Permissions:             &[]string{permissionID},
		IncludeInInitialRoleSet: clerk.Bool(false),
	})
	require.NoError(t, err)
	assert.Equal(t, roleID, role.ID)
	assert.Equal(t, name, role.Name)
	assert.Equal(t, key, role.Key)
	assert.Equal(t, description, *role.Description)
	assert.Len(t, role.Permissions, 1)
	assert.Equal(t, permissionID, role.Permissions[0].ID)
	assert.False(t, role.IsCreatorEligible)
}

func TestOrganizationRoleClient_Get(t *testing.T) {
	roleID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	name := "Billing Admin"
	key := "org:billing_admin"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role","id":"%s","name":"%s","key":"%s","description":"Role for billing administrators","permissions":[],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567890}`, roleID, name, key)),
			Method: http.MethodGet,
			Path:   "/v1/organization_roles/" + url.PathEscape(roleID),
		},
	}
	client := NewClient(config)
	role, err := client.Get(context.Background(), roleID)
	require.NoError(t, err)
	assert.Equal(t, roleID, role.ID)
	assert.Equal(t, name, role.Name)
	assert.Equal(t, key, role.Key)
}

func TestOrganizationRoleClient_Update(t *testing.T) {
	roleID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	oldName := "Billing Admin"
	newName := "Senior Billing Admin"
	key := "org:billing_admin"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, newName)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role","id":"%s","name":"%s","key":"%s","description":"Role for billing administrators","permissions":[],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567891}`, roleID, newName, key)),
			Method: http.MethodPatch,
			Path:   "/v1/organization_roles/" + url.PathEscape(roleID),
		},
	}
	client := NewClient(config)
	role, err := client.Update(context.Background(), roleID, &UpdateParams{
		Name: clerk.String(newName),
	})
	require.NoError(t, err)
	assert.Equal(t, roleID, role.ID)
	assert.Equal(t, newName, role.Name)
	assert.NotEqual(t, oldName, role.Name)
}

func TestOrganizationRoleClient_Delete(t *testing.T) {
	roleID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role","id":"%s","deleted":true}`, roleID)),
			Method: http.MethodDelete,
			Path:   "/v1/organization_roles/" + url.PathEscape(roleID),
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), roleID)
	require.NoError(t, err)
	assert.Equal(t, roleID, deletedResource.ID)
	assert.True(t, deletedResource.Deleted)
}

func TestOrganizationRoleClient_List(t *testing.T) {
	role1ID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	role2ID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKf"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
				"data": [
					{"object":"role","id":"%s","name":"Billing Admin","key":"org:billing_admin","description":"Billing role","permissions":[],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567890},
					{"object":"role","id":"%s","name":"Support Admin","key":"org:support_admin","description":"Support role","permissions":[],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567890}
				],
				"total_count": 2
			}`, role1ID, role2ID)),
			Method: http.MethodGet,
			Path:   "/v1/organization_roles",
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	assert.Len(t, list.OrganizationRoles, 2)
	assert.Equal(t, int64(2), list.TotalCount)
	assert.Equal(t, role1ID, list.OrganizationRoles[0].ID)
	assert.Equal(t, role2ID, list.OrganizationRoles[1].ID)
}

func TestOrganizationRoleClient_AssignPermission(t *testing.T) {
	roleID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	permissionID := "perm_456"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role","id":"%s","name":"Billing Admin","key":"org:billing_admin","description":"Billing role","permissions":[{"object":"permission","id":"%s","name":"Manage Settings","key":"org:settings:manage","type":"user","created_at":1234567890,"updated_at":1234567890}],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567891}`, roleID, permissionID)),
			Method: http.MethodPost,
			Path:   "/v1/organization_roles/" + url.PathEscape(roleID) + "/permissions/" + url.PathEscape(permissionID),
		},
	}
	client := NewClient(config)
	role, err := client.AssignPermission(context.Background(), roleID, permissionID)
	require.NoError(t, err)
	assert.Equal(t, roleID, role.ID)
	assert.Len(t, role.Permissions, 1)
	assert.Equal(t, permissionID, role.Permissions[0].ID)
}

func TestOrganizationRoleClient_RemovePermission(t *testing.T) {
	roleID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	permissionID := "perm_456"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"role","id":"%s","name":"Billing Admin","key":"org:billing_admin","description":"Billing role","permissions":[],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567891}`, roleID)),
			Method: http.MethodDelete,
			Path:   "/v1/organization_roles/" + url.PathEscape(roleID) + "/permissions/" + url.PathEscape(permissionID),
		},
	}
	client := NewClient(config)
	role, err := client.RemovePermission(context.Background(), roleID, permissionID)
	require.NoError(t, err)
	assert.Equal(t, roleID, role.ID)
	assert.Len(t, role.Permissions, 0)
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

func TestOrganizationRoleClient_ListWithQueryParams(t *testing.T) {
	t.Parallel()
	role1ID := "role_2b6E7b8FdHPjQKsrrakHLUPOzKe"

	expectedQuery := url.Values{}
	expectedQuery.Set("query", "billing")
	expectedQuery.Set("order_by", "created_at")
	expectedQuery.Set("limit", "10")

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
				"data": [
					{"object":"role","id":"%s","name":"Billing Admin","key":"org:billing_admin","description":"Billing role","permissions":[],"is_creator_eligible":false,"created_at":1234567890,"updated_at":1234567890}
				],
				"total_count": 1
			}`, role1ID)),
			Method: http.MethodGet,
			Path:   "/v1/organization_roles",
			Query:  &expectedQuery,
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{
		ListParams: clerk.ListParams{
			Limit: clerk.Int64(10),
		},
		Query:   clerk.String("billing"),
		OrderBy: clerk.String("created_at"),
	})
	require.NoError(t, err)
	assert.Len(t, list.OrganizationRoles, 1)
	assert.Equal(t, int64(1), list.TotalCount)
	assert.Equal(t, role1ID, list.OrganizationRoles[0].ID)
}
