package clerk

type RedirectURL struct {
	APIResource
	Object    string `json:"object"`
	ID        string `json:"id"`
	URL       string `json:"url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type RedirectURLList struct {
	APIResource
	RedirectURLs []*RedirectURL `json:"data"`
	TotalCount   int64          `json:"total_count"`
}
