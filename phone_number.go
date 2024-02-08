package clerk

type PhoneNumber struct {
	APIResource

	Object                  string `json:"object"`
	ID                      string `json:"id"`
	PhoneNumber             string `json:"phone_number"`
	ReservedForSecondFactor bool   `json:"reserved_for_second_factor"`
	DefaultSecondFactor     bool   `json:"default_second_factor"`
	Reserved                bool   `json:"reserved"`

	Verification *struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
		Attempts int    `json:"attempts"`
		ExpireAt int    `json:"expire_at"`
	} `json:"verification"`

	LinkedTo []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"linked_to"`

	BackupCodes []string `json:"backup_codes"`
}
