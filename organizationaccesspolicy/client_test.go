package organizationaccesspolicy

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

func TestOrganizationAccessPolicyClientCreate(t *testing.T) {
	t.Parallel()
	id := "orgap_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(`{"domains":["example.com"],"enrollment_mode":"sso","authentication":{"connections":["entconn_123"],"directories":["dir_123"],"provisioning":"jit","connection_required":true},"status":"draft"}`),
			Out: json.RawMessage(fmt.Sprintf(`{"object":"organization_access_policy","id":"%s","organization_id":"%s","default":false,"domains":[{"name":"example.com","affiliation_verified":false,"ownership_verified":true}],"enrollment_mode":"sso","status":"draft","authentication":{"connections":[{"id":"entconn_123","name":"Okta","active":true}],"directories":[{"id":"dir_123"}],"provisioning":null,"connection_required":false},"created_at":1700000000000,"updated_at":1700000000000}`,
				id, organizationID)),
			Method: http.MethodPost,
			Path:   "/v1/organizations/" + organizationID + "/access_policies",
		},
	}
	client := NewClient(config)
	policy, err := client.Create(context.Background(), organizationID, &CreateParams{
		Domains:        &[]string{"example.com"},
		EnrollmentMode: clerk.String("sso"),
		Authentication: &AuthenticationParams{
			Connections:        &[]string{"entconn_123"},
			Directories:        &[]string{"dir_123"},
			Provisioning:       clerk.String("jit"),
			ConnectionRequired: clerk.Bool(true),
		},
		Status: clerk.String("draft"),
	})
	require.NoError(t, err)
	require.Equal(t, id, policy.ID)
	require.Equal(t, organizationID, policy.OrganizationID)
	require.False(t, policy.Default)
	require.Equal(t, "sso", policy.EnrollmentMode)
	require.Equal(t, "example.com", policy.Domains[0].Name)
	require.False(t, policy.Domains[0].AffiliationVerified)
	require.True(t, policy.Domains[0].OwnershipVerified)
	require.Equal(t, "draft", *policy.Status)
	require.NotNil(t, policy.Authentication)
	require.Equal(t, "entconn_123", policy.Authentication.Connections[0].ID)
	require.Equal(t, "Okta", policy.Authentication.Connections[0].Name)
	require.True(t, policy.Authentication.Connections[0].Active)
	require.Equal(t, "dir_123", policy.Authentication.Directories[0].ID)
	require.Nil(t, policy.Authentication.Provisioning)
	require.Equal(t, int64(1700000000000), policy.CreatedAt)
}

func TestOrganizationAccessPolicyClientCreate_Error(t *testing.T) {
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
	_, err := client.Create(context.Background(), "org_123", &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationAccessPolicyClientGet(t *testing.T) {
	t.Parallel()
	id := "orgap_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"organization_access_policy","id":"%s","organization_id":"%s","default":true,"domains":[],"enrollment_mode":"manual_invitation","status":null,"authentication":null}`, id, organizationID)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + organizationID + "/access_policies/" + id,
		},
	}
	client := NewClient(config)
	policy, err := client.Get(context.Background(), &GetParams{
		OrganizationID: organizationID,
		PolicyID:       id,
	})
	require.NoError(t, err)
	require.Equal(t, id, policy.ID)
	require.True(t, policy.Default)
	require.Empty(t, policy.Domains)
	require.Nil(t, policy.Status)
	require.Nil(t, policy.Authentication)
}

func TestOrganizationAccessPolicyClientUpdate(t *testing.T) {
	t.Parallel()
	id := "orgap_123"
	organizationID := "org_123"
	for _, tc := range []struct {
		name   string
		params UpdateParams
		in     string
	}{
		{
			name:   "sends only the set fields",
			params: UpdateParams{EnrollmentMode: clerk.String("automatic_invitation")},
			in:     `{"enrollment_mode":"automatic_invitation"}`,
		},
		{
			name:   "sends an empty domains list to clear them",
			params: UpdateParams{Domains: &[]string{}},
			in:     `{"domains":[]}`,
		},
		{
			name:   "sends a partial authentication object",
			params: UpdateParams{Authentication: &AuthenticationParams{ConnectionRequired: clerk.Bool(false)}},
			in:     `{"authentication":{"connection_required":false}}`,
		},
		{
			name:   "sends the status",
			params: UpdateParams{Status: clerk.String("active")},
			in:     `{"status":"active"}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := &clerk.ClientConfig{}
			config.HTTPClient = &http.Client{
				Transport: &clerktest.RoundTripper{
					T:      t,
					In:     json.RawMessage(tc.in),
					Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","enrollment_mode":"automatic_invitation"}`, id)),
					Method: http.MethodPatch,
					Path:   "/v1/organizations/" + organizationID + "/access_policies/" + id,
				},
			}
			client := NewClient(config)
			params := tc.params
			params.OrganizationID = organizationID
			params.PolicyID = id
			policy, err := client.Update(context.Background(), &params)
			require.NoError(t, err)
			require.Equal(t, id, policy.ID)
			require.Equal(t, "automatic_invitation", policy.EnrollmentMode)
		})
	}
}

func TestOrganizationAccessPolicyClientUpdate_Error(t *testing.T) {
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

func TestOrganizationAccessPolicyClientDelete(t *testing.T) {
	t.Parallel()
	id := "orgap_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"organization_access_policy","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + organizationID + "/access_policies/" + id,
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), &DeleteParams{
		OrganizationID: organizationID,
		PolicyID:       id,
	})
	require.NoError(t, err)
	require.Equal(t, id, deletedResource.ID)
	require.True(t, deletedResource.Deleted)
}

func TestOrganizationAccessPolicyClientList(t *testing.T) {
	t.Parallel()
	id := "orgap_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
"data": [
  {"object":"organization_access_policy","id":"%s","organization_id":"%s","default":false,"domains":[{"name":"example.com","affiliation_verified":true,"ownership_verified":false}],"enrollment_mode":"automatic_suggestion","authentication":null}
],
"total_count": 1
}`,
				id, organizationID)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + organizationID + "/access_policies",
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
	list, err := client.List(context.Background(), organizationID, params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, id, list.OrganizationAccessPolicies[0].ID)
	require.Equal(t, organizationID, list.OrganizationAccessPolicies[0].OrganizationID)
	require.True(t, list.OrganizationAccessPolicies[0].Domains[0].AffiliationVerified)
}

func TestOrganizationAccessPolicyClientVerify(t *testing.T) {
	t.Parallel()
	id := "orgap_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"domain":"example.com","kind":"ownership"}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"organization_access_policy","id":"%s","organization_id":"%s","domains":[{"name":"example.com","affiliation_verified":false,"ownership_verified":true}]}`, id, organizationID)),
			Method: http.MethodPost,
			Path:   "/v1/organizations/" + organizationID + "/access_policies/" + id + "/verify",
		},
	}
	client := NewClient(config)
	policy, err := client.Verify(context.Background(), &VerifyParams{
		OrganizationID: organizationID,
		PolicyID:       id,
		Domain:         clerk.String("example.com"),
		Kind:           clerk.String("ownership"),
	})
	require.NoError(t, err)
	require.Equal(t, id, policy.ID)
	require.True(t, policy.Domains[0].OwnershipVerified)
}

func TestOrganizationAccessPolicyClientVerify_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusNotFound,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"verify-error-code"
	}],
	"clerk_trace_id":"verify-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.Verify(context.Background(), &VerifyParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "verify-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "verify-error-code", apiErr.Errors[0].Code)
}
