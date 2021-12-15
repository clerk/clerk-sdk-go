package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JWTTemplatesService service

type JWTTemplate struct {
	ID               string          `json:"id"`
	Object           string          `json:"object"`
	Name             string          `json:"name"`
	Claims           json.RawMessage `json:"claims"`
	Lifetime         int             `json:"lifetime"`
	AllowedClockSkew int             `json:"allowed_clock_skew"`
	CreatedAt        int64           `json:"created_at"`
	UpdatedAt        int64           `json:"updated_at"`
}

type CreateUpdateJWTTemplate struct {
	Claims           map[string]interface{}
	Name             string
	Lifetime         *int
	AllowedClockSkew *int
}

func (t CreateUpdateJWTTemplate) toRequest() (*createUpdateJWTTemplateRequest, error) {
	claimsBytes, err := json.Marshal(t.Claims)
	if err != nil {
		return nil, err
	}

	return &createUpdateJWTTemplateRequest{
		Claims:           string(claimsBytes),
		Name:             t.Name,
		Lifetime:         t.Lifetime,
		AllowedClockSkew: t.AllowedClockSkew,
	}, nil
}

type createUpdateJWTTemplateRequest struct {
	Claims           string `json:"claims"`
	Name             string `json:"name"`
	Lifetime         *int   `json:"lifetime,omitempty"`
	AllowedClockSkew *int   `json:"allowed_clock_skew,omitempty"`
}

func (s JWTTemplatesService) ListAll() ([]JWTTemplate, error) {
	req, err := s.client.NewRequest(http.MethodGet, JWTTemplatesUrl)
	if err != nil {
		return nil, err
	}

	jwtTemplates := make([]JWTTemplate, 0)
	if _, err = s.client.Do(req, &jwtTemplates); err != nil {
		return nil, err
	}

	return jwtTemplates, nil
}

func (s JWTTemplatesService) Read(id string) (*JWTTemplate, error) {
	reqURL := fmt.Sprintf("%s/%s", JWTTemplatesUrl, id)
	req, err := s.client.NewRequest(http.MethodGet, reqURL)
	if err != nil {
		return nil, err
	}

	jwtTemplate := &JWTTemplate{}
	if _, err = s.client.Do(req, jwtTemplate); err != nil {
		return nil, err
	}

	return jwtTemplate, nil
}

func (s JWTTemplatesService) Create(params *CreateUpdateJWTTemplate) (*JWTTemplate, error) {
	reqBody, err := params.toRequest()
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, JWTTemplatesUrl, reqBody)
	if err != nil {
		return nil, err
	}

	resp := JWTTemplate{}
	if _, err = s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s JWTTemplatesService) Update(id string, params *CreateUpdateJWTTemplate) (*JWTTemplate, error) {
	reqURL := fmt.Sprintf("%s/%s", JWTTemplatesUrl, id)

	reqBody, err := params.toRequest()
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPatch, reqURL, reqBody)
	if err != nil {
		return nil, err
	}

	resp := JWTTemplate{}
	if _, err = s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s JWTTemplatesService) Delete(id string) (*DeleteResponse, error) {
	reqURL := fmt.Sprintf("%s/%s", JWTTemplatesUrl, id)
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
