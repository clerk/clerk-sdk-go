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

// NewClient creates a new Clerk client.
// Because the token supplied will be used for all authenticated requests,
// the created client should not be used across different users
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

// NewRequest creates an API request.
// The urlStr is a URL which is resolved relative to the BaseURL of the client.
func (c *client) NewRequest(method string, urlStr string) (*http.Request, error) {
	fullUrl, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fullUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
