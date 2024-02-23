package organizationinvitation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
