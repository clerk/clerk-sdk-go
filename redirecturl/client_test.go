package redirecturl

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

func TestRedirectURLClientCreate(t *testing.T) {
	t.Parallel()
	id := "ru_123"
	url := "https://example.com"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"url":"%s"}`, url)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","url":"%s"}`, id, url)),
			Method: http.MethodPost,
			Path:   "/v1/redirect_urls",
		},
	}
	client := NewClient(config)
	redirectURL, err := client.Create(context.Background(), &CreateParams{
		URL: clerk.String(url),
	})
	require.NoError(t, err)
	require.Equal(t, id, redirectURL.ID)
	require.Equal(t, url, redirectURL.URL)
}

func TestRedirectURLClientCreate_Error(t *testing.T) {
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

func TestRedirectURLClientGet(t *testing.T) {
	t.Parallel()
	id := "ru_123"
	url := "https://example.com"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","url":"%s"}`, id, url)),
			Method: http.MethodGet,
			Path:   "/v1/redirect_urls/" + id,
		},
	}
	client := NewClient(config)
	redirectURL, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, redirectURL.ID)
	require.Equal(t, url, redirectURL.URL)
}

func TestRedirectURLClientDelete(t *testing.T) {
	t.Parallel()
	id := "ru_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/redirect_urls/" + id,
		},
	}
	client := NewClient(config)
	redirectURL, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, redirectURL.ID)
	require.True(t, redirectURL.Deleted)
}

func TestRedirectURLClientList(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
"data": [{"id":"ru_123","url":"https://example.com"}],
"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/redirect_urls",
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.RedirectURLs))
	require.Equal(t, "ru_123", list.RedirectURLs[0].ID)
	require.Equal(t, "https://example.com", list.RedirectURLs[0].URL)
}
