package clerk

type Domain struct {
	APIResource
	ID                string        `json:"id"`
	Object            string        `json:"object"`
	Name              string        `json:"name"`
	IsSatellite       bool          `json:"is_satellite"`
	FrontendAPIURL    string        `json:"frontend_api_url"`
	AccountPortalURL  *string       `json:"accounts_portal_url,omitempty"`
	ProxyURL          *string       `json:"proxy_url,omitempty"`
	CNAMETargets      []CNAMETarget `json:"cname_targets,omitempty"`
	DevelopmentOrigin string        `json:"development_origin"`
}

type CNAMETarget struct {
	Host  string `json:"host"`
	Value string `json:"value"`
}

type DomainList struct {
	APIResource
	Domains    []*Domain `json:"data"`
	TotalCount int64     `json:"total_count"`
}
