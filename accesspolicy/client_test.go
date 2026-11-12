package accesspolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestAccessPolicyClientCreate(t *testing.T) {
	t.Parallel()
	id := "acpol_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(`{"domains":["example.com"],"authentication":{"connections":["ent_123"],"directories":["dir_123"],"provisioning":"jit","connection_required":true},"status":"draft"}`),
			Out: json.RawMessage(fmt.Sprintf(`{"object":"access_policy","id":"%s","default":false,"organization_id":null,"domains":[{"name":"example.com"}],"status":"draft","authentication":{"connections":[{"id":"ent_123","name":"Okta","active":true}],"directories":[{"id":"dir_123"}],"provisioning":"jit","connection_required":true},"created_at":1700000000000,"updated_at":1700000000000}`,
				id)),
			Method: http.MethodPost,
			Path:   "/v1/access_policies",
		},
	}
	client := NewClient(config)
	policy, err := client.Create(context.Background(), &CreateParams{
		Domains: &[]string{"example.com"},
		Authentication: &AuthenticationParams{
			Connections:        &[]string{"ent_123"},
			Directories:        &[]string{"dir_123"},
			Provisioning:       clerk.String("jit"),
			ConnectionRequired: clerk.Bool(true),
		},
		Status: clerk.String("draft"),
	})
	require.NoError(t, err)
	require.Equal(t, id, policy.ID)
	require.False(t, policy.Default)
	require.Nil(t, policy.OrganizationID)
	require.Equal(t, "example.com", policy.Domains[0].Name)
	require.Equal(t, "draft", policy.Status)
	require.NotNil(t, policy.Authentication)
	require.Equal(t, "ent_123", policy.Authentication.Connections[0].ID)
	require.Equal(t, "Okta", policy.Authentication.Connections[0].Name)
	require.True(t, policy.Authentication.Connections[0].Active)
	require.Equal(t, "dir_123", policy.Authentication.Directories[0].ID)
	require.Equal(t, "jit", policy.Authentication.Provisioning)
	require.True(t, policy.Authentication.ConnectionRequired)
	require.Equal(t, int64(1700000000000), policy.CreatedAt)
}

func TestAccessPolicyClientCreate_Error(t *testing.T) {
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

func TestAccessPolicyClientGet(t *testing.T) {
	t.Parallel()
	id := "acpol_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"access_policy","id":"%s","default":false,"organization_id":"%s","domains":[{"name":"example.com"}],"status":"active","authentication":null}`, id, organizationID)),
			Method: http.MethodGet,
			Path:   "/v1/access_policies/" + id,
		},
	}
	client := NewClient(config)
	policy, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, policy.ID)
	require.NotNil(t, policy.OrganizationID)
	require.Equal(t, organizationID, *policy.OrganizationID)
	require.Equal(t, "active", policy.Status)
	require.Nil(t, policy.Authentication)
}

func TestAccessPolicyClientUpdate(t *testing.T) {
	t.Parallel()
	id := "acpol_123"
	for _, tc := range []struct {
		name   string
		params UpdateParams
		in     string
	}{
		{
			name:   "sends only the set fields",
			params: UpdateParams{Status: clerk.String("active")},
			in:     `{"status":"active"}`,
		},
		{
			name:   "sends the domains to replace",
			params: UpdateParams{Domains: &[]string{"example.org"}},
			in:     `{"domains":["example.org"]}`,
		},
		{
			name:   "sends a partial authentication object",
			params: UpdateParams{Authentication: &AuthenticationParams{ConnectionRequired: clerk.Bool(false)}},
			in:     `{"authentication":{"connection_required":false}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := &clerk.ClientConfig{}
			config.HTTPClient = &http.Client{
				Transport: &clerktest.RoundTripper{
					T:      t,
					In:     json.RawMessage(tc.in),
					Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"active"}`, id)),
					Method: http.MethodPatch,
					Path:   "/v1/access_policies/" + id,
				},
			}
			client := NewClient(config)
			params := tc.params
			policy, err := client.Update(context.Background(), id, &params)
			require.NoError(t, err)
			require.Equal(t, id, policy.ID)
			require.Equal(t, "active", policy.Status)
		})
	}
}

func TestAccessPolicyClientUpdate_Error(t *testing.T) {
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
	_, err := client.Update(context.Background(), "acpol_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestAccessPolicyClientDelete(t *testing.T) {
	t.Parallel()
	id := "acpol_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"access_policy","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/access_policies/" + id,
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, deletedResource.ID)
	require.True(t, deletedResource.Deleted)
}

func TestAccessPolicyClientList(t *testing.T) {
	t.Parallel()
	id := "acpol_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
"data": [
  {"object":"access_policy","id":"%s","default":true,"organization_id":null,"domains":[],"status":"active","authentication":null}
],
"total_count": 1
}`,
				id)),
			Method: http.MethodGet,
			Path:   "/v1/access_policies",
			Query: &url.Values{
				"limit":  []string{"1"},
				"offset": []string{"2"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, id, list.AccessPolicies[0].ID)
	require.True(t, list.AccessPolicies[0].Default)
	require.Empty(t, list.AccessPolicies[0].Domains)
}
