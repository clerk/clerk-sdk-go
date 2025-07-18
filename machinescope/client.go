// Package machinescope provides the Machine Scopes API.
package machinescope

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/machines"

// Client is used to invoke the Machine Scopes API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type CreateScopeParams struct {
	clerk.APIParams
	ToMachineID string `json:"to_machine_id"`
}

// CreateScope creates a new machine scope.
func (c *Client) CreateScope(ctx context.Context, machineID string, params *CreateScopeParams) (*clerk.MachineScope, error) {
	path, err := clerk.JoinPath(path, machineID, "scopes")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	machineScope := &clerk.MachineScope{}
	err = c.Backend.Call(ctx, req, machineScope)
	return machineScope, err
}

// DeleteScope deletes a machine scope.
func (c *Client) DeleteScope(ctx context.Context, machineID string, otherMachineID string) (*clerk.DeletedMachineScope, error) {
	path, err := clerk.JoinPath(path, machineID, "scopes", otherMachineID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	deletedMachineScope := &clerk.DeletedMachineScope{}
	err = c.Backend.Call(ctx, req, deletedMachineScope)
	return deletedMachineScope, err
}
