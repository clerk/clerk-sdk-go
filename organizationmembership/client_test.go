package organizationmembership

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestOrganizationMembershipClientCreate(t *testing.T) {
	t.Parallel()
	id := "orgmem_123"
	organizationID := "org_123"
	userID := "user_123"
	role := "admin"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(fmt.Sprintf(`{"user_id":"%s","role":"%s"}`, userID, role)),
			Out: json.RawMessage(fmt.Sprintf(`{
"id":"%s",
"role":"%s",
"organization":{"id":"%s"},
"public_user_data":{"user_id":"%s"}
}`,
				id, role, organizationID, userID)),
			Method: http.MethodPost,
			Path:   "/v1/organizations/" + organizationID + "/memberships",
		},
	}
	client := NewClient(config)
	membership, err := client.Create(context.Background(), &CreateParams{
		UserID:         clerk.String(userID),
		Role:           clerk.String(role),
		OrganizationID: organizationID,
	})
	require.NoError(t, err)
	require.Equal(t, id, membership.ID)
	require.Equal(t, role, membership.Role)
	require.Equal(t, organizationID, membership.Organization.ID)
	require.Equal(t, userID, membership.PublicUserData.UserID)
}

func TestOrganizationMembershipClientCreate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"create-error-code"
	}],
	"clerk_trace_id":"create-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationMembershipClientUpdate(t *testing.T) {
	t.Parallel()
	id := "orgmem_123"
	organizationID := "org_123"
	userID := "user_123"
	role := "admin"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(fmt.Sprintf(`{"role":"%s"}`, role)),
			Out: json.RawMessage(fmt.Sprintf(`{
"id":"%s",
"role":"%s",
"organization":{"id":"%s"},
"public_user_data":{"user_id":"%s"}
}`,
				id, role, organizationID, userID)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + organizationID + "/memberships/" + userID,
		},
	}
	client := NewClient(config)
	membership, err := client.Update(context.Background(), &UpdateParams{
		Role:           clerk.String(role),
		OrganizationID: organizationID,
		UserID:         userID,
	})
	require.NoError(t, err)
	require.Equal(t, id, membership.ID)
	require.Equal(t, role, membership.Role)
	require.Equal(t, organizationID, membership.Organization.ID)
	require.Equal(t, userID, membership.PublicUserData.UserID)
}

func TestOrganizationMembershipClientUpdate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"update-error-code"
	}],
	"clerk_trace_id":"update-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.Update(context.Background(), &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationMembershipClientDelete(t *testing.T) {
	t.Parallel()
	id := "orgmem_123"
	organizationID := "org_123"
	userID := "user_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
"id":"%s",
"organization":{"id":"%s"},
"public_user_data":{"user_id":"%s"}
}`,
				id, organizationID, userID)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + organizationID + "/memberships/" + userID,
		},
	}
	client := NewClient(config)
	membership, err := client.Delete(context.Background(), &DeleteParams{
		UserID:         userID,
		OrganizationID: organizationID,
	})
	require.NoError(t, err)
	require.Equal(t, id, membership.ID)
	require.Equal(t, organizationID, membership.Organization.ID)
	require.Equal(t, userID, membership.PublicUserData.UserID)
}

func TestOrganizationMembershipClientList(t *testing.T) {
	t.Parallel()
	id := "orgmem_123"
	organizationID := "org_123"
	userID := "user_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
"data": [{
	"id":"%s",
	"organization":{"id":"%s"},
	"public_user_data":{"user_id":"%s"},
	"role": "string",
	"role_name": "string"
}],
"total_count": 1
}`,
				id, organizationID, userID)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + organizationID + "/memberships",
			Query: &url.Values{
				"limit":    []string{"1"},
				"offset":   []string{"2"},
				"role":     []string{"admin", "member"},
				"order_by": []string{"-created_at"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		OrganizationID: organizationID,
		OrderBy:        clerk.String("-created_at"),
		Roles:          []string{"admin", "member"},
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, id, list.OrganizationMemberships[0].ID)
	require.Equal(t, organizationID, list.OrganizationMemberships[0].Organization.ID)
	require.Equal(t, userID, list.OrganizationMemberships[0].PublicUserData.UserID)
	require.Equal(t, "string", list.OrganizationMemberships[0].RoleName)
	require.Equal(t, "string", list.OrganizationMemberships[0].Role)
}
