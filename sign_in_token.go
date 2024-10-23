package clerk

type SignInToken struct {
	APIResource
	Object    string  `json:"object"`
	ID        string  `json:"id"`
	Status    string  `json:"status"`
	UserID    string  `json:"user_id"`
	Token     string  `json:"token,omitempty"`
	URL       *string `json:"url,omitempty"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
