package organizationpermission

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

func TestOrganizationPermissionClient_Create(t *testing.T) {
	permissionID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	name := "Manage Billing"
	key := "org:billing:manage"
	description := "Permission to manage billing settings"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","key":"%s","description":"%s"}`, name, key, description)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"permission","id":"%s","name":"%s","key":"%s","description":"%s","type":"user","created_at":1234567890,"updated_at":1234567890}`, permissionID, name, key, description)),
			Method: http.MethodPost,
			Path:   "/v1/organization_permissions",
		},
	}
	client := NewClient(config)
	permission, err := client.Create(context.Background(), &CreateParams{
		Name:        clerk.String(name),
		Key:         clerk.String(key),
		Description: clerk.String(description),
	})
	require.NoError(t, err)
	assert.Equal(t, permissionID, permission.ID)
	assert.Equal(t, name, permission.Name)
	assert.Equal(t, key, permission.Key)
	assert.Equal(t, description, *permission.Description)
	assert.Equal(t, "user", permission.Type)
}

func TestOrganizationPermissionClient_Get(t *testing.T) {
	permissionID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	name := "Manage Billing"
	key := "org:billing:manage"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"permission","id":"%s","name":"%s","key":"%s","description":"Permission to manage billing","type":"user","created_at":1234567890,"updated_at":1234567890}`, permissionID, name, key)),
			Method: http.MethodGet,
			Path:   "/v1/organization_permissions/" + url.PathEscape(permissionID),
		},
	}
	client := NewClient(config)
	permission, err := client.Get(context.Background(), permissionID)
	require.NoError(t, err)
	assert.Equal(t, permissionID, permission.ID)
	assert.Equal(t, name, permission.Name)
	assert.Equal(t, key, permission.Key)
}

func TestOrganizationPermissionClient_Update(t *testing.T) {
	permissionID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	oldName := "Manage Billing"
	newName := "Manage All Billing"
	key := "org:billing:manage"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, newName)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"permission","id":"%s","name":"%s","key":"%s","description":"Permission to manage billing","type":"user","created_at":1234567890,"updated_at":1234567891}`, permissionID, newName, key)),
			Method: http.MethodPatch,
			Path:   "/v1/organization_permissions/" + url.PathEscape(permissionID),
		},
	}
	client := NewClient(config)
	permission, err := client.Update(context.Background(), permissionID, &UpdateParams{
		Name: clerk.String(newName),
	})
	require.NoError(t, err)
	assert.Equal(t, permissionID, permission.ID)
	assert.Equal(t, newName, permission.Name)
	assert.NotEqual(t, oldName, permission.Name)
}

func TestOrganizationPermissionClient_Delete(t *testing.T) {
	permissionID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKe"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"permission","id":"%s","deleted":true}`, permissionID)),
			Method: http.MethodDelete,
			Path:   "/v1/organization_permissions/" + url.PathEscape(permissionID),
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), permissionID)
	require.NoError(t, err)
	assert.Equal(t, permissionID, deletedResource.ID)
	assert.True(t, deletedResource.Deleted)
}

func TestOrganizationPermissionClient_List(t *testing.T) {
	perm1ID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKe"
	perm2ID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKf"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
				"data": [
					{"object":"permission","id":"%s","name":"Manage Billing","key":"org:billing:manage","description":"Billing permission","type":"user","created_at":1234567890,"updated_at":1234567890},
					{"object":"permission","id":"%s","name":"Manage Settings","key":"org:settings:manage","description":"Settings permission","type":"user","created_at":1234567890,"updated_at":1234567890}
				],
				"total_count": 2
			}`, perm1ID, perm2ID)),
			Method: http.MethodGet,
			Path:   "/v1/organization_permissions",
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	assert.Len(t, list.OrganizationPermissions, 2)
	assert.Equal(t, int64(2), list.TotalCount)
	assert.Equal(t, perm1ID, list.OrganizationPermissions[0].ID)
	assert.Equal(t, perm2ID, list.OrganizationPermissions[1].ID)
}

func TestListParams_ToQuery(t *testing.T) {
	t.Parallel()

	t.Run("with all parameters", func(t *testing.T) {
		params := &ListParams{
			ListParams: clerk.ListParams{
				Limit:  clerk.Int64(10),
				Offset: clerk.Int64(20),
			},
			Query:   clerk.String("billing"),
			OrderBy: clerk.String("created_at"),
		}

		query := params.ToQuery()
		assert.Equal(t, "10", query.Get("limit"))
		assert.Equal(t, "20", query.Get("offset"))
		assert.Equal(t, "billing", query.Get("query"))
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

func TestOrganizationPermissionClient_ListWithQueryParams(t *testing.T) {
	t.Parallel()
	perm1ID := "perm_2b6E7b8FdHPjQKsrrakHLUPOzKe"

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
					{"object":"permission","id":"%s","name":"Manage Billing","key":"org:billing:manage","description":"Billing permission","type":"user","created_at":1234567890,"updated_at":1234567890}
				],
				"total_count": 1
			}`, perm1ID)),
			Method: http.MethodGet,
			Path:   "/v1/organization_permissions",
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
	assert.Len(t, list.OrganizationPermissions, 1)
	assert.Equal(t, int64(1), list.TotalCount)
	assert.Equal(t, perm1ID, list.OrganizationPermissions[0].ID)
}
