package clerk

import (
	"fmt"
	"net/http"
)

type ImpersonationTokenService service

type ImpersonationTokenRespose struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	SubjectID string `json:"subject_id"`
	ActorID   string `json:"actor_id"`
	Token     string `json:"token,omitempty"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CreateImpersonationTokenParams struct {
	SubjectID        string `json:"subject_id"`
	ActorID          string `json:"actor_id"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
}

func (s *ImpersonationTokenService) Create(params CreateImpersonationTokenParams) (*ImpersonationTokenRespose, error) {
	req, _ := s.client.NewRequest(http.MethodPost, ImpersonationTokensUrl, &params)

	var impersonationTokenResponse ImpersonationTokenRespose
	_, err := s.client.Do(req, &impersonationTokenResponse)
	if err != nil {
		return nil, err
	}
	return &impersonationTokenResponse, nil
}

func (s *ImpersonationTokenService) Revoke(impersonationTokenID string) (*ImpersonationTokenRespose, error) {
	req, _ := s.client.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/revoke", ImpersonationTokensUrl, impersonationTokenID))

	var impersonationTokenResponse ImpersonationTokenRespose
	_, err := s.client.Do(req, &impersonationTokenResponse)
	if err != nil {
		return nil, err
	}
	return &impersonationTokenResponse, nil
}
