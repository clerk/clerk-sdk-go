package clerk

type ProxyCheck struct {
	APIResource
	Object     string `json:"object"`
	ID         string `json:"id"`
	DomainID   string `json:"domain_id"`
	ProxyURL   string `json:"proxy_url"`
	Successful bool   `json:"successful"`
	LastRunAt  *int64 `json:"last_run_at"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}
