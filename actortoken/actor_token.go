// Package actortoken provides the Actor Tokens API.
package actortoken

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

const path = "/actor_tokens"

type CreateParams struct {
	clerk.APIParams
	UserID                      *string         `json:"user_id,omitempty"`
	Actor                       json.RawMessage `json:"actor,omitempty"`
	ExpiresInSeconds            *int64          `json:"expires_in_seconds,omitempty"`
	SessionMaxDurationInSeconds *int64          `json:"session_max_duration_in_seconds,omitempty"`
}

// Create creates a new actor token.
func Create(ctx context.Context, params *CreateParams) (*clerk.ActorToken, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	token := &clerk.ActorToken{}
	err := clerk.GetBackend().Call(ctx, req, token)
	return token, err
}

// Revoke revokes a pending actor token.
func Revoke(ctx context.Context, id string) (*clerk.ActorToken, error) {
	token := &clerk.ActorToken{}
	path, err := clerk.JoinPath(path, id, "revoke")
	if err != nil {
		return token, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	err = clerk.GetBackend().Call(ctx, req, token)
	return token, err
}
