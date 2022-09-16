package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ActorTokenService service

type ActorTokenResponse struct {
	Object    string          `json:"object"`
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Actor     json.RawMessage `json:"actor"`
	Token     string          `json:"token,omitempty"`
	Status    string          `json:"status"`
	CreatedAt int64           `json:"created_at"`
	UpdatedAt int64           `json:"updated_at"`
}

type CreateActorTokenParams struct {
	UserID                      string          `json:"user_id"`
	Actor                       json.RawMessage `json:"actor"`
	ExpiresInSeconds            *int            `json:"expires_in_seconds"`
	SessionMaxDurationInSeconds *int            `json:"session_max_duration_in_seconds"`
}

func (s *ActorTokenService) Create(params CreateActorTokenParams) (*ActorTokenResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPost, ActorTokensUrl, &params)

	var actorTokenResponse ActorTokenResponse
	_, err := s.client.Do(req, &actorTokenResponse)
	if err != nil {
		return nil, err
	}
	return &actorTokenResponse, nil
}

func (s *ActorTokenService) Revoke(actorTokenID string) (*ActorTokenResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/revoke", ActorTokensUrl, actorTokenID))

	var actorTokenResponse ActorTokenResponse
	_, err := s.client.Do(req, &actorTokenResponse)
	if err != nil {
		return nil, err
	}
	return &actorTokenResponse, nil
}
