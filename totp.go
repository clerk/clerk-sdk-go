package clerk

// Describes a Time-based One-time Password (TOTP) for a user.
type TOTP struct {
	APIResource
	Object      string   `json:"object"`
	ID          string   `json:"id"`
	Secret      *string  `json:"secret"`
	URI         *string  `json:"uri" `
	Verified    bool     `json:"verified"`
	BackupCodes []string `json:"backup_codes"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}
