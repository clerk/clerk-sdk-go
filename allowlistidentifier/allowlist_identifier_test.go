package allowlistidentifier

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

func TestAllowlistIdentifierCreate(t *testing.T) {
	identifier := "foo@bar.com"
	id := "alid_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"identifier":"%s"}`, identifier)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","identifier":"%s"}`, id, identifier)),
				Method: http.MethodPost,
				Path:   "/v1/allowlist_identifiers",
			},
		},
	}))

	allowlistIdentifier, err := Create(context.Background(), &CreateParams{
		Identifier: clerk.String(identifier),
	})
	require.NoError(t, err)
	require.Equal(t, id, allowlistIdentifier.ID)
	require.Equal(t, identifier, allowlistIdentifier.Identifier)
}

func TestAllowlistIdentifierCreate_Error(t *testing.T) {
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

func TestAllowlistIdentifierDelete(t *testing.T) {
	id := "alid_456"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
				Method: http.MethodDelete,
				Path:   "/v1/allowlist_identifiers/" + id,
			},
		},
	}))

	allowlistIdentifier, err := Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, allowlistIdentifier.ID)
	require.True(t, allowlistIdentifier.Deleted)
}

func TestAllowlistIdentifierList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [{"id":"alid_123","identifier":"foo@bar.com"}],
	"total_count": 1
}`),
				Method: http.MethodGet,
				Path:   "/v1/allowlist_identifiers",
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.AllowlistIdentifiers))
	require.Equal(t, "alid_123", list.AllowlistIdentifiers[0].ID)
	require.Equal(t, "foo@bar.com", list.AllowlistIdentifiers[0].Identifier)
}
