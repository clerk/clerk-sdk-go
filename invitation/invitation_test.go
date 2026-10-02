package invitation

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

func TestInvitationListWithParams(t *testing.T) {
	limit := int64(10)
	offset := int64(20)
	orderBy := "-created_at"
	query := "example@email.com"
	status1 := "pending"
	status2 := "accepted"

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [
		{"id":"inv_123","email_address":"foo@bar.com"},
		{"id":"inv_124","email_address":"baz@qux.com"}
	],
	"total_count": 2
}`),
				Path:   "/v1/invitations",
				Method: http.MethodGet,
				Query: &url.Values{
					"limit":     []string{fmt.Sprintf("%d", limit)},
					"offset":    []string{fmt.Sprintf("%d", offset)},
					"order_by":  []string{orderBy},
					"query":     []string{query},
					"status":    []string{status1, status2},
					"paginated": []string{"true"},
				},
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{
		ListParams: clerk.ListParams{
			Limit:  &limit,
			Offset: &offset,
		},
		OrderBy:  &orderBy,
		Query:    &query,
		Statuses: []string{status1, status2},
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), list.TotalCount)
	require.Equal(t, 2, len(list.Invitations))
	require.Equal(t, "inv_123", list.Invitations[0].ID)
	require.Equal(t, "foo@bar.com", list.Invitations[0].EmailAddress)
	require.Equal(t, "inv_124", list.Invitations[1].ID)
	require.Equal(t, "baz@qux.com", list.Invitations[1].EmailAddress)
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

func TestInvitationCreateWithExpiration(t *testing.T) {
	emailAddress := "foo@bar.com"
	id := "inv_123"
	expiresInDays := int64(7)
	expiresAt := int64(1700000000)
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"email_address":"%s","expires_in_days":%d}`, emailAddress, expiresInDays)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","email_address":"%s","expires_at":%d}`, id, emailAddress, expiresAt)),
				Method: http.MethodPost,
				Path:   "/v1/invitations",
			},
		},
	}))

	invitation, err := Create(context.Background(), &CreateParams{
		EmailAddress:  emailAddress,
		ExpiresInDays: &expiresInDays,
	})
	require.NoError(t, err)
	require.Equal(t, id, invitation.ID)
	require.Equal(t, emailAddress, invitation.EmailAddress)
	require.Equal(t, expiresAt, *invitation.ExpiresAt)
}

func TestInvitationCreateWithTemplateSlug(t *testing.T) {
	emailAddress := "foo@bar.com"
	id := "inv_123"
	templateSlug := "template-slug"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"email_address":"%s","template_slug":"%s"}`, emailAddress, templateSlug)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","email_address":"%s"}`, id, emailAddress)),
				Method: http.MethodPost,
				Path:   "/v1/invitations",
			},
		},
	}))

	invitation, err := Create(context.Background(), &CreateParams{
		EmailAddress: emailAddress,
		TemplateSlug: &templateSlug,
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

func TestBulkInvitationCreate(t *testing.T) {
	emailAddresses := []string{"foo@bar.com", "bar@foo.com"}
	ids := []string{"inv_123", "inv_456"}
	invitations := []*clerk.Invitation{
		{ID: ids[0], EmailAddress: emailAddresses[0]},
		{ID: ids[1], EmailAddress: emailAddresses[1]},
	}

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:  t,
				In: json.RawMessage(fmt.Sprintf(`[{"email_address":"%s"},{"email_address":"%s"}]`, emailAddresses[0], emailAddresses[1])),
				Out: json.RawMessage(fmt.Sprintf(
					`[{"id":"%s","email_address":"%s"},{"id":"%s","email_address":"%s"}]`,
					ids[0], emailAddresses[0], ids[1], emailAddresses[1],
				)),
				Method: http.MethodPost,
				Path:   "/v1/invitations/bulk",
			},
		},
	}))

	params := BulkCreateParams{
		Invitations: []*CreateParams{
			{EmailAddress: emailAddresses[0]},
			{EmailAddress: emailAddresses[1]},
		},
	}

	response, err := BulkCreate(context.Background(), &params)
	require.NoError(t, err)
	require.Len(t, invitations, 2)

	for i, invitation := range response.Invitations {
		require.Equal(t, ids[i], invitation.ID)
		require.Equal(t, emailAddresses[i], invitation.EmailAddress)
	}
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

func TestInvitationDelete(t *testing.T) {
	id := "inv_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"invitation","deleted":true}`, id)),
				Path:   fmt.Sprintf("/v1/invitations/%s", id),
				Method: http.MethodDelete,
			},
		},
	}))

	invitation, err := Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, invitation.ID)
	require.True(t, invitation.Deleted)
}
