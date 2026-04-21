package oauthapplication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/clerk/clerk-sdk-go/v3/optional"
	"github.com/stretchr/testify/require"
)

func TestOAuthApplicationClientCreate(t *testing.T) {
	t.Parallel()
	id := "oauth_app_123"
	name := "Test Application"
	callbackURL := "https://callback.url"
	scopes := "read,write"
	public := true
	clientSecret := "secret_123"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","callback_url":"https://callback.url","scopes":"read,write","public":true,"consent_screen_enabled":true,"pkce_required":true}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","callback_url":"%s","scopes":"%s","public":%t,"client_secret":"%s","consent_screen_enabled":true,"pkce_required":true}`, id, name, callbackURL, scopes, public, clientSecret)),
			Method: http.MethodPost,
			Path:   "/v1/oauth_applications",
		},
	}
	client := NewClient(config)

	params := &CreateParams{
		Name:                 name,
		CallbackURL:          clerk.String(callbackURL),
		Scopes:               clerk.String(scopes),
		Public:               clerk.Bool(public),
		ConsentScreenEnabled: clerk.Bool(true),
		PKCERequired:         clerk.Bool(true),
	}
	oauthApp, err := client.Create(t.Context(), params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, name, oauthApp.Name)
	//nolint:staticcheck
	require.Equal(t, callbackURL, oauthApp.CallbackURL)
	require.Equal(t, scopes, oauthApp.Scopes)
	require.Equal(t, public, oauthApp.Public)
	require.Equal(t, clientSecret, *oauthApp.ClientSecret)
	require.Equal(t, true, oauthApp.ConsentScreenEnabled)
	require.Equal(t, true, oauthApp.PKCERequired)
}

func TestOrganizationClientCreate_Error(t *testing.T) {
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
	_, err := client.Create(t.Context(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestOAuthApplicationClientGet(t *testing.T) {
	t.Parallel()
	id := "app_123"
	name := "Test Application"
	callbackURL := "https://callback.url"
	scopes := "read,write"
	public := true

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","callback_url":"%s","scopes":"%s","public":%t}`, id, name, callbackURL, scopes, public)),
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/v1/oauth_applications/%s", id),
		},
	}

	client := NewClient(config)
	oauthApp, err := client.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, name, oauthApp.Name)
	//nolint:staticcheck
	require.Equal(t, callbackURL, oauthApp.CallbackURL)
	require.Equal(t, scopes, oauthApp.Scopes)
	require.Equal(t, public, oauthApp.Public)
}

func TestOAuthApplicationClientUpdate(t *testing.T) {
	t.Parallel()
	id := "app_123"
	updatedName := "Updated Application"
	callbackURL1 := "https://updated.callback.url"
	callbackURL2 := "https://new.callback.url"
	public := true

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","redirect_uris":["%s","%s"], "public":%t,"consent_screen_enabled":%t,"pkce_required":true}`, updatedName, callbackURL1, callbackURL2, public, false)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","redirect_uris":["%s","%s"],"callback_url":"%s","public":%t,"consent_screen_enabled":%t,"pkce_required":true}`, id, updatedName, callbackURL1, callbackURL2, callbackURL1, public, false)),
			Method: http.MethodPatch,
			Path:   fmt.Sprintf("/v1/oauth_applications/%s", id),
		},
	}

	client := NewClient(config)
	params := &UpdateParams{
		Name:                 clerk.String(updatedName),
		RedirectURIs:         []string{callbackURL1, callbackURL2},
		Public:               clerk.Bool(public),
		ConsentScreenEnabled: clerk.Bool(false),
		PKCERequired:         clerk.Bool(true),
	}
	oauthApp, err := client.Update(t.Context(), id, params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, updatedName, oauthApp.Name)
	//nolint:staticcheck
	require.Equal(t, callbackURL1, oauthApp.CallbackURL)
	require.Equal(t, []string{callbackURL1, callbackURL2}, oauthApp.RedirectURIs)
	require.False(t, oauthApp.ConsentScreenEnabled)
	require.Equal(t, true, oauthApp.PKCERequired)
	require.Nil(t, oauthApp.AccessTokenTTL)
}

func TestOAuthApplicationClientUpdate_AccessTokenTTL(t *testing.T) {
	t.Parallel()
	id := "app_123"
	updatedTTL := int64(3600)
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"access_token_ttl":%d}`, updatedTTL)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","access_token_ttl":%d}`, id, updatedTTL)),
			Method: http.MethodPatch,
			Path:   fmt.Sprintf("/v1/oauth_applications/%s", id),
		},
	}
	client := NewClient(config)
	params := &UpdateParams{
		AccessTokenTTL: optional.Int64FromPtr(&updatedTTL),
	}
	oauthApp, err := client.Update(t.Context(), id, params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, updatedTTL, *oauthApp.AccessTokenTTL)
}

// Demonstrates setting access token TTL to null to reset to default.
func TestOAuthApplicationClientUpdate_AccessTokenTTL_ResetToDefault_Null(t *testing.T) {
	t.Parallel()
	id := "app_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage((`{"access_token_ttl":null}`)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","access_token_ttl":null}`, id)),
			Method: http.MethodPatch,
			Path:   fmt.Sprintf("/v1/oauth_applications/%s", id),
		},
	}
	client := NewClient(config)
	params := &UpdateParams{
		// Specify null to reset to default
		AccessTokenTTL: optional.Int64FromPtr(nil),
	}
	oauthApp, err := client.Update(t.Context(), id, params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Nil(t, oauthApp.AccessTokenTTL)
}

// Demonstrates that omitting access token TTL does not pass the field to the API, and leaves it unchanged.
func TestOAuthApplicationClientUpdate_AccessTokenTTL_Omit(t *testing.T) {
	t.Parallel()
	id := "app_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage((`{"name":"Test Application"}`)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"Test Application","access_token_ttl":42}`, id)),
			Method: http.MethodPatch,
			Path:   fmt.Sprintf("/v1/oauth_applications/%s", id),
		},
	}
	client := NewClient(config)
	params := &UpdateParams{
		Name: clerk.String("Test Application"),
		// Omit access token TTL
	}
	oauthApp, err := client.Update(t.Context(), id, params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, int64(42), *oauthApp.AccessTokenTTL)
	require.Equal(t, "Test Application", oauthApp.Name)
}

func TestOrganizationClientUpdate_Error(t *testing.T) {
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
	_, err := client.Update(t.Context(), "oauth_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestOAuthApplicationClientDelete(t *testing.T) {
	t.Parallel()
	id := "oauth_app_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/oauth_applications/" + id,
		},
	}
	client := NewClient(config)
	deletedApp, err := client.DeleteOAuthApplication(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, id, deletedApp.ID)
	require.True(t, deletedApp.Deleted)
}

func TestOAuthApplicationClientList(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
				"data": [{"id":"oauth_app_123","name":"OAuth App 1"}],
				"total_count": 1
			}`),
			Method: http.MethodGet,
			Path:   "/v1/oauth_applications",
			Query: &url.Values{
				"limit":      []string{"10"},
				"offset":     []string{"0"},
				"order_by":   []string{"-name"},
				"name_query": []string{"app_123"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{}
	params.Limit = clerk.Int64(10)
	params.Offset = clerk.Int64(0)
	params.OrderBy = clerk.String("-name")
	params.NameQuery = clerk.String("app_123")
	appList, err := client.List(t.Context(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), appList.TotalCount)
	require.Equal(t, "oauth_app_123", appList.OAuthApplications[0].ID)
	require.Equal(t, "OAuth App 1", appList.OAuthApplications[0].Name)
}

func TestOAuthApplicationClientRotateClientSecret(t *testing.T) {
	t.Parallel()
	id := "oauth_app_123"
	newSecret := "new_client_secret"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","client_secret":"%s"}`, id, newSecret)),
			Method: http.MethodPost,
			Path:   "/v1/oauth_applications/" + id + "/rotate_secret",
		},
	}
	client := NewClient(config)
	updatedApp, err := client.RotateClientSecret(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, id, updatedApp.ID)
	require.Equal(t, newSecret, *updatedApp.ClientSecret)
}
