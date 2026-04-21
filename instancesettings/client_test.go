package instancesettings

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestInstanceClientUpdate(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"test_mode":true}`),
			Out:    nil,
			Method: http.MethodPatch,
			Path:   "/v1/instance",
		},
	}
	client := NewClient(config)
	err := client.Update(t.Context(), &UpdateParams{
		TestMode: clerk.Bool(true),
	})
	require.NoError(t, err)
}

func TestInstanceClientUpdate_Error(t *testing.T) {
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
	err := client.Update(t.Context(), &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestInstanceClientUpdateRestrictions(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"allowlist":true}`),
			Out:    json.RawMessage(`{"allowlist":true,"blocklist":false}`),
			Method: http.MethodPatch,
			Path:   "/v1/instance/restrictions",
		},
	}
	client := NewClient(config)
	restrictions, err := client.UpdateRestrictions(t.Context(), &UpdateRestrictionsParams{
		Allowlist: clerk.Bool(true),
	})
	require.NoError(t, err)
	require.True(t, restrictions.Allowlist)
	require.False(t, restrictions.Blocklist)
}

func TestInstanceClientUpdateRestrictions_Error(t *testing.T) {
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
	_, err := client.UpdateRestrictions(t.Context(), &UpdateRestrictionsParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestInstanceClientUpdateOrganizationSettings(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"enabled":true,"max_allowed_domains":20,"force_organization_selection": true,"slug_disabled": true}`),
			Out:    json.RawMessage(`{"enabled":true,"max_allowed_memberships":3,"max_allowed_domains":20}`),
			Method: http.MethodPatch,
			Path:   "/v1/instance/organization_settings",
		},
	}
	client := NewClient(config)
	orgSettings, err := client.UpdateOrganizationSettings(t.Context(), &UpdateOrganizationSettingsParams{
		Enabled:                    clerk.Bool(true),
		MaxAllowedDomains:          clerk.Int64(20),
		ForceOrganizationSelection: clerk.Bool(true),
		SlugDisabled:               clerk.Bool(true),
	})
	require.NoError(t, err)
	require.True(t, orgSettings.Enabled)
	require.Equal(t, int64(3), orgSettings.MaxAllowedMemberships)
	require.Equal(t, int64(20), orgSettings.MaxAllowedDomains)
}

func TestInstanceClientUpdateOrganizationSettingsWithInitialRoleSetKey(t *testing.T) {
	t.Parallel()
	initialRoleSetKey := "admin-roles"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"enabled":true,"initial_role_set_key":"admin-roles"}`),
			Out:    json.RawMessage(`{"enabled":true,"max_allowed_memberships":5,"initial_role_set_key":"admin-roles", "max_role_sets_allowed": 30}`),
			Method: http.MethodPatch,
			Path:   "/v1/instance/organization_settings",
		},
	}
	client := NewClient(config)
	orgSettings, err := client.UpdateOrganizationSettings(t.Context(), &UpdateOrganizationSettingsParams{
		Enabled:           clerk.Bool(true),
		InitialRoleSetKey: clerk.String(initialRoleSetKey),
	})
	require.NoError(t, err)
	require.True(t, orgSettings.Enabled)
	require.Equal(t, int64(5), orgSettings.MaxAllowedMemberships)
	require.NotNil(t, orgSettings.InitialRoleSetKey)
	require.Equal(t, initialRoleSetKey, *orgSettings.InitialRoleSetKey)
}

func TestInstanceClientUpdateOrganizationSettings_Error(t *testing.T) {
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
	_, err := client.UpdateOrganizationSettings(t.Context(), &UpdateOrganizationSettingsParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestInstanceClientUpdateOrganizationSettingsWithCreationDefaults(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"enabled":true,"organization_creation_defaults":{"enabled":true,"automatic_organization_creation":{"enabled":true},"detect_from_email_domain":{"enabled":false},"organization_name_template":{"enabled":true,"template":"{{user.first_name}}'s Organization"},"fallback":{"name":"My Organization"}}}`),
			Out:    json.RawMessage(`{"enabled":true,"max_allowed_memberships":3,"organization_creation_defaults":{"enabled":true,"automatic_organization_creation":{"enabled":true},"detect_from_email_domain":{"enabled":false},"organization_name_template":{"enabled":true,"template":"{{user.first_name}}'s Organization"},"fallback":{"name":"My Organization"}}}`),
			Method: http.MethodPatch,
			Path:   "/v1/instance/organization_settings",
		},
	}
	client := NewClient(config)
	orgSettings, err := client.UpdateOrganizationSettings(t.Context(), &UpdateOrganizationSettingsParams{
		Enabled: clerk.Bool(true),
		OrganizationCreationDefaults: &UpdateOrganizationCreationDefaultsParams{
			Enabled: clerk.Bool(true),
			AutomaticOrganizationCreation: &AutomaticOrganizationCreationSettingsParams{
				Enabled: clerk.Bool(true),
			},
			DetectFromEmailDomain: &DetectFromEmailDomainSettingsParams{
				Enabled: clerk.Bool(false),
			},
			OrganizationNameTemplate: &OrganizationNameTemplateSettingsParams{
				Enabled:  clerk.Bool(true),
				Template: clerk.String("{{user.first_name}}'s Organization"),
			},
			Fallback: &FallbackSettingsParams{
				Name: clerk.String("My Organization"),
			},
		},
	})
	require.NoError(t, err)
	require.True(t, orgSettings.Enabled)
	require.Equal(t, int64(3), orgSettings.MaxAllowedMemberships)
	require.NotNil(t, orgSettings.OrganizationCreationDefaults)
	require.True(t, orgSettings.OrganizationCreationDefaults.Enabled)
	require.True(t, orgSettings.OrganizationCreationDefaults.AutomaticOrganizationCreation.Enabled)
	require.False(t, orgSettings.OrganizationCreationDefaults.DetectFromEmailDomain.Enabled)
	require.True(t, orgSettings.OrganizationCreationDefaults.OrganizationNameTemplate.Enabled)
	require.Equal(t, "{{user.first_name}}'s Organization", orgSettings.OrganizationCreationDefaults.OrganizationNameTemplate.Template)
	require.Equal(t, "My Organization", orgSettings.OrganizationCreationDefaults.Fallback.Name)
}

func TestInstanceClientReadOAuthApplicationSettings(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"object":"oauth_application_settings","dynamic_oauth_client_registration":false,"oauth_jwt_access_tokens":true,"oidc_sign_out_enabled":true}`),
			Method: http.MethodGet,
			Path:   "/v1/instance/oauth_application_settings",
		},
	}
	client := NewClient(config)
	settings, err := client.ReadOAuthApplicationSettings(t.Context())
	require.NoError(t, err)
	require.Equal(t, "oauth_application_settings", settings.Object)
	require.False(t, settings.DynamicOauthClientRegistration)
	require.True(t, settings.OAuthJWTAccessTokens)
	require.True(t, settings.OIDCSignOutEnabled)
}

func TestInstanceClientUpdateOAuthApplicationSettings(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"dynamic_oauth_client_registration":true,"oauth_jwt_access_tokens":false,"oidc_sign_out_enabled":true}`),
			Out:    json.RawMessage(`{"object":"oauth_application_settings","dynamic_oauth_client_registration":true,"oauth_jwt_access_tokens":false,"oidc_sign_out_enabled":true}`),
			Method: http.MethodPatch,
			Path:   "/v1/instance/oauth_application_settings",
		},
	}
	client := NewClient(config)
	settings, err := client.UpdateOAuthApplicationSettings(t.Context(), &UpdateOAuthApplicationSettingsParams{
		DynamicOAuthClientRegistration: clerk.Bool(true),
		OAuthJWTAccessTokens:           clerk.Bool(false),
		OIDCSignOutEnabled:             clerk.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, "oauth_application_settings", settings.Object)
	require.True(t, settings.DynamicOauthClientRegistration)
	require.False(t, settings.OAuthJWTAccessTokens)
	require.True(t, settings.OIDCSignOutEnabled)
}

func TestInstanceClientUpdateOAuthApplicationSettings_Error(t *testing.T) {
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
	_, err := client.UpdateOAuthApplicationSettings(t.Context(), &UpdateOAuthApplicationSettingsParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestInstanceClientGetOrganizationSettings(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"enabled":true,"max_allowed_memberships":3}`),
			Method: http.MethodGet,
			Path:   "/v1/instance/organization_settings",
		},
	}
	client := NewClient(config)
	orgSettings, err := client.GetOrganizationSettings(t.Context())
	require.NoError(t, err)
	require.True(t, orgSettings.Enabled)
	require.Equal(t, int64(3), orgSettings.MaxAllowedMemberships)
}

func TestInstanceClientGetOrganizationSettings_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"get-error-code"
	}],
	"clerk_trace_id":"get-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.GetOrganizationSettings(t.Context())
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "get-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "get-error-code", apiErr.Errors[0].Code)
}
