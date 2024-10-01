package clerk

import "encoding/json"

type SessionActivity struct {
	Object         string  `json:"object"`
	ID             string  `json:"id"`
	DeviceType     *string `json:"device_type,omitempty"`
	IsMobile       bool    `json:"is_mobile"`
	BrowserName    *string `json:"browser_name,omitempty"`
	BrowserVersion *string `json:"browser_version,omitempty"`
	IPAddress      *string `json:"ip_address,omitempty"`
	City           *string `json:"city,omitempty"`
	Country        *string `json:"country,omitempty"`
}

type Session struct {
	APIResource
	Object                   string           `json:"object"`
	ID                       string           `json:"id"`
	ClientID                 string           `json:"client_id"`
	UserID                   string           `json:"user_id"`
	Status                   string           `json:"status"`
	LastActiveOrganizationID string           `json:"last_active_organization_id,omitempty"`
	LatestActivity           *SessionActivity `json:"latest_activity,omitempty"`
	Actor                    json.RawMessage  `json:"actor,omitempty"`
	LastActiveAt             int64            `json:"last_active_at"`
	ExpireAt                 int64            `json:"expire_at"`
	AbandonAt                int64            `json:"abandon_at"`
	CreatedAt                int64            `json:"created_at"`
	UpdatedAt                int64            `json:"updated_at"`
}

type SessionList struct {
	APIResource
	Sessions   []*Session `json:"data"`
	TotalCount int64      `json:"total_count"`
}
