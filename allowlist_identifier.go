package clerk

type AllowlistIdentifier struct {
	APIResource
	Object         string  `json:"object"`
	ID             string  `json:"id"`
	Identifier     string  `json:"identifier"`
	IdentifierType string  `json:"identifier_type"`
	InvitationID   *string `json:"invitation_id,omitempty"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
}

type AllowlistIdentifierList struct {
	APIResource
	AllowlistIdentifiers []*AllowlistIdentifier `json:"data"`
	TotalCount           int64                  `json:"total_count"`
}
