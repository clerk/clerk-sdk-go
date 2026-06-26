package samlconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestSAMLConnectionClientCreate(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	forceAuthn := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s","provider":"%s","force_authn":%t}`, name, domain, provider, forceAuthn)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","force_authn":%t}`, id, name, domain, provider, forceAuthn)),
			Method: http.MethodPost,
			Path:   "/v1/saml_connections",
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Create(context.Background(), &CreateParams{
		Name:       clerk.String(name),
		Domain:     clerk.String(domain),
		Provider:   clerk.String(provider),
		ForceAuthn: clerk.Bool(forceAuthn),
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	// nolint:staticcheck // we want to test the .Domain behavior when it's deprecated
	require.Equal(t, domain, samlConnection.Domain)
	require.Equal(t, provider, samlConnection.Provider)
	require.Equal(t, forceAuthn, samlConnection.ForceAuthn)
}

// TestSAMLConnectionClientCreate_WithBothDomainAndDomains tests that the client can not create a SAML connection
// When providing both domain and domains. An error is returned.
func TestSAMLConnectionClientCreate_WithBothDomainAndDomains(t *testing.T) {
	t.Parallel()
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s", "domains": ["%s"], "provider":"%s"}`, name, domain, domain, provider)),
			Out:    json.RawMessage(`{ "clerk_trace_id": "trace-id", "errors": [{"code": "form_conditional_param_disallowed", "short_message": "is not allowed", "long_message": "domain isn't allowed when domains is present.", "meta": {"param_name": "domain"}}]}`),
			Method: http.MethodPost,
			Path:   "/v1/saml_connections",
			Status: http.StatusUnprocessableEntity,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Create(context.Background(), &CreateParams{
		Name:     clerk.String(name),
		Domain:   clerk.String(domain),
		Domains:  &[]string{domain},
		Provider: clerk.String(provider),
	})
	require.Error(t, err)
	require.Empty(t, samlConnection.ID)
}

// TestSAMLConnectionClientCreate_WithDomains tests that the client can create a SAML connection
// When providing only domains.
func TestSAMLConnectionClientCreate_WithDomains(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domainA := "example.com"
	domainB := "example.org"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domains": ["%s", "%s"], "provider":"%s"}`, name, domainA, domainB, provider)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain": "%s", "domains": ["%s", "%s"],"provider":"%s"}`, id, name, domainA, domainA, domainB, provider)),
			Method: http.MethodPost,
			Path:   "/v1/saml_connections",
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Create(context.Background(), &CreateParams{
		Name:     clerk.String(name),
		Domains:  &[]string{domainA, domainB},
		Provider: clerk.String(provider),
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	// nolint:staticcheck // we want to test the .Domain behavior when it's deprecated
	require.Equal(t, domainA, samlConnection.Domain)
	require.Equal(t, []string{domainA, domainB}, samlConnection.Domains)
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
	forceAuthn := false
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s", "disable_additional_identifications": %t, "force_authn": %t, "enterprise_connection_id": "entconn_2abc123def456", "custom_attributes": [{"name": "custom_attribute_name", "key": "custom_attribute_key", "sso_path": "custom_attribute_sso_path", "scim_path": "custom_attribute_scim_path"}]}`, id, name, domain, provider, disableAdditionalIdentifications, forceAuthn)),
			Method: http.MethodGet,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	// nolint:staticcheck // we want to test the .Domain behavior when it's deprecated
	require.Equal(t, domain, samlConnection.Domain)
	require.Equal(t, provider, samlConnection.Provider)
	require.Equal(t, disableAdditionalIdentifications, samlConnection.DisableAdditionalIdentifications)
	require.Equal(t, forceAuthn, samlConnection.ForceAuthn)
	require.Equal(t, &[]clerk.CustomAttribute{
		{
			Name:     clerk.String("custom_attribute_name"),
			Key:      clerk.String("custom_attribute_key"),
			SSOPath:  clerk.String("custom_attribute_sso_path"),
			SCIMPath: clerk.String("custom_attribute_scim_path"),
		},
	}, samlConnection.CustomAttributes)
}

func TestSAMLConnectionClientUpdate(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	disableAdditionalIdentifications := true
	forceAuthn := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s", "disable_additional_identifications": %t, "organization_id": "", "force_authn": %t, "custom_attributes": [{"name": "custom_attribute_name", "key": "custom_attribute_key", "sso_path": "custom_attribute_sso_path", "scim_path": "custom_attribute_scim_path"}]}`, name, disableAdditionalIdentifications, forceAuthn)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","disable_additional_identifications": %t,"force_authn": %t, "enterprise_connection_id": null, "custom_attributes": [{"name": "custom_attribute_name", "key": "custom_attribute_key", "sso_path": "custom_attribute_sso_path", "scim_path": "custom_attribute_scim_path"}]}`, id, name, domain, provider, disableAdditionalIdentifications, forceAuthn)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		Name:                             clerk.String(name),
		DisableAdditionalIdentifications: clerk.Bool(disableAdditionalIdentifications),
		OrganizationID:                   clerk.String(""),
		ForceAuthn:                       clerk.Bool(forceAuthn),
		CustomAttributes: &[]clerk.CustomAttribute{
			{
				Name:     clerk.String("custom_attribute_name"),
				Key:      clerk.String("custom_attribute_key"),
				SSOPath:  clerk.String("custom_attribute_sso_path"),
				SCIMPath: clerk.String("custom_attribute_scim_path"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, &[]clerk.CustomAttribute{
		{
			Name:     clerk.String("custom_attribute_name"),
			Key:      clerk.String("custom_attribute_key"),
			SSOPath:  clerk.String("custom_attribute_sso_path"),
			SCIMPath: clerk.String("custom_attribute_scim_path"),
		},
	}, samlConnection.CustomAttributes)
	require.Equal(t, name, samlConnection.Name)
	require.Equal(t, disableAdditionalIdentifications, samlConnection.DisableAdditionalIdentifications)
	require.Equal(t, forceAuthn, samlConnection.ForceAuthn)
}

func TestSAMLConnectionClientUpdate_DisableJitProvisioning(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	disableJitProvisioning := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","disable_jit_provisioning":%t}`, name, disableJitProvisioning)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","disable_jit_provisioning":%t}`, id, name, domain, provider, disableJitProvisioning)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		Name:                   clerk.String(name),
		DisableJitProvisioning: clerk.Bool(disableJitProvisioning),
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.NotNil(t, samlConnection.DisableJitProvisioning)
	require.Equal(t, disableJitProvisioning, *samlConnection.DisableJitProvisioning)
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

// TestSAMLConnectionClientUpdate_WithBothDomainAndDomains tests that the client can not update a SAML connection
// When providing both domain and domains. An error is returned.
func TestSAMLConnectionClientUpdate_WithBothDomainAndDomains(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s", "domains": ["%s"]}`, name, domain, domain)),
			Out:    json.RawMessage(`{ "clerk_trace_id": "trace-id", "errors": [{"code": "form_conditional_param_disallowed", "short_message": "is not allowed", "long_message": "domain isn't allowed when domains is present.", "meta": {"param_name": "domain"}}]}`),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
			Status: http.StatusUnprocessableEntity,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		Name:    clerk.String(name),
		Domain:  clerk.String(domain),
		Domains: &[]string{domain},
	})
	require.Error(t, err)
	require.Empty(t, samlConnection.ID)
}

// TestSAMLConnectionClientUpdate_WithDomains tests that the client can update a SAML connection
// When providing only domains.
func TestSAMLConnectionClientUpdate_WithDomains(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domainA := "example.com"
	domainB := "example.org"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domains": ["%s", "%s"]}`, name, domainA, domainB)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain": "%s", "domains": ["%s", "%s"],"provider":"%s"}`, id, name, domainA, domainA, domainB, provider)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		Name:    clerk.String(name),
		Domains: &[]string{domainA, domainB},
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	// nolint:staticcheck // we want to test the .Domain behavior when it's deprecated
	require.Equal(t, domainA, samlConnection.Domain)
	require.Equal(t, []string{domainA, domainB}, samlConnection.Domains)
	require.Equal(t, provider, samlConnection.Provider)
}

// TestSAMLConnectionClientCreate_WithMultiValuedCustomAttributes verifies the client
// can create a SAML connection with custom attributes that include the `multi_valued`
// flag, covering both true and false values to ensure the boolean serializes correctly
// in either state.
func TestSAMLConnectionClientCreate_WithMultiValuedCustomAttributes(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s","provider":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"manager","key":"manager","sso_path":"$.manager","scim_path":"manager","multi_valued":false}]}`, name, domain, provider)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"manager","key":"manager","sso_path":"$.manager","scim_path":"manager","multi_valued":false}]}`, id, name, domain, provider)),
			Method: http.MethodPost,
			Path:   "/v1/saml_connections",
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Create(context.Background(), &CreateParams{
		Name:     clerk.String(name),
		Domain:   clerk.String(domain),
		Provider: clerk.String(provider),
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
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.NotNil(t, samlConnection.CustomAttributes)
	require.Equal(t, 2, len(*samlConnection.CustomAttributes))
	require.Equal(t, true, *(*samlConnection.CustomAttributes)[0].MultiValued)
	require.Equal(t, false, *(*samlConnection.CustomAttributes)[1].MultiValued)
}

// TestSAMLConnectionClientUpdate_WithMultiValuedCustomAttributes verifies the client
// can update a SAML connection with custom attributes that include the `multi_valued`
// flag, covering both true and false values.
func TestSAMLConnectionClientUpdate_WithMultiValuedCustomAttributes(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"department","key":"department","sso_path":"$.department","scim_path":"department","multi_valued":false}]}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","multi_valued":true},{"name":"department","key":"department","sso_path":"$.department","scim_path":"department","multi_valued":false}]}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
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
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.NotNil(t, samlConnection.CustomAttributes)
	require.Equal(t, 2, len(*samlConnection.CustomAttributes))
	require.Equal(t, true, *(*samlConnection.CustomAttributes)[0].MultiValued)
	require.Equal(t, false, *(*samlConnection.CustomAttributes)[1].MultiValued)
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
	"data": [{"id":"%s","name":"%s","domain":"%s","domains": ["%s"],"provider":"%s"}],
	"total_count": 1
}`, id, name, domain, domain, provider)),
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
	// nolint:staticcheck // we want to test the .Domain behavior when it's deprecated
	require.Equal(t, domain, list.SAMLConnections[0].Domain)
	require.Equal(t, domain, list.SAMLConnections[0].Domains[0])
	require.Equal(t, provider, list.SAMLConnections[0].Provider)
}
