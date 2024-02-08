package organization

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

func TestOrganizationClientCreate(t *testing.T) {
	t.Parallel()
	id := "org_123"
	name := "Acme Inc"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodPost,
			Path:   "/v1/organizations",
		},
	}
	client := NewClient(config)
	organization, err := client.Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.Equal(t, name, organization.Name)
}

func TestOrganizationClientCreate_Error(t *testing.T) {
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

func TestOrganizationClientGet(t *testing.T) {
	t.Parallel()
	id := "org_123"
	name := "Acme Inc"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodGet,
			Path:   "/v1/organizations/" + id,
		},
	}
	client := NewClient(config)
	organization, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.Equal(t, name, organization.Name)
}

func TestOrganizationClientUpdate(t *testing.T) {
	t.Parallel()
	id := "org_123"
	name := "Acme Inc"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s"}`, id, name)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + id,
		},
	}
	client := NewClient(config)
	organization, err := client.Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.Equal(t, name, organization.Name)
}

func TestOrganizationClientUpdate_Error(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
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
	_, err := client.Update(context.Background(), "org_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestOrganizationClientDelete(t *testing.T) {
	t.Parallel()
	id := "org_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + id,
		},
	}
	client := NewClient(config)
	organization, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.True(t, organization.Deleted)
}

func TestOrganizationClientList(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
"data": [{"id":"org_123","name":"Acme Inc"}],
"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/organizations",
			Query: &url.Values{
				"limit":    []string{"1"},
				"offset":   []string{"2"},
				"order_by": []string{"-created_at"},
				"query":    []string{"Acme"},
				"user_id":  []string{"user_123", "user_456"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		OrderBy: clerk.String("-created_at"),
		Query:   clerk.String("Acme"),
		UserIDs: []string{"user_123", "user_456"},
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Organizations))
	require.Equal(t, "org_123", list.Organizations[0].ID)
	require.Equal(t, "Acme Inc", list.Organizations[0].Name)
}

func TestOrganizationClientDeleteLogo(t *testing.T) {
	t.Parallel()
	id := "org_123"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/organizations/" + id + "/logo",
		},
	}
	client := NewClient(config)
	organization, err := client.DeleteLogo(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
}

func TestOrganizationClientUpdateMetadata(t *testing.T) {
	t.Parallel()
	id := "org_123"
	metadata := `{"foo":"bar"}`
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"private_metadata":%s}`, metadata)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","private_metadata":%s}`, id, metadata)),
			Method: http.MethodPatch,
			Path:   "/v1/organizations/" + id + "/metadata",
		},
	}
	client := NewClient(config)
	metadataParam := json.RawMessage(metadata)
	organization, err := client.UpdateMetadata(context.Background(), id, &UpdateMetadataParams{
		PrivateMetadata: &metadataParam,
	})
	require.NoError(t, err)
	require.Equal(t, id, organization.ID)
	require.JSONEq(t, metadata, string(organization.PrivateMetadata))
}
