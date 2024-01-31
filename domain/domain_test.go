package domain

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

func TestDomainCreate(t *testing.T) {
	name := "clerk.com"
	id := "dmn_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
				Path:   "/v1/domains",
				Method: http.MethodPost,
			},
		},
	}))

	dmn, err := Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, dmn.ID)
	require.Equal(t, name, dmn.Name)
}

func TestDomainCreate_Error(t *testing.T) {
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

func TestDomainUpdate(t *testing.T) {
	id := "dmn_456"
	name := "clerk.dev"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
				Path:   fmt.Sprintf("/v1/domains/%s", id),
				Method: http.MethodPatch,
			},
		},
	}))

	dmn, err := Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, dmn.ID)
	require.Equal(t, name, dmn.Name)
}

func TestDomainUpdate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Status: http.StatusBadRequest,
				Out: json.RawMessage(`{
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
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestDomainDelete(t *testing.T) {
	id := "dmn_789"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
				Path:   fmt.Sprintf("/v1/domains/%s", id),
				Method: http.MethodDelete,
			},
		},
	}))

	dmn, err := Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, dmn.ID)
	require.True(t, dmn.Deleted)
}

func TestDomainList(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
	"data": [{"id":"dmn_123","name":"clerk.com"}],
	"total_count": 1
}`),
				Path:   "/v1/domains",
				Method: http.MethodGet,
			},
		},
	}))

	list, err := List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Domains))
	require.Equal(t, "dmn_123", list.Domains[0].ID)
	require.Equal(t, "clerk.com", list.Domains[0].Name)
}
