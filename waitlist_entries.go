package clerk

type WaitlistEntry struct {
	APIResource
	Object       string      `json:"object"`
	ID           string      `json:"id"`
	EmailAddress string      `json:"email_address"`
	Status       string      `json:"status"`
	IsLocked     bool        `json:"is_locked"`
	CreatedAt    int64       `json:"created_at"`
	UpdatedAt    int64       `json:"updated_at"`
	Invitation   *Invitation `json:"invitation"`
}

type WaitlistEntries struct {
	APIResource
	WaitlistEntries []*WaitlistEntry `json:"data"`
}

type WaitlistEntriesList struct {
	APIResource
	WaitlistEntries []*WaitlistEntry `json:"data"`
	TotalCount      int64            `json:"total_count"`
}
