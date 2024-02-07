package instancesettings

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestInstanceClientUpdate(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
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
	err := client.Update(context.Background(), &UpdateParams{
		TestMode: clerk.Bool(true),
	})
	require.NoError(t, err)
}

func TestInstanceClientUpdate_Error(t *testing.T) {
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
	err := client.Update(context.Background(), &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestInstanceClientUpdateRestrictions(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
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
	restrictions, err := client.UpdateRestrictions(context.Background(), &UpdateRestrictionsParams{
		Allowlist: clerk.Bool(true),
	})
	require.NoError(t, err)
	require.True(t, restrictions.Allowlist)
	require.False(t, restrictions.Blocklist)
}

func TestInstanceClientUpdateRestrictions_Error(t *testing.T) {
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
	_, err := client.UpdateRestrictions(context.Background(), &UpdateRestrictionsParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestInstanceClientUpdateOrganizationSettings(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"enabled":true}`),
			Out:    json.RawMessage(`{"enabled":true,"max_allowed_memberships":3}`),
			Method: http.MethodPatch,
			Path:   "/v1/instance/organization_settings",
		},
	}
	client := NewClient(config)
	orgSettings, err := client.UpdateOrganizationSettings(context.Background(), &UpdateOrganizationSettingsParams{
		Enabled: clerk.Bool(true),
	})
	require.NoError(t, err)
	require.True(t, orgSettings.Enabled)
	require.Equal(t, int64(3), orgSettings.MaxAllowedMemberships)
}

func TestInstanceClientUpdateOrganizationSettings_Error(t *testing.T) {
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
	_, err := client.UpdateOrganizationSettings(context.Background(), &UpdateOrganizationSettingsParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}
