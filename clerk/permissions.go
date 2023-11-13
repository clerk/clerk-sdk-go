package clerk

type Permission struct {
	Object      string `json:"object"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type PermissionsResponse struct {
	Data       []Permission `json:"data"`
	TotalCount int64        `json:"total_count"`
}
