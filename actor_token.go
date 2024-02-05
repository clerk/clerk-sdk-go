package clerk

import "encoding/json"

type ActorToken struct {
	APIResource
	Object    string          `json:"object"`
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Actor     json.RawMessage `json:"actor"`
	Token     string          `json:"token,omitempty"`
	URL       *string         `json:"url,omitempty"`
	Status    string          `json:"status"`
	CreatedAt int64           `json:"created_at"`
	UpdatedAt int64           `json:"updated_at"`
}
