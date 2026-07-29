package clerk

// TrustedDevice represents a device enrolled by a user for trusted-device authentication.
type TrustedDevice struct {
	APIResource
	ID            string  `json:"id"`
	Object        string  `json:"object"`
	Platform      string  `json:"platform"`
	AppIdentifier string  `json:"app_identifier"`
	Name          *string `json:"name"`
	Algorithm     string  `json:"algorithm"`
	Status        string  `json:"status"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
	LastUsedAt    *int64  `json:"last_used_at"`
	RevokedAt     *int64  `json:"revoked_at"`
}
