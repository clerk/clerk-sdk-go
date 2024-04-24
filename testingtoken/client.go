// Package testingtoken provides the Testing Tokens API.
//
// https://clerk.com/docs/reference/backend-api/tag/Testing-Tokens
package testingtoken

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/testing_tokens"

// Client is used to invoke the Testing Tokens API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// Create creates a new testing token.
func (c *Client) Create(ctx context.Context) (*clerk.TestingToken, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	token := &clerk.TestingToken{}
	err := c.Backend.Call(ctx, req, token)
	return token, err
}
