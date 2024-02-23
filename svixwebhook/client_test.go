package svixwebhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestSvixWebhookClientCreate(t *testing.T) {
	t.Parallel()
	svixURL := "https://foo.com/webhook"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"svix_url":"%s"}`, svixURL)),
			Method: http.MethodPost,
			Path:   "/v1/webhooks/svix",
		},
	}
	client := NewClient(config)
	webhook, err := client.Create(context.Background())
	require.NoError(t, err)
	require.Equal(t, svixURL, webhook.SvixURL)
}

func TestSvixWebhookClientDelete(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Method: http.MethodDelete,
			Path:   "/v1/webhooks/svix",
		},
	}
	client := NewClient(config)
	_, err := client.Delete(context.Background())
	require.NoError(t, err)
}

func TestSvixWebhookClientRefreshURL(t *testing.T) {
	t.Parallel()
	svixURL := "https://foo.com/webhook"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"svix_url":"%s"}`, svixURL)),
			Method: http.MethodPost,
			Path:   "/v1/webhooks/svix_url",
		},
	}
	client := NewClient(config)
	webhook, err := client.RefreshURL(context.Background())
	require.NoError(t, err)
	require.Equal(t, svixURL, webhook.SvixURL)
}
