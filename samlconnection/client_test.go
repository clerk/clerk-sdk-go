package samlconnection

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

func TestSAMLConnectionClientCreate(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s","provider":"%s"}`, name, domain, provider)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s"}`, id, name, domain, provider)),
			Method: http.MethodPost,
			Path:   "/v1/saml_connections",
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Create(context.Background(), &CreateParams{
		Name:     clerk.String(name),
		Domain:   clerk.String(domain),
		Provider: clerk.String(provider),
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	require.Equal(t, domain, samlConnection.Domain)
	require.Equal(t, provider, samlConnection.Provider)
}

func TestSAMLConnectionClientCreate_Error(t *testing.T) {
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

func TestSAMLConnectionClientGet(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	disableAdditionalIdentifications := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:   t,
			Out: json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s", "disable_additional_identifications": %t}`, id, name, domain, provider, disableAdditionalIdentifications)), Method: http.MethodGet,
			Path: "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	require.Equal(t, domain, samlConnection.Domain)
	require.Equal(t, provider, samlConnection.Provider)
	require.Equal(t, disableAdditionalIdentifications, samlConnection.DisableAdditionalIdentifications)
}

func TestSAMLConnectionClientUpdate(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	disableAdditionalIdentifications := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s", "disable_additional_identifications": %t, "organization_id": ""}`, name, disableAdditionalIdentifications)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","disable_additional_identifications": %t}`, id, name, domain, provider, disableAdditionalIdentifications)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		Name:                             clerk.String(name),
		DisableAdditionalIdentifications: clerk.Bool(disableAdditionalIdentifications),
		OrganizationID:                   clerk.String(""),
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	require.Equal(t, disableAdditionalIdentifications, samlConnection.DisableAdditionalIdentifications)
}

func TestSAMLConnectionClientUpdate_Error(t *testing.T) {
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

func TestSAMLConnectionClientDelete(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.True(t, samlConnection.Deleted)
}

func TestSAMLConnectionClientList(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{
	"data": [{"id":"%s","name":"%s","domain":"%s","provider":"%s"}],
	"total_count": 1
}`, id, name, domain, provider)),
			Method: http.MethodGet,
			Path:   "/v1/saml_connections",
			Query: &url.Values{
				"limit":    []string{"1"},
				"query":    []string{"Acme"},
				"order_by": []string{"-created_at"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		ListParams: clerk.ListParams{
			Limit: clerk.Int64(1),
		},
		Query:   clerk.String("Acme"),
		OrderBy: clerk.String("-created_at"),
	}
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.SAMLConnections))
	require.Equal(t, id, list.SAMLConnections[0].ID)
	require.Equal(t, name, list.SAMLConnections[0].Name)
	require.Equal(t, domain, list.SAMLConnections[0].Domain)
	require.Equal(t, provider, list.SAMLConnections[0].Provider)
}
