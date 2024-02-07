package clerk

type BlocklistIdentifier struct {
	APIResource
	Object         string `json:"object"`
	ID             string `json:"id"`
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifier_type"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}

type BlocklistIdentifierList struct {
	APIResource
	BlocklistIdentifiers []*BlocklistIdentifier `json:"data"`
	TotalCount           int64                  `json:"total_count"`
}
