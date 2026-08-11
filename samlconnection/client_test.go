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
	loginHintMode := "custom_attribute"
	loginHintSource := "employee_id"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s","provider":"%s","force_authn":%t,"login_hint":{"mode":"%s","source":"%s"}}`, name, domain, provider, forceAuthn, loginHintMode, loginHintSource)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","force_authn":%t,"login_hint":{"mode":"%s","source":"%s"}}`, id, name, domain, provider, forceAuthn, loginHintMode, loginHintSource)),
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
		LoginHint: &LoginHintParams{
			Mode:   clerk.String(loginHintMode),
			Source: clerk.String(loginHintSource),
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, name, samlConnection.Name)
	// nolint:staticcheck // we want to test the .Domain behavior when it's deprecated
	require.Equal(t, domain, samlConnection.Domain)
	require.Equal(t, provider, samlConnection.Provider)
	require.Equal(t, forceAuthn, samlConnection.ForceAuthn)
	require.NotNil(t, samlConnection.LoginHint)
	require.Equal(t, loginHintMode, samlConnection.LoginHint.Mode)
	require.NotNil(t, samlConnection.LoginHint.Source)
	require.Equal(t, loginHintSource, *samlConnection.LoginHint.Source)
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
	loginHintMode := "off"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s", "disable_additional_identifications": %t, "force_authn": %t, "login_hint": {"mode": "%s"}, "enterprise_connection_id": "entconn_2abc123def456", "custom_attributes": [{"name": "custom_attribute_name", "key": "custom_attribute_key", "sso_path": "custom_attribute_sso_path", "directory_path": "custom_attribute_directory_path"}]}`, id, name, domain, provider, disableAdditionalIdentifications, forceAuthn, loginHintMode)),
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
	require.NotNil(t, samlConnection.LoginHint)
	require.Equal(t, loginHintMode, samlConnection.LoginHint.Mode)
	require.Nil(t, samlConnection.LoginHint.Source)
	// The response carries the path as `directory_path`; it decodes into
	// SCIMPath, and DirectoryPath stays nil.
	require.Equal(t, &[]clerk.CustomAttribute{
		{
			Name:    clerk.String("custom_attribute_name"),
			Key:     clerk.String("custom_attribute_key"),
			SSOPath: clerk.String("custom_attribute_sso_path"),
			// nolint:staticcheck // SCIMPath is the field a decoded path lands in
			SCIMPath: clerk.String("custom_attribute_directory_path"),
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
	loginHintMode := "email_address"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s", "disable_additional_identifications": %t, "organization_id": "", "force_authn": %t, "login_hint": {"mode": "%s"}, "custom_attributes": [{"name": "custom_attribute_name", "key": "custom_attribute_key", "sso_path": "custom_attribute_sso_path", "directory_path": "custom_attribute_directory_path"}]}`, name, disableAdditionalIdentifications, forceAuthn, loginHintMode)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","disable_additional_identifications": %t,"force_authn": %t, "login_hint": {"mode": "%s"}, "enterprise_connection_id": null, "custom_attributes": [{"name": "custom_attribute_name", "key": "custom_attribute_key", "sso_path": "custom_attribute_sso_path", "directory_path": "custom_attribute_directory_path"}]}`, id, name, domain, provider, disableAdditionalIdentifications, forceAuthn, loginHintMode)),
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
		LoginHint: &LoginHintParams{
			Mode: clerk.String(loginHintMode),
		},
		CustomAttributes: &[]clerk.CustomAttribute{
			{
				Name:          clerk.String("custom_attribute_name"),
				Key:           clerk.String("custom_attribute_key"),
				SSOPath:       clerk.String("custom_attribute_sso_path"),
				DirectoryPath: clerk.String("custom_attribute_directory_path"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.Equal(t, &[]clerk.CustomAttribute{
		{
			Name:    clerk.String("custom_attribute_name"),
			Key:     clerk.String("custom_attribute_key"),
			SSOPath: clerk.String("custom_attribute_sso_path"),
			// nolint:staticcheck // SCIMPath is the field a decoded path lands in
			SCIMPath: clerk.String("custom_attribute_directory_path"),
		},
	}, samlConnection.CustomAttributes)
	require.Equal(t, name, samlConnection.Name)
	require.Equal(t, disableAdditionalIdentifications, samlConnection.DisableAdditionalIdentifications)
	require.Equal(t, forceAuthn, samlConnection.ForceAuthn)
	require.NotNil(t, samlConnection.LoginHint)
	require.Equal(t, loginHintMode, samlConnection.LoginHint.Mode)
	require.Nil(t, samlConnection.LoginHint.Source)
}

func TestSAMLConnectionClientUpdate_DisableJITProvisioning(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	disableJITProvisioning := true
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","disable_jit_provisioning":%t}`, name, disableJITProvisioning)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","disable_jit_provisioning":%t}`, id, name, domain, provider, disableJITProvisioning)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		Name:                   clerk.String(name),
		DisableJITProvisioning: clerk.Bool(disableJITProvisioning),
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	require.NotNil(t, samlConnection.DisableJITProvisioning)
	require.Equal(t, disableJITProvisioning, *samlConnection.DisableJITProvisioning)
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
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s","provider":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","directory_path":"groups","multi_valued":true},{"name":"manager","key":"manager","sso_path":"$.manager","directory_path":"manager","multi_valued":false}]}`, name, domain, provider)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","directory_path":"groups","multi_valued":true},{"name":"manager","key":"manager","sso_path":"$.manager","directory_path":"manager","multi_valued":false}]}`, id, name, domain, provider)),
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
				Name:          clerk.String("groups"),
				Key:           clerk.String("groups"),
				SSOPath:       clerk.String("$.groups"),
				DirectoryPath: clerk.String("groups"),
				MultiValued:   clerk.Bool(true),
			},
			{
				Name:          clerk.String("manager"),
				Key:           clerk.String("manager"),
				SSOPath:       clerk.String("$.manager"),
				DirectoryPath: clerk.String("manager"),
				MultiValued:   clerk.Bool(false),
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
			In:     json.RawMessage(`{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","directory_path":"groups","multi_valued":true},{"name":"department","key":"department","sso_path":"$.department","directory_path":"department","multi_valued":false}]}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","directory_path":"groups","multi_valued":true},{"name":"department","key":"department","sso_path":"$.department","directory_path":"department","multi_valued":false}]}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	client := NewClient(config)
	samlConnection, err := client.Update(context.Background(), id, &UpdateParams{
		CustomAttributes: &[]clerk.CustomAttribute{
			{
				Name:          clerk.String("groups"),
				Key:           clerk.String("groups"),
				SSOPath:       clerk.String("$.groups"),
				DirectoryPath: clerk.String("groups"),
				MultiValued:   clerk.Bool(true),
			},
			{
				Name:          clerk.String("department"),
				Key:           clerk.String("department"),
				SSOPath:       clerk.String("$.department"),
				DirectoryPath: clerk.String("department"),
				MultiValued:   clerk.Bool(false),
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

// TestSAMLConnectionClientCreate_WithLegacySCIMPath verifies that the legacy
// CustomAttribute.SCIMPath field still works. A client that sets only SCIMPath
// sends `scim_path` and no `directory_path`. The response carries both names;
// the path decodes into SCIMPath and DirectoryPath stays nil.
func TestSAMLConnectionClientCreate_WithLegacySCIMPath(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	name := "the-name"
	domain := "example.com"
	provider := "saml_custom"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","domain":"%s","provider":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups"}]}`, name, domain, provider)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","domain":"%s","provider":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"groups","directory_path":"groups"}]}`, id, name, domain, provider)),
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
				Name:    clerk.String("groups"),
				Key:     clerk.String("groups"),
				SSOPath: clerk.String("$.groups"),
				// nolint:staticcheck // exercising the deprecated field on purpose
				SCIMPath: clerk.String("groups"),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, samlConnection.ID)
	attrs := *samlConnection.CustomAttributes
	require.Len(t, attrs, 1)
	// nolint:staticcheck // exercising the deprecated field on purpose
	require.Equal(t, "groups", *attrs[0].SCIMPath)
	require.Nil(t, attrs[0].DirectoryPath)
}

// TestSAMLConnectionClientUpdate_ReadModifyWriteSCIMPath covers the pattern the
// dual path names would otherwise break: read a connection, edit the path
// through the legacy SCIMPath field, and send the attributes straight back. The
// request must carry only `scim_path` with the edited value. Sending the stale
// `directory_path` alongside it is what the API rejects with a 422.
func TestSAMLConnectionClientUpdate_ReadModifyWriteSCIMPath(t *testing.T) {
	t.Parallel()
	id := "samlc__123"

	getConfig := &clerk.ClientConfig{}
	getConfig.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"old.path","directory_path":"old.path"}]}`, id)),
			Method: http.MethodGet,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	samlConnection, err := NewClient(getConfig).Get(context.Background(), id)
	require.NoError(t, err)

	attrs := *samlConnection.CustomAttributes
	// nolint:staticcheck // exercising the deprecated field on purpose
	attrs[0].SCIMPath = clerk.String("new.path")

	updateConfig := &clerk.ClientConfig{}
	updateConfig.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"new.path"}]}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"new.path","directory_path":"new.path"}]}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	updated, err := NewClient(updateConfig).Update(context.Background(), id, &UpdateParams{
		CustomAttributes: &attrs,
	})
	require.NoError(t, err)
	// nolint:staticcheck // exercising the deprecated field on purpose
	require.Equal(t, "new.path", *(*updated.CustomAttributes)[0].SCIMPath)
}

// TestSAMLConnectionClientUpdate_ClearPath verifies that clearing the path still
// clears it. Both path fields nil means neither name is sent, which the API
// reads as an empty path.
func TestSAMLConnectionClientUpdate_ClearPath(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups"}]}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"","directory_path":""}]}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	updated, err := NewClient(config).Update(context.Background(), id, &UpdateParams{
		CustomAttributes: &[]clerk.CustomAttribute{
			{
				Name:    clerk.String("groups"),
				Key:     clerk.String("groups"),
				SSOPath: clerk.String("$.groups"),
			},
		},
	})
	require.NoError(t, err)
	// nolint:staticcheck // exercising the deprecated field on purpose
	require.Equal(t, "", *(*updated.CustomAttributes)[0].SCIMPath)
}

// TestSAMLConnectionClientUpdate_DirectoryPathWins verifies that DirectoryPath
// takes precedence on write and is sent under the new name only, even when a
// stale SCIMPath is still set from a previous read.
func TestSAMLConnectionClientUpdate_DirectoryPathWins(t *testing.T) {
	t.Parallel()
	id := "samlc__123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","directory_path":"new.path"}]}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","custom_attributes":[{"name":"groups","key":"groups","sso_path":"$.groups","scim_path":"new.path","directory_path":"new.path"}]}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/saml_connections/" + id,
		},
	}
	_, err := NewClient(config).Update(context.Background(), id, &UpdateParams{
		CustomAttributes: &[]clerk.CustomAttribute{
			{
				Name:    clerk.String("groups"),
				Key:     clerk.String("groups"),
				SSOPath: clerk.String("$.groups"),
				// nolint:staticcheck // stale value from a previous read
				SCIMPath:      clerk.String("old.path"),
				DirectoryPath: clerk.String("new.path"),
			},
		},
	})
	require.NoError(t, err)
}
