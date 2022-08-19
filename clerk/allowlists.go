package clerk

import "net/http"

type AllowlistsService service

type AllowlistIdentifierResponse struct {
	Object       string `json:"object"`
	ID           string `json:"id"`
	InvitationID string `json:"invitation_id,omitempty"`
	Identifier   string `json:"identifier"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type CreateAllowlistIdentifierParams struct {
	Identifier string `json:"identifier"`
	Notify     bool   `json:"notify"`
}

func (s *AllowlistsService) CreateIdentifier(params CreateAllowlistIdentifierParams) (*AllowlistIdentifierResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPost, AllowlistsUrl, &params)

	var allowlistIdentifierResponse AllowlistIdentifierResponse
	_, err := s.client.Do(req, &allowlistIdentifierResponse)
	if err != nil {
		return nil, err
	}
	return &allowlistIdentifierResponse, nil
}

func (s *AllowlistsService) DeleteIdentifier(identifierID string) (*DeleteResponse, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, AllowlistsUrl+"/"+identifierID)

	var deleteResponse DeleteResponse
	_, err := s.client.Do(req, &deleteResponse)
	if err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}

type AllowlistIdentifiersResponse struct {
	Data       []*AllowlistIdentifierResponse `json:"data"`
	TotalCount int64                          `json:"total_count"`
}

func (s *AllowlistsService) ListAllIdentifiers() (*AllowlistIdentifiersResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, AllowlistsUrl)

	var allowlistIdentifiersResponse []*AllowlistIdentifierResponse
	_, err := s.client.Do(req, &allowlistIdentifiersResponse)
	if err != nil {
		return nil, err
	}
	return &AllowlistIdentifiersResponse{
		Data:       allowlistIdentifiersResponse,
		TotalCount: int64(len(allowlistIdentifiersResponse)),
	}, nil
}
