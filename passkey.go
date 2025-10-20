package clerk

type Passkey struct {
	ID           string        `json:"id"`
	Object       string        `json:"object"`
	Name         string        `json:"name"`
	LastUsedAt   *int64        `json:"last_used_at,omitempty"`
	Verification *Verification `json:"verification"`
	CreatedAt    int64         `json:"created_at"`
	UpdatedAt    int64         `json:"updated_at"`
}
