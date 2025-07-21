package organizationdomain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
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
	response, err := client.Create(context.Background(), organizationID, &CreateParams{
		Name:           clerk.String(domain),
		EnrollmentMode: clerk.String("automatic_invitation"),
		Verified:       clerk.Bool(verified),
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
	_, err := client.Create(context.Background(), "org_123", &CreateParams{})
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
			In:     json.RawMessage(fmt.Sprintf(`{"verified": %s}`, strconv.FormatBool(verified))),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","verification":{"status": "verified"}}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + organizationID + "/domains/" + id,
		},
	}
	client := NewClient(config)
	domain, err := client.Update(context.Background(), &UpdateParams{
		OrganizationID: organizationID,
		DomainID:       id,
		Verified:       clerk.Bool(verified),
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
		Verified:        &verified,
		EnrollmentModes: &[]string{"automatic_invitation"},
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), organizationID, params)
	require.NoError(t, err)
	require.Equal(t, id, list.OrganizationDomains[0].ID)
	require.Equal(t, organizationID, list.OrganizationDomains[0].OrganizationID)
}

func TestOrganizationDomainClientListFromInstance(t *testing.T) {
	t.Parallel()
	id1 := "orgdm_123"
	id2 := "orgdm_456"
	domain1 := "mydomain.com"
	domain2 := "anotherdomain.com"
	organizationID1 := "org_123"
	organizationID2 := "org_456"
	verified := true
	query := "mydomain.com"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
"data": [
  {"enrollment_mode":"automatic_suggestion","id":"%s","name":"%s","object":"organization_domain","organization_id":"%s","verification":{"status":"verified"}},
  {"enrollment_mode":"automatic_invitation","id":"%s","name":"%s","object":"organization_domain","organization_id":"%s","verification":{"status":"unverified"}}
],
"total_count": 2
}`,
				id1, domain1, organizationID1, id2, domain2, organizationID2)),
			Method: http.MethodGet,
			Path:   "/v1/organization_domains",
			Query: &url.Values{
				"limit":           []string{"10"},
				"offset":          []string{"0"},
				"verified":        []string{"true"},
				"enrollment_mode": []string{"automatic_invitation"},
				"organization_id": []string{"org_123"},
				"order_by":        []string{"-created_at"},
				"query":           []string{"mydomain.com"},
			},
		},
	}
	client := NewClient(config)
	params := &ListAllFromInstanceParams{
		Verified:       &verified,
		EnrollmentMode: clerk.String("automatic_invitation"),
		OrganizationID: clerk.String("org_123"),
		OrderBy:        clerk.String("-created_at"),
		Query:          clerk.String(query),
	}
	params.Limit = clerk.Int64(10)
	params.Offset = clerk.Int64(0)
	list, err := client.ListAllFromInstance(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, int64(2), list.TotalCount)
	require.Equal(t, 2, len(list.OrganizationDomains))
	require.Equal(t, id1, list.OrganizationDomains[0].ID)
	require.Equal(t, organizationID1, list.OrganizationDomains[0].OrganizationID)
	require.Equal(t, id2, list.OrganizationDomains[1].ID)
	require.Equal(t, organizationID2, list.OrganizationDomains[1].OrganizationID)
}
