package clerk

type Role struct {
	Object      string       `json:"object"`
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Key         string       `json:"key"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   int64        `json:"created_at"`
	UpdatedAt   int64        `json:"updated_at"`
}

type RolesResponse struct {
	Data       []Role `json:"data"`
	TotalCount int64  `json:"total_count"`
}
