package clerk

import (
	"errors"
	"net/http"
)

// ClerkOption describes a functional parameter for the clerk client constructor
type ClerkOption func(*client) error

// WithHTTPClient allows the overriding of the http client
func WithHTTPClient(httpClient *http.Client) ClerkOption {
	return func(c *client) error {
		if httpClient == nil {
			return errors.New("http client can't be nil")
		}

		c.client = httpClient
		return nil
	}
}

// WithBaseURL allows the overriding of the base URL
func WithBaseURL(rawURL string) ClerkOption {
	return func(c *client) error {
		if rawURL == "" {
			return errors.New("base url can't be empty")
		}

		baseURL, err := toURLWithEndingSlash(rawURL)
		if err != nil {
			return err
		}

		c.baseURL = baseURL
		return nil
	}
}
