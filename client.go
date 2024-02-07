package clerk

type Client struct {
	APIResource
	Object              string     `json:"object"`
	ID                  string     `json:"id"`
	LastActiveSessionID *string    `json:"last_active_session_id"`
	SignInID            *string    `json:"sign_in_id"`
	SignUpID            *string    `json:"sign_up_id"`
	SessionIDs          []string   `json:"session_ids"`
	Sessions            []*Session `json:"sessions"`
	CreatedAt           int64      `json:"created_at"`
	UpdatedAt           int64      `json:"updated_at"`
}

type ClientList struct {
	APIResource
	Clients    []*Client `json:"data"`
	TotalCount int64     `json:"total_count"`
}
