// Package machine provides the Machines API.
package machine

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/machines"

// Client is used to invoke the Machines API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type CreateParams struct {
	clerk.APIParams
	Name            string   `json:"name"`
	ScopedMachines  []string `json:"scoped_machines,omitempty"`
	DefaultTokenTTL *int64   `json:"default_token_ttl,omitempty"`
}

// Create creates a new machine.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.MachineWithScopedMachinesAndSecretKey, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	machine := &clerk.MachineWithScopedMachinesAndSecretKey{}
	err := c.Backend.Call(ctx, req, machine)
	return machine, err
}

// Get retrieves details for a machine.
func (c *Client) Get(ctx context.Context, id string) (*clerk.MachineWithScopedMachines, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	machine := &clerk.MachineWithScopedMachines{}
	err = c.Backend.Call(ctx, req, machine)
	return machine, err
}

// GetSecretKey retrieves the secret key for a machine.
func (c *Client) GetSecretKey(ctx context.Context, id string) (*clerk.MachineSecretKey, error) {
	path, err := clerk.JoinPath(path, id, "secret_key")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	secretKey := &clerk.MachineSecretKey{}
	err = c.Backend.Call(ctx, req, secretKey)
	return secretKey, err
}

type UpdateParams struct {
	clerk.APIParams
	Name            *string `json:"name,omitempty"`
	DefaultTokenTTL *int64  `json:"default_token_ttl,omitempty"`
}

// Update updates a machine.
func (c *Client) Update(ctx context.Context, id string, params *UpdateParams) (*clerk.MachineWithScopedMachines, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)
	machine := &clerk.MachineWithScopedMachines{}
	err = c.Backend.Call(ctx, req, machine)
	return machine, err
}

// Delete deletes a machine.
func (c *Client) Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	deletedResource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, deletedResource)
	return deletedResource, err
}

type ListParams struct {
	clerk.APIParams
	clerk.ListParams
	OrderBy *string `json:"order_by,omitempty"`
	Query   *string `json:"query,omitempty"`
}

func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}
	return q
}

// List returns a list of machines.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.MachineList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	machineList := &clerk.MachineList{}
	err := c.Backend.Call(ctx, req, machineList)
	return machineList, err
}
