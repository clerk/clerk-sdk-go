package clerk

import (
	"net/http"
)

type BlocklistsService service

type BlocklistIdentifierResponse struct {
	Object         string `json:"object"`
	ID             string `json:"id"`
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifier_type"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

type CreateBlocklistIdentifierParams struct {
	Identifier string `json:"identifier"`
}

func (s *BlocklistsService) CreateIdentifier(params CreateBlocklistIdentifierParams) (*BlocklistIdentifierResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPost, BlocklistsUrl, &params)

	var blocklistIdentifierResponse BlocklistIdentifierResponse
	_, err := s.client.Do(req, &blocklistIdentifierResponse)
	if err != nil {
		return nil, err
	}
	return &blocklistIdentifierResponse, nil
}

func (s *BlocklistsService) DeleteIdentifier(identifierID string) (*DeleteResponse, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, BlocklistsUrl+"/"+identifierID)

	var deleteResponse DeleteResponse
	_, err := s.client.Do(req, &deleteResponse)
	if err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}

type BlocklistIdentifiersResponse struct {
	Data       []*BlocklistIdentifierResponse `json:"data"`
	TotalCount int64                          `json:"total_count"`
}

func (s *BlocklistsService) ListAllIdentifiers() (*BlocklistIdentifiersResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, BlocklistsUrl)

	var blocklistIdentifiersResponse *BlocklistIdentifiersResponse
	_, err := s.client.Do(req, &blocklistIdentifiersResponse)
	if err != nil {
		return nil, err
	}
	return blocklistIdentifiersResponse, nil
}
