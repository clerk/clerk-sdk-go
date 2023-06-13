package clerk

import (
	"fmt"
	"net/http"
)

type SAMLConnectionsService service

type SAMLConnection struct {
	ID                 string  `json:"id"`
	Object             string  `json:"object"`
	Name               string  `json:"name"`
	Domain             string  `json:"domain"`
	IdpEntityID        *string `json:"idp_entity_id"`
	IdpSsoURL          *string `json:"idp_sso_url"`
	IdpCertificate     *string `json:"idp_certificate"`
	AcsURL             string  `json:"acs_url"`
	SPEntityID         string  `json:"sp_entity_id"`
	Active             bool    `json:"active"`
	Provider           string  `json:"provider"`
	UserCount          int64   `json:"user_count"`
	SyncUserAttributes bool    `json:"sync_user_attributes"`
	CreatedAt          int64   `json:"created_at"`
	UpdatedAt          int64   `json:"updated_at"`
}

type ListSAMLConnectionsResponse struct {
	Data       []SAMLConnection `json:"data"`
	TotalCount int64            `json:"total_count"`
}

type ListSAMLConnectionsParams struct {
	Limit   *int
	Offset  *int
	Query   *string
	OrderBy *string
}

func (s SAMLConnectionsService) ListAll(params ListSAMLConnectionsParams) (*ListSAMLConnectionsResponse, error) {
	req, err := s.client.NewRequest(http.MethodGet, SAMLConnectionsUrl)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	addPaginationParams(query, PaginationParams{Limit: params.Limit, Offset: params.Offset})

	if params.Query != nil {
		query.Set("query", *params.Query)
	}
	if params.OrderBy != nil {
		query.Set("order_by", *params.OrderBy)
	}

	req.URL.RawQuery = query.Encode()

	samlConnections := &ListSAMLConnectionsResponse{}
	if _, err = s.client.Do(req, samlConnections); err != nil {
		return nil, err
	}

	return samlConnections, nil
}

func (s SAMLConnectionsService) Read(id string) (*SAMLConnection, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", SAMLConnectionsUrl, id))
	if err != nil {
		return nil, err
	}

	samlConnection := &SAMLConnection{}
	if _, err = s.client.Do(req, samlConnection); err != nil {
		return nil, err
	}

	return samlConnection, nil
}

type CreateSAMLConnectionParams struct {
	Name           string  `json:"name"`
	Domain         string  `json:"domain"`
	IdpEntityID    *string `json:"idp_entity_id,omitempty"`
	IdpSsoURL      *string `json:"idp_sso_url,omitempty"`
	IdpCertificate *string `json:"idp_certificate,omitempty"`
}

func (s SAMLConnectionsService) Create(params *CreateSAMLConnectionParams) (*SAMLConnection, error) {
	req, err := s.client.NewRequest(http.MethodPost, SAMLConnectionsUrl, params)
	if err != nil {
		return nil, err
	}

	resp := SAMLConnection{}
	if _, err = s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type UpdateSAMLConnectionParams struct {
	Name               *string `json:"name,omitempty"`
	Domain             *string `json:"domain,omitempty"`
	IdpEntityID        *string `json:"idp_entity_id,omitempty"`
	IdpSsoURL          *string `json:"idp_sso_url,omitempty"`
	IdpCertificate     *string `json:"idp_certificate,omitempty"`
	Active             *bool   `json:"active,omitempty"`
	SyncUserAttributes *bool   `json:"sync_user_attributes,omitempty"`
}

func (s SAMLConnectionsService) Update(id string, params *UpdateSAMLConnectionParams) (*SAMLConnection, error) {
	req, err := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%s", SAMLConnectionsUrl, id), params)
	if err != nil {
		return nil, err
	}

	resp := SAMLConnection{}
	if _, err = s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s SAMLConnectionsService) Delete(id string) (*DeleteResponse, error) {
	reqURL := fmt.Sprintf("%s/%s", SAMLConnectionsUrl, id)
	req, err := s.client.NewRequest(http.MethodDelete, reqURL)
	if err != nil {
		return nil, err
	}

	resp := &DeleteResponse{}
	if _, err = s.client.Do(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
