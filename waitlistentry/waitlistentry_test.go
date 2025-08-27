package waitlistentry

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

func TestWaitlistList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [{"id":"wle_123","email_address":"foo@bar.com"}],
	"total_count": 1
}`),
				Path:   "/v1/waitlist_entries",
				Method: http.MethodGet,
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.WaitlistEntries))
	require.Equal(t, "wle_123", list.WaitlistEntries[0].ID)
	require.Equal(t, "foo@bar.com", list.WaitlistEntries[0].EmailAddress)
}

func TestWaitlistEntriesListWithParams(t *testing.T) {
	limit := int64(10)
	offset := int64(20)
	orderBy := "-created_at"
	query := "example@email.com"
	status1 := "pending"
	status2 := "invited"

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [
		{"id":"wle_123","email_address":"foo@bar.com"},
		{"id":"wle_124","email_address":"baz@qux.com","invitation":{"id":"inv_124","email_address":"baz@qux.com"}}
	],
	"total_count": 2
}`),
				Path:   "/v1/waitlist_entries",
				Method: http.MethodGet,
				Query: &url.Values{
					"limit":    []string{fmt.Sprintf("%d", limit)},
					"offset":   []string{fmt.Sprintf("%d", offset)},
					"order_by": []string{orderBy},
					"query":    []string{query},
					"status":   []string{status1, status2},
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
	require.Len(t, list.WaitlistEntries, 2)
	require.Equal(t, "wle_123", list.WaitlistEntries[0].ID)
	require.Equal(t, "foo@bar.com", list.WaitlistEntries[0].EmailAddress)
	require.Nil(t, list.WaitlistEntries[0].Invitation)
	require.Equal(t, "wle_124", list.WaitlistEntries[1].ID)
	require.Equal(t, "baz@qux.com", list.WaitlistEntries[1].EmailAddress)
	require.NotNil(t, list.WaitlistEntries[1].Invitation)
	require.Equal(t, "inv_124", list.WaitlistEntries[1].Invitation.ID)
	require.Equal(t, "baz@qux.com", list.WaitlistEntries[1].Invitation.EmailAddress)
}

func TestWaitlistEntryCreate(t *testing.T) {
	emailAddress := "foo@bar.com"
	id := "inv_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"email_address":"%s"}`, emailAddress)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","email_address":"%s"}`, id, emailAddress)),
				Method: http.MethodPost,
				Path:   "/v1/waitlist_entries",
			},
		},
	}))

	entry, err := Create(context.Background(), &CreateParams{
		EmailAddress: emailAddress,
	})
	require.NoError(t, err)
	require.Equal(t, id, entry.ID)
	require.Equal(t, emailAddress, entry.EmailAddress)
}

func TestWaitlistEntryCreate_Error(t *testing.T) {
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

func TestWaitlistEntryInvite(t *testing.T) {
	id := "wle_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"invited"}`, id)),
				Method: http.MethodPost,
				Path:   "/v1/waitlist_entries/" + id + "/invite",
			},
		},
	}))

	entry, err := Invite(context.Background(), id, &InviteParams{})
	require.NoError(t, err)
	require.Equal(t, id, entry.ID)
	require.Equal(t, "invited", entry.Status)
}

func TestWaitlistEntryInvite_Error(t *testing.T) {
	id := "wle_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Status: http.StatusBadRequest,
				Out: json.RawMessage(`{
  "errors":[{
		"code":"invite-error-code"
	}],
	"clerk_trace_id":"invite-trace-id"
}`),
				Path: "/v1/waitlist_entries/" + id + "/invite",
			},
		},
	}))

	_, err := Invite(context.Background(), id, &InviteParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "invite-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "invite-error-code", apiErr.Errors[0].Code)
}

func TestWaitlistEntryReject(t *testing.T) {
	id := "wle_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"rejected"}`, id)),
				Method: http.MethodPost,
				Path:   "/v1/waitlist_entries/" + id + "/reject",
			},
		},
	}))

	entry, err := Reject(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, entry.ID)
	require.Equal(t, "rejected", entry.Status)
}

func TestWaitlistEntryReject_Error(t *testing.T) {
	id := "wle_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Status: http.StatusBadRequest,
				Out: json.RawMessage(`{
  "errors":[{
		"code":"reject-error-code"
	}],
	"clerk_trace_id":"reject-trace-id"
}`),
				Path: "/v1/waitlist_entries/" + id + "/reject",
			},
		},
	}))

	_, err := Reject(context.Background(), id)
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "reject-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "reject-error-code", apiErr.Errors[0].Code)
}
