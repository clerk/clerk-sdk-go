package clerk

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
)

const (
	clerkBaseUrl = "https://api.clerk.dev"
)

type client struct {
	client *http.Client

	BaseURL *url.URL
}

func NewClient(token string) (*client, error) {
	baseURL, _ := url.Parse(clerkBaseUrl)
	ctx := context.Background()
	httpClient := createTokenClient(ctx, token)

	client := &client{client: httpClient, BaseURL: baseURL}
	return client, nil
}

func createTokenClient(ctx context.Context, token string) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return oauth2.NewClient(ctx, ts)
}
