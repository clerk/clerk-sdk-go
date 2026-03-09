package enterpriseconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestEnterpriseConnectionClientCreate(t *testing.T) {
	t.Parallel()
	id := "entconn_123"
	name := "My Enterprise"
	provider := "saml_custom"
	protocol := "saml"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(fmt.Sprintf(`{"name":"%s","provider":"%s","protocol":"%s"}`, name, provider, protocol)),
			Out: json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","provider":"%s","protocol":"%s","object":"enterprise_connection","active":true}`, id, name, provider, protocol)),
			Method: http.MethodPost,
			Path:   "/v1/enterprise_connections",
		},
	}
	client := NewClient(config)
	conn, err := client.Create(context.Background(), &CreateParams{
		Name:     clerk.String(name),
		Provider: clerk.String(provider),
		Protocol: clerk.String(protocol),
	})
	require.NoError(t, err)
	require.Equal(t, id, conn.ID)
	require.Equal(t, name, conn.Name)
	require.Equal(t, provider, conn.Provider)
	require.Equal(t, protocol, conn.Protocol)
}

func TestEnterpriseConnectionClientGet(t *testing.T) {
	t.Parallel()
	id := "entconn_456"
	name := "Acme Corp"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"enterprise_connection","protocol":"saml","provider":"okta","active":true}`, id, name)),
			Method: http.MethodGet,
			Path:   "/v1/enterprise_connections/" + id,
		},
	}
	client := NewClient(config)
	conn, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, conn.ID)
	require.Equal(t, name, conn.Name)
}

func TestEnterpriseConnectionClientList(t *testing.T) {
	t.Parallel()
	id := "entconn_789"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{"data":[{"id":"%s","name":"Test","object":"enterprise_connection","protocol":"saml","provider":"okta","active":true}],"total_count":1}`, id)),
			Method: http.MethodGet,
			Path:   "/v1/enterprise_connections",
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, 1, len(list.EnterpriseConnections))
	require.Equal(t, id, list.EnterpriseConnections[0].ID)
	require.Equal(t, int64(1), list.TotalCount)
}

func TestEnterpriseConnectionClientUpdate(t *testing.T) {
	t.Parallel()
	id := "entconn_abc"
	name := "Updated Name"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:  t,
			In: json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out: json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"enterprise_connection","protocol":"saml","provider":"okta","active":true}`, id, name)),
			Method: http.MethodPatch,
			Path:   "/v1/enterprise_connections/" + id,
		},
	}
	client := NewClient(config)
	conn, err := client.Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, conn.ID)
	require.Equal(t, name, conn.Name)
}

func TestEnterpriseConnectionClientDelete(t *testing.T) {
	t.Parallel()
	id := "entconn_del"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_connection","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/enterprise_connections/" + id,
		},
	}
	client := NewClient(config)
	res, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.True(t, res.Deleted)
	require.Equal(t, id, res.ID)
}
