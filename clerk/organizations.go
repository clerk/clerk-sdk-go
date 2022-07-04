package clerk

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type OrganizationsService service

type Organization struct {
	Object          string          `json:"object"`
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Slug            *string         `json:"slug"`
	LogoURL         *string         `json:"logo_url"`
	PublicMetadata  json.RawMessage `json:"public_metadata"`
	PrivateMetadata json.RawMessage `json:"private_metadata,omitempty"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`
}

type OrganizationsResponse struct {
	Data       []Organization `json:"data"`
	TotalCount int64          `json:"total_count"`
}

type ListAllOrganizationsParams struct {
	Limit  *int
	Offset *int
}

func (s *OrganizationsService) ListAll(params ListAllOrganizationsParams) (*OrganizationsResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, OrganizationsUrl)

	query := req.URL.Query()
	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		query.Set("offset", strconv.Itoa(*params.Offset))
	}
	req.URL.RawQuery = query.Encode()

	var organizationsResponse *OrganizationsResponse
	_, err := s.client.Do(req, &organizationsResponse)
	if err != nil {
		return nil, err
	}
	return organizationsResponse, nil
}
