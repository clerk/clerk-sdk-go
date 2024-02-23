package jwttemplate

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

func TestJWTTemplateClientCreate(t *testing.T) {
	t.Parallel()
	name := "the-name"
	id := "jtmpl_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodPost,
			Path:   "/v1/jwt_templates",
		},
	}
	client := NewClient(config)
	jwtTemplate, err := client.Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, jwtTemplate.ID)
	require.Equal(t, name, jwtTemplate.Name)
}

func TestJWTTemplateClientCreate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
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

func TestJWTTemplateClientGet(t *testing.T) {
	t.Parallel()
	id := "jtmpl_123"
	name := "the-name"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodGet,
			Path:   "/v1/jwt_templates/" + id,
		},
	}
	client := NewClient(config)
	jwtTemplate, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, jwtTemplate.ID)
	require.Equal(t, name, jwtTemplate.Name)
}

func TestJWTTemplateClientUpdate(t *testing.T) {
	t.Parallel()
	id := "jtmpl_123"
	name := "the-name"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodPatch,
			Path:   "/v1/jwt_templates/" + id,
		},
	}
	client := NewClient(config)
	jwtTemplate, err := client.Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, jwtTemplate.ID)
	require.Equal(t, name, jwtTemplate.Name)
}

func TestJWTTemplateClientUpdate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
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
	}
	client := NewClient(config)
	_, err := client.Update(context.Background(), "jtmpl_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestJWTTemplateClientDelete(t *testing.T) {
	t.Parallel()
	id := "jtmpl_456"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/jwt_templates/" + id,
		},
	}
	client := NewClient(config)
	jwtTemplate, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, jwtTemplate.ID)
	require.True(t, jwtTemplate.Deleted)
}

func TestJWTTemplateClientList(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
	"data": [{"id":"jtmpl_123","name":"the-name"}],
	"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/jwt_templates",
			Query: &url.Values{
				"paginated": []string{"true"},
			},
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.JWTTemplates))
	require.Equal(t, "jtmpl_123", list.JWTTemplates[0].ID)
	require.Equal(t, "the-name", list.JWTTemplates[0].Name)
}
