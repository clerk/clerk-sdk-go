package organizationinvitation

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

func TestOrganizationInvitationClientCreate(t *testing.T) {
	t.Parallel()
	id := "orginv_123"
	organizationID := "org_123"
	emailAddress := "foo@bar.com"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"email_address":"%s"}`, emailAddress)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","email_address":"%s","organization_id":"%s"}`, id, emailAddress, organizationID)),
			Method: http.MethodPost,
			Path:   "/v1/organizations/" + organizationID + "/invitations",
		},
	}
	client := NewClient(config)
	invitation, err := client.Create(context.Background(), &CreateParams{
		OrganizationID: organizationID,
		EmailAddress:   clerk.String(emailAddress),
	})
	require.NoError(t, err)
	require.Equal(t, id, invitation.ID)
	require.Equal(t, organizationID, invitation.OrganizationID)
	require.Equal(t, emailAddress, invitation.EmailAddress)
}

func TestOrganizationInvitationClientCreate_Error(t *testing.T) {
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

func TestOrganizationInvitationClientList(t *testing.T) {
	t.Parallel()
	organizationID := "org_123"
	id := "orginv_123"
	statuses := []string{"pending", "accepted"}
	limit := int64(10)
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"data":[{"id":"%s","object":"organization_invitation","email_address":"string","role":"string","organization_id":"%s","status":"string","public_metadata":{},"private_metadata":{},"created_at":0,"updated_at":0}],"total_count":1}`, id, organizationID)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + organizationID + "/invitations",
			Query: &url.Values{
				"limit":  []string{fmt.Sprintf("%d", limit)},
				"status": statuses,
			},
		},
	}
	client := NewClient(config)
	response, err := client.List(context.Background(), &ListParams{
		OrganizationID: organizationID,
		ListParams: clerk.ListParams{
			Limit: clerk.Int64(limit),
		},
		Statuses: &statuses,
	})
	require.NoError(t, err)
	require.Len(t, response.OrganizationInvitations, 1)
	require.Equal(t, id, response.OrganizationInvitations[0].ID)
	require.Equal(t, organizationID, response.OrganizationInvitations[0].OrganizationID)
	require.Equal(t, int64(1), response.TotalCount)
}

func TestOrganizationInvitationClientList_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
				"errors":[{
					"code":"list-error-code"
				}],
				"clerk_trace_id":"list-trace-id"
			}`),
		},
	}
	client := NewClient(config)
	_, err := client.List(context.Background(), &ListParams{OrganizationID: "org_123"})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "list-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "list-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationInvitationClientGet(t *testing.T) {
	t.Parallel()
	organizationID := "org_123"
	id := "orginv_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"organization_invitation","email_address":"string","role":"string","organization_id":"%s","status":"string","public_metadata":{},"private_metadata":{},"created_at":0,"updated_at":0}`, id, organizationID)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + organizationID + "/invitations/" + id,
		},
	}
	client := NewClient(config)
	response, err := client.Get(context.Background(), &GetParams{
		OrganizationID: organizationID,
		ID:             id,
	})
	require.NoError(t, err)
	require.Equal(t, id, response.ID)
	require.Equal(t, organizationID, response.OrganizationID)
}

func TestOrganizationInvitationClientGet_Error(t *testing.T) {
	t.Parallel()
	organizationID := "org_123"
	id := "orginv_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
				"errors":[{
					"code":"get-error-code"
				}],
				"clerk_trace_id":"get-trace-id"
			}`),
		},
	}
	client := NewClient(config)
	_, err := client.Get(context.Background(), &GetParams{
		OrganizationID: organizationID,
		ID:             id,
	})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "get-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "get-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationInvitationClientRevoke(t *testing.T) {
	t.Parallel()
	organizationID := "org_123"
	id := "orginv_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"organization_invitation","email_address":"string","role":"string","organization_id":"%s","status":"string","public_metadata":{},"private_metadata":{},"created_at":0,"updated_at":0}`, id, organizationID)),
			Method: http.MethodPost,
			Path:   "/v1/organizations/" + organizationID + "/invitations/" + id + "/revoke",
		},
	}
	client := NewClient(config)
	response, err := client.Revoke(context.Background(), &RevokeParams{
		OrganizationID: organizationID,
		ID:             id,
	})
	require.NoError(t, err)
	require.Equal(t, id, response.ID)
	require.Equal(t, organizationID, response.OrganizationID)
}

func TestOrganizationInvitationClientRevoke_Error(t *testing.T) {
	t.Parallel()
	organizationID := "org_123"
	id := "orginv_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
				"errors":[{
					"code":"revoke-error-code"
				}],
				"clerk_trace_id":"revoke-trace-id"
			}`),
		},
	}
	client := NewClient(config)
	_, err := client.Revoke(context.Background(), &RevokeParams{
		OrganizationID: organizationID,
		ID:             id,
	})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "revoke-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "revoke-error-code", apiErr.Errors[0].Code)
}
