package clerk

import (
	"fmt"
	"net/http"
)

type DomainsService service

type DomainListResponse struct {
	Data       []Domain `json:"data"`
	TotalCount int      `json:"total_count"`
}

type DomainCNameTarget struct {
	Host  string `json:"host"`
	Value string `json:"value"`
}

type Domain struct {
	Object            string              `json:"object"`
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	IsSatellite       bool                `json:"is_satellite"`
	FapiURL           string              `json:"frontend_api_url"`
	AccountsPortalURL *string             `json:"accounts_portal_url,omitempty"`
	ProxyURL          *string             `json:"proxy_url,omitempty"`
	CNameTargets      []DomainCNameTarget `json:"cname_targets,omitempty"`
	DevelopmentOrigin string              `json:"development_origin"`
}

type CreateDomainParams struct {
	Name        string  `json:"name"`
	IsSatellite bool    `json:"is_satellite"`
	ProxyURL    *string `json:"proxy_url,omitempty"`
}

type UpdateDomainParams struct {
	Name     *string `json:"name,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
}

func (s *DomainsService) ListAll() (*DomainListResponse, error) {
	req, _ := s.client.NewRequest(http.MethodGet, DomainsURL)

	var response *DomainListResponse
	_, err := s.client.Do(req, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *DomainsService) Create(
	params CreateDomainParams,
) (*Domain, error) {
	req, _ := s.client.NewRequest(http.MethodPost, DomainsURL, &params)

	var domain Domain
	_, err := s.client.Do(req, &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (s *DomainsService) Update(
	domainID string,
	params UpdateDomainParams,
) (*Domain, error) {
	url := fmt.Sprintf("%s/%s", DomainsURL, domainID)
	req, _ := s.client.NewRequest(http.MethodPatch, url, &params)

	var domain Domain
	_, err := s.client.Do(req, &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (s *DomainsService) Delete(domainID string) (*DeleteResponse, error) {
	url := fmt.Sprintf("%s/%s", DomainsURL, domainID)
	req, _ := s.client.NewRequest(http.MethodDelete, url)

	var delResponse DeleteResponse
	if _, err := s.client.Do(req, &delResponse); err != nil {
		return nil, err
	}
	return &delResponse, nil
}
