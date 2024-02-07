package blocklistidentifier

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

func TestBlocklistIdentifierClientCreate(t *testing.T) {
	t.Parallel()
	identifier := "foo@bar.com"
	id := "blid_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"identifier":"%s"}`, identifier)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","identifier":"%s"}`, id, identifier)),
			Method: http.MethodPost,
			Path:   "/v1/blocklist_identifiers",
		},
	}
	client := NewClient(config)
	blocklistIdentifier, err := client.Create(context.Background(), &CreateParams{
		Identifier: clerk.String(identifier),
	})
	require.NoError(t, err)
	require.Equal(t, id, blocklistIdentifier.ID)
	require.Equal(t, identifier, blocklistIdentifier.Identifier)
}

func TestBlocklistIdentifierClientCreate_Error(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
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

func TestBlocklistIdentifierClientDelete(t *testing.T) {
	t.Parallel()
	id := "blid_456"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/blocklist_identifiers/" + id,
		},
	}
	client := NewClient(config)
	blocklistIdentifier, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, blocklistIdentifier.ID)
	require.True(t, blocklistIdentifier.Deleted)
}

func TestBlocklistIdentifierClientList(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
	"data": [{"id":"blid_123","identifier":"foo@bar.com"}],
	"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/blocklist_identifiers",
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.BlocklistIdentifiers))
	require.Equal(t, "blid_123", list.BlocklistIdentifiers[0].ID)
	require.Equal(t, "foo@bar.com", list.BlocklistIdentifiers[0].Identifier)
}
