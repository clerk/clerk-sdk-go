package clerk

import "net/http"

type RedirectURLsService service

type RedirectURLResponse struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	URL       string `json:"url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CreateRedirectURLParams struct {
	URL string `json:"url"`
}

func (s *RedirectURLsService) Create(params CreateRedirectURLParams) (*RedirectURLResponse, error) {
	req, _ := s.client.NewRequest(http.MethodPost, RedirectURLsUrl, &params)

	var redirectURLResponse RedirectURLResponse
	_, err := s.client.Do(req, &redirectURLResponse)
	if err != nil {
		return nil, err
	}
	return &redirectURLResponse, nil
}

func (s *RedirectURLsService) ListAll() ([]*RedirectURLResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, RedirectURLsUrl, nil)

	var redirectURLResponses []*RedirectURLResponse
	_, err := s.client.Do(req, &redirectURLResponses)
	if err != nil {
		return nil, err
	}
	return redirectURLResponses, nil
}

func (s *RedirectURLsService) Delete(redirectURLID string) (*DeleteResponse, error) {
	req, _ := s.client.NewRequest(http.MethodDelete, RedirectURLsUrl+"/"+redirectURLID, nil)

	var deleteResponse DeleteResponse
	_, err := s.client.Do(req, &deleteResponse)
	if err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}
