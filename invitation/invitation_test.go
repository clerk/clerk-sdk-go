package invitation

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

func TestInvitationList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [{"id":"inv_123","email_address":"foo@bar.com"}],
	"total_count": 1
}`),
				Path:   "/v1/invitations",
				Method: http.MethodGet,
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Invitations))
	require.Equal(t, "inv_123", list.Invitations[0].ID)
	require.Equal(t, "foo@bar.com", list.Invitations[0].EmailAddress)
}

func TestInvitationCreate(t *testing.T) {
	emailAddress := "foo@bar.com"
	id := "inv_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"email_address":"%s"}`, emailAddress)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","email_address":"%s"}`, id, emailAddress)),
				Method: http.MethodPost,
				Path:   "/v1/invitations",
			},
		},
	}))

	invitation, err := Create(context.Background(), &CreateParams{
		EmailAddress: emailAddress,
	})
	require.NoError(t, err)
	require.Equal(t, id, invitation.ID)
	require.Equal(t, emailAddress, invitation.EmailAddress)
}

func TestInvitationCreate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
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
		},
	}))

	_, err := Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestInvitationRevoke(t *testing.T) {
	id := "inv_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","revoked":true,"status":"revoked"}`, id)),
				Method: http.MethodPost,
				Path:   "/v1/invitations/" + id + "/revoke",
			},
		},
	}))

	invitation, err := Revoke(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, invitation.ID)
	require.True(t, invitation.Revoked)
	require.Equal(t, "revoked", invitation.Status)
}

func TestInvitationRevoke_Error(t *testing.T) {
	id := "inv_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
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
		},
	}))

	_, err := Revoke(context.Background(), id)
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "revoke-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "revoke-error-code", apiErr.Errors[0].Code)
}
