package oauthapplication

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
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","callback_url":"https://callback.url","scopes":"read,write","public":true}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","callback_url":"%s","scopes":"%s","public":%t,"client_secret":"%s"}`, id, name, callbackURL, scopes, public, clientSecret)),
			Method: http.MethodPost,
			Path:   "/v1/oauth_applications",
		},
	}
	client := NewClient(config)

	params := &CreateParams{
		Name:        name,
		CallbackURL: callbackURL,
		Scopes:      scopes,
		Public:      public,
	}
	oauthApp, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, name, oauthApp.Name)
	require.Equal(t, callbackURL, oauthApp.CallbackURL)
	require.Equal(t, scopes, oauthApp.Scopes)
	require.Equal(t, public, oauthApp.Public)
	require.Equal(t, clientSecret, *oauthApp.ClientSecret)
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
	_, err := client.Create(context.Background(), &CreateParams{})
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
	oauthApp, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, name, oauthApp.Name)
	require.Equal(t, callbackURL, oauthApp.CallbackURL)
	require.Equal(t, scopes, oauthApp.Scopes)
	require.Equal(t, public, oauthApp.Public)
}

func TestOAuthApplicationClientUpdate(t *testing.T) {
	t.Parallel()
	id := "app_123"
	updatedName := "Updated Application"
	callbackURL := "https://updated.callback.url"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","callback_url":"%s","consent_screen_enabled":true}`, updatedName, callbackURL)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","callback_url":"%s","consent_screen_enabled":true}`, id, updatedName, callbackURL)),
			Method: http.MethodPatch,
			Path:   fmt.Sprintf("/v1/oauth_applications/%s", id),
		},
	}

	client := NewClient(config)
	params := &UpdateParams{
		Name:                 clerk.String(updatedName),
		CallbackURL:          clerk.String(callbackURL),
		ConsentScreenEnabled: clerk.Bool(true),
	}
	oauthApp, err := client.Update(context.Background(), id, params)
	require.NoError(t, err)
	require.Equal(t, id, oauthApp.ID)
	require.Equal(t, updatedName, oauthApp.Name)
	require.Equal(t, callbackURL, oauthApp.CallbackURL)
	require.Equal(t, true, oauthApp.ConsentScreenEnabled)
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
	_, err := client.Update(context.Background(), "oauth_123", &UpdateParams{})
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
	deletedApp, err := client.DeleteOAuthApplication(context.Background(), id)
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
				"limit":  []string{"10"},
				"offset": []string{"0"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{}
	params.Limit = clerk.Int64(10)
	params.Offset = clerk.Int64(0)
	appList, err := client.List(context.Background(), params)
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
	updatedApp, err := client.RotateClientSecret(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, updatedApp.ID)
	require.Equal(t, newSecret, *updatedApp.ClientSecret)
}
