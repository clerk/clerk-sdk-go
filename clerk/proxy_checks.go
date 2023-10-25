package clerk

import "net/http"

type ProxyChecksService service

type ProxyCheck struct {
	Object     string `json:"object"`
	ID         string `json:"id"`
	DomainID   string `json:"domain_id"`
	ProxyURL   string `json:"proxy_url"`
	Successful bool   `json:"successful"`
	LastRunAt  *int64 `json:"last_run_at"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type CreateProxyCheckParams struct {
	DomainID string `json:"domain_id"`
	ProxyURL string `json:"proxy_url"`
}

func (s *ProxyChecksService) Create(params CreateProxyCheckParams) (*ProxyCheck, error) {
	req, _ := s.client.NewRequest(http.MethodPost, ProxyChecksURL, &params)
	var proxyCheck ProxyCheck
	_, err := s.client.Do(req, &proxyCheck)
	if err != nil {
		return nil, err
	}
	return &proxyCheck, nil
}
