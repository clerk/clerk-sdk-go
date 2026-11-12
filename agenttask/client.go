// Package agenttask provides the Agent Tasks API.
package agenttask

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/agents/tasks"

// Client is used to invoke the Agent Tasks API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

// OnBehalfOf identifies the user an agent task is created for. Provide
// exactly one of UserID or Identifier.
type OnBehalfOf struct {
	UserID *string `json:"user_id,omitempty"`
	// Identifier is a verified identifier of the user, such as an email
	// address.
	Identifier *string `json:"identifier,omitempty"`
}

type CreateParams struct {
	clerk.APIParams
	OnBehalfOf *OnBehalfOf `json:"on_behalf_of,omitempty"`
	// Permissions must be "*", which grants all permissions.
	Permissions *string `json:"permissions,omitempty"`
	// AgentName identifies the agent. Tasks with the same agent name get the
	// same AgentID.
	AgentName       *string `json:"agent_name,omitempty"`
	TaskDescription *string `json:"task_description,omitempty"`
	// RedirectURL is where the user lands after visiting the task URL.
	RedirectURL                 *string `json:"redirect_url,omitempty"`
	SessionMaxDurationInSeconds *int64  `json:"session_max_duration_in_seconds,omitempty"`
	// SessionRestricted makes Clerk's Frontend API refuse the session's
	// writes.
	SessionRestricted *bool `json:"session_restricted,omitempty"`
	// OrganizationID is the ID of the organization the session starts in.
	// The user must be a member of it. If nil, the session starts in the
	// user's last active organization, or in the personal account when the
	// user has none.
	OrganizationID *string `json:"session_organization_id,omitempty"`
}

// Create creates an agent task.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.AgentTask, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	task := &clerk.AgentTask{}
	err := c.Backend.Call(ctx, req, task)
	return task, err
}

// Revoke revokes a pending agent task.
func (c *Client) Revoke(ctx context.Context, id string) (*clerk.AgentTask, error) {
	task := &clerk.AgentTask{}
	path, err := clerk.JoinPath(path, id, "revoke")
	if err != nil {
		return task, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	err = c.Backend.Call(ctx, req, task)
	return task, err
}
