package organizationdomain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestOrganizationDomainClientCreate(t *testing.T) {
	t.Parallel()
	id := "orgdm_123"
	organizationID := "org_123"
	domain := "mydomain.com"
	verified := false
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(fmt.Sprintf(`{"name": "%s", "enrollment_mode": "automatic_invitation", "verified": %s}`, domain, strconv.FormatBool(verified))),
			Out: json.RawMessage(fmt.Sprintf(`{"enrollment_mode":"automatic_invitation","id":"%s","name":"%s","object":"organization_domain","organization_id":"%s","verification":{"status":"unverified"}}`,
				id, domain, organizationID)),
			Method: http.MethodPost,
			Path:   "/v1/organizations/" + organizationID + "/domains",
		},
	}
	client := NewClient(config)
	response, err := client.Create(context.Background(), &CreateParams{
		OrganizationID: organizationID,
		Name:           domain,
		EnrollmentMode: "automatic_invitation",
		Verified:       &verified,
	})
	require.NoError(t, err)
	require.Equal(t, id, response.ID)
	require.Equal(t, domain, response.Name)
	require.Equal(t, "automatic_invitation", response.EnrollmentMode)
	require.Equal(t, "unverified", response.Verification.Status)
}

func TestOrganizationDomainClientCreate_Error(t *testing.T) {
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

func TestOrganizationDomainClientUpdate(t *testing.T) {
	t.Parallel()
	id := "orgdm_123"
	organizationID := "org_123"
	verified := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"verified": %s, "enrollment_mode": null}`, strconv.FormatBool(verified))),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","verification":{"status": "verified"}}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + organizationID + "/domains/" + id,
		},
	}
	client := NewClient(config)
	domain, err := client.Update(context.Background(), &UpdateParams{
		OrganizationID: organizationID,
		DomainID:       id,
		Verified:       &verified,
		EnrollmentMode: nil,
	})
	require.NoError(t, err)
	require.Equal(t, id, domain.ID)
	require.Equal(t, "verified", domain.Verification.Status)
}

func TestOrganizationDomainClientUpdate_Error(t *testing.T) {
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

func TestOrganizationDomainClientDelete(t *testing.T) {
	t.Parallel()
	id := "orgdm_123"
	organizationID := "org_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"organization_domain"}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + organizationID + "/domains/" + id,
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), &DeleteParams{
		OrganizationID: organizationID,
		DomainID:       id,
	})
	require.NoError(t, err)
	require.Equal(t, id, deletedResource.ID)
}

func TestOrganizationDomainClientList(t *testing.T) {
	t.Parallel()
	id := "orgdm_123"
	domain := "mydomain.com"
	organizationID := "org_123"
	verified := true

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
"data": [
  {"enrollment_mode":"automatic_suggestion","id":"%s","name":"%s","object":"organization_domain","organization_id":"%s","verification":{"status":"unverified"}}
],
"total_count": 1
}`,
				id, domain, organizationID)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + organizationID + "/domains",
			Query: &url.Values{
				"limit":           []string{"1"},
				"offset":          []string{"2"},
				"verified":        []string{"true"},
				"enrollment_mode": []string{"automatic_invitation"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		OrganizationID:  organizationID,
		Verified:        &verified,
		EnrollmentModes: []string{"automatic_invitation"},
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, id, list.OrganizationDomains[0].ID)
	require.Equal(t, organizationID, list.OrganizationDomains[0].OrganizationID)
}
