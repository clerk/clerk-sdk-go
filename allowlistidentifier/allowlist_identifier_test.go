package allowlistidentifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllowlistIdentifierCreate(t *testing.T) {
	identifier := "foo@bar.com"
	id := "alid_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:   t,
				in:  json.RawMessage(fmt.Sprintf(`{"identifier":"%s"}`, identifier)),
				out: json.RawMessage(fmt.Sprintf(`{"id":"%s","identifier":"%s"}`, id, identifier)),
			},
		},
	}))

	allowlistIdentifier, err := Create(context.Background(), &CreateParams{
		Identifier: clerk.String(identifier),
	})
	require.NoError(t, err)
	assert.Equal(t, id, allowlistIdentifier.ID)
	assert.Equal(t, identifier, allowlistIdentifier.Identifier)
}

func TestAllowlistIdentifierCreate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:      t,
				status: http.StatusBadRequest,
				out: json.RawMessage(`{
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
	assert.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	assert.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestAllowlistIdentifierDelete(t *testing.T) {
	id := "alid_456"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:   t,
				out: json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			},
		},
	}))

	allowlistIdentifier, err := Delete(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, allowlistIdentifier.ID)
	assert.True(t, allowlistIdentifier.Deleted)
}

func TestAllowlistIdentifierList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T: t,
				out: json.RawMessage(`{
	"data": [{"id":"alid_123","identifier":"foo@bar.com"}],
	"total_count": 1
}`),
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), list.TotalCount)
	assert.Equal(t, 1, len(list.AllowlistIdentifiers))
	assert.Equal(t, "alid_123", list.AllowlistIdentifiers[0].ID)
	assert.Equal(t, "foo@bar.com", list.AllowlistIdentifiers[0].Identifier)
}

type mockRoundTripper struct {
	T      *testing.T
	status int
	in     json.RawMessage
	out    json.RawMessage
}

func (rt *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt.status == 0 {
		rt.status = http.StatusOK
	}
	if rt.in != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		assert.JSONEq(rt.T, string(rt.in), string(body))
	}
	return &http.Response{
		StatusCode: rt.status,
		Body:       io.NopCloser(bytes.NewReader(rt.out)),
	}, nil
}
