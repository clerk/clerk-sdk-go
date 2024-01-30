package domain

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

func TestDomainCreate(t *testing.T) {
	name := "clerk.com"
	id := "dmn_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:      t,
				in:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
				out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
				path:   "/v1/domains",
				method: http.MethodPost,
			},
		},
	}))

	dmn, err := Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	assert.Equal(t, id, dmn.ID)
	assert.Equal(t, name, dmn.Name)
}

func TestDomainCreate_Error(t *testing.T) {
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

func TestDomainUpdate(t *testing.T) {
	id := "dmn_456"
	name := "clerk.dev"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:      t,
				in:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
				out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
				path:   fmt.Sprintf("/v1/domains/%s", id),
				method: http.MethodPatch,
			},
		},
	}))

	dmn, err := Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	assert.Equal(t, id, dmn.ID)
	assert.Equal(t, name, dmn.Name)
}

func TestDomainUpdate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:      t,
				status: http.StatusBadRequest,
				out: json.RawMessage(`{
  "errors":[{
		"code":"update-error-code"
	}],
	"clerk_trace_id":"update-trace-id"
}`),
			},
		},
	}))

	_, err := Update(context.Background(), "dmn_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	assert.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	assert.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestDomainDelete(t *testing.T) {
	id := "dmn_789"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T:      t,
				out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
				path:   fmt.Sprintf("/v1/domains/%s", id),
				method: http.MethodDelete,
			},
		},
	}))

	dmn, err := Delete(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, dmn.ID)
	assert.True(t, dmn.Deleted)
}

func TestDomainList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				T: t,
				out: json.RawMessage(`{
	"data": [{"id":"dmn_123","name":"clerk.com"}],
	"total_count": 1
}`),
				path:   "/v1/domains",
				method: http.MethodGet,
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), list.TotalCount)
	assert.Equal(t, 1, len(list.Domains))
	assert.Equal(t, "dmn_123", list.Domains[0].ID)
	assert.Equal(t, "clerk.com", list.Domains[0].Name)
}

type mockRoundTripper struct {
	T *testing.T
	// Response status.
	status int
	// Response body.
	out json.RawMessage
	// If set, we'll assert that the request body
	// matches.
	in json.RawMessage
	// If set, we'll assert the request path matches.
	path string
	// If set, we'll assert that the request method matches.
	method string
}

func (rt *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt.status == 0 {
		rt.status = http.StatusOK
	}

	if rt.method != "" {
		require.Equal(rt.T, rt.method, r.Method)
	}
	if rt.path != "" {
		require.Equal(rt.T, rt.path, r.URL.Path)
	}
	if rt.in != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		require.JSONEq(rt.T, string(rt.in), string(body))
	}

	return &http.Response{
		StatusCode: rt.status,
		Body:       io.NopCloser(bytes.NewReader(rt.out)),
	}, nil
}
