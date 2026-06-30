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
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","provider":"%s","protocol":"%s"}`, name, provider, protocol)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","provider":"%s","protocol":"%s","object":"enterprise_connection","active":true}`, id, name, provider, protocol)),
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
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"enterprise_connection","protocol":"saml","provider":"okta","active":true}`, id, name)),
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
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"data":[{"id":"%s","name":"Test","object":"enterprise_connection","protocol":"saml","provider":"okta","active":true}],"total_count":1}`, id)),
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
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"enterprise_connection","protocol":"saml","provider":"okta","active":true}`, id, name)),
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

func TestEnterpriseConnectionClientUpdate_DisableJITProvisioning(t *testing.T) {
	t.Parallel()
	id := "entconn_abc"
	disableJITProvisioning := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"disable_jit_provisioning":%t}`, disableJITProvisioning)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_connection","protocol":"saml","provider":"okta","disable_jit_provisioning":%t}`, id, disableJITProvisioning)),
			Method: http.MethodPatch,
			Path:   "/v1/enterprise_connections/" + id,
		},
	}
	client := NewClient(config)
	conn, err := client.Update(context.Background(), id, &UpdateParams{
		DisableJITProvisioning: clerk.Bool(disableJITProvisioning),
	})
	require.NoError(t, err)
	require.Equal(t, id, conn.ID)
	require.NotNil(t, conn.DisableJITProvisioning)
	require.Equal(t, disableJITProvisioning, *conn.DisableJITProvisioning)
}

// TestEnterpriseConnectionClientCreate_WithMultiValuedCustomAttributes verifies the
// SAML enterprise connection create endpoint accepts custom attributes that carry
// the `multi_valued` flag, covering both true and false values.
func TestEnterpriseConnectionClientCreate_WithMultiValuedCustomAttributes(t *testing.T) {
	t.Parallel()
	id := "entconn_456"
	name := "Acme SAML"
	provider := "saml_custom"
	protocol := "saml"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","provider":"%s","protocol":"%s","saml":{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"manager","key":"manager","sso_path":"$.manager","scim_path":"manager","multi_valued":false}]}}`, name, provider, protocol)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","provider":"%s","protocol":"%s","object":"enterprise_connection","active":true,"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"manager","key":"manager","sso_path":"$.manager","scim_path":"manager","multi_valued":false}]}`, id, name, provider, protocol)),
			Method: http.MethodPost,
			Path:   "/v1/enterprise_connections",
		},
	}
	client := NewClient(config)
	conn, err := client.Create(context.Background(), &CreateParams{
		Name:     clerk.String(name),
		Provider: clerk.String(provider),
		Protocol: clerk.String(protocol),
		Saml: &CreateParamsSaml{
			CustomAttributes: &[]clerk.CustomAttribute{
				{
					Name:        clerk.String("groups"),
					Key:         clerk.String("groups"),
					SSOPath:     clerk.String("$.groups"),
					SCIMPath:    clerk.String("groups"),
					MultiValued: clerk.Bool(true),
				},
				{
					Name:        clerk.String("manager"),
					Key:         clerk.String("manager"),
					SSOPath:     clerk.String("$.manager"),
					SCIMPath:    clerk.String("manager"),
					MultiValued: clerk.Bool(false),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, conn.ID)
	require.NotNil(t, conn.CustomAttributes)
	require.Equal(t, 2, len(*conn.CustomAttributes))
	require.Equal(t, true, *(*conn.CustomAttributes)[0].MultiValued)
	require.Equal(t, false, *(*conn.CustomAttributes)[1].MultiValued)
}

// TestEnterpriseConnectionClientUpdate_WithMultiValuedCustomAttributes verifies the
// SAML enterprise connection update endpoint accepts custom attributes that carry
// the `multi_valued` flag, covering both true and false values.
func TestEnterpriseConnectionClientUpdate_WithMultiValuedCustomAttributes(t *testing.T) {
	t.Parallel()
	id := "entconn_789"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"saml":{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"department","key":"department","sso_path":"$.department","scim_path":"department","multi_valued":false}]}}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_connection","protocol":"saml","provider":"saml_custom","active":true,"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"department","key":"department","sso_path":"$.department","scim_path":"department","multi_valued":false}]}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/enterprise_connections/" + id,
		},
	}
	client := NewClient(config)
	conn, err := client.Update(context.Background(), id, &UpdateParams{
		Saml: &UpdateParamsSaml{
			CustomAttributes: &[]clerk.CustomAttribute{
				{
					Name:        clerk.String("groups"),
					Key:         clerk.String("groups"),
					SSOPath:     clerk.String("$.groups"),
					SCIMPath:    clerk.String("groups"),
					MultiValued: clerk.Bool(true),
				},
				{
					Name:        clerk.String("department"),
					Key:         clerk.String("department"),
					SSOPath:     clerk.String("$.department"),
					SCIMPath:    clerk.String("department"),
					MultiValued: clerk.Bool(false),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, conn.ID)
	require.NotNil(t, conn.CustomAttributes)
	require.Equal(t, 2, len(*conn.CustomAttributes))
	require.Equal(t, true, *(*conn.CustomAttributes)[0].MultiValued)
	require.Equal(t, false, *(*conn.CustomAttributes)[1].MultiValued)
}

func TestEnterpriseConnectionClientDelete(t *testing.T) {
	t.Parallel()
	id := "entconn_del"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","object":"enterprise_connection","deleted":true}`, id)),
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
