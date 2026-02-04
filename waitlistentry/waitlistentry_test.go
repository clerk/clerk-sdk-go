package waitlistentry

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

func TestWaitlistEntryBulkCreate(t *testing.T) {
	emailAddresses := []string{"foo@bar.com", "bar@foo.com"}
	ids := []string{"wle_123", "wle_456"}
	createdAt := int64(1700000000)
	updatedAt := int64(1700000100)

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:  t,
				In: json.RawMessage(fmt.Sprintf(`[{"email_address":"%s"},{"email_address":"%s"}]`, emailAddresses[0], emailAddresses[1])),
				Out: json.RawMessage(fmt.Sprintf(
					`[{"object":"waitlist_entry","id":"%s","email_address":"%s","status":"pending","is_locked":false,"created_at":%d,"updated_at":%d,"invitation":null},{"object":"waitlist_entry","id":"%s","email_address":"%s","status":"pending","is_locked":false,"created_at":%d,"updated_at":%d,"invitation":null}]`,
					ids[0], emailAddresses[0], createdAt, updatedAt, ids[1], emailAddresses[1], createdAt, updatedAt,
				)),
				Method: http.MethodPost,
				Path:   "/v1/waitlist_entries/bulk",
			},
		},
	}))

	params := BulkCreateParams{
		WaitlistEntries: []*CreateParams{
			{EmailAddress: emailAddresses[0]},
			{EmailAddress: emailAddresses[1]},
		},
	}

	response, err := BulkCreate(context.Background(), &params)
	require.NoError(t, err)
	require.Len(t, response.WaitlistEntries, 2)

	for i, entry := range response.WaitlistEntries {
		require.Equal(t, "waitlist_entry", entry.Object)
		require.Equal(t, ids[i], entry.ID)
		require.Equal(t, emailAddresses[i], entry.EmailAddress)
		require.Equal(t, "pending", entry.Status)
		require.False(t, entry.IsLocked)
		require.Equal(t, createdAt, entry.CreatedAt)
		require.Equal(t, updatedAt, entry.UpdatedAt)
		require.Nil(t, entry.Invitation)
	}
}
