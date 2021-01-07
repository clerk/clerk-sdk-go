package clerk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	clerkBaseUrl = "https://api.clerk.dev"
)

type Client interface {
	NewRequest(method string, url string) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)
}

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
// A relative URL can be specified which is resolved relative to the BaseURL of the client.
// Relative URLs should be specified without a preceding slash.
func (c *client) NewRequest(method string, url string) (*http.Request, error) {
	fullUrl, err := c.BaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fullUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do will send the given request using the client `c` on which it is called.
// If the response contains a body, it will be unmarshalled in `v`.
func (c *client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkForErrors(resp)
	if err != nil {
		return resp, err
	}

	if resp.Body != nil && v != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		err = json.Unmarshal(body, &v)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func checkForErrors(resp *http.Response) error {
	if c := resp.StatusCode; c >= 200 && c < 400 {
		return nil
	}

	data, _ := ioutil.ReadAll(resp.Body)
	if data != nil && len(data) > 0 {
		return errors.New(string(data))
	}
	return errors.New(fmt.Sprintf("Server returned unexpected error with status code %d", resp.StatusCode))
}
