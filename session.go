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

// AgentActor returns the agent of an agent session, read from the session's
// actor, or nil for any other session. It is non-nil whenever the actor has
// type "agent", Poppy sessions included; fields that cannot be read from a
// malformed actor are left empty, so a check that refuses agent sessions fails
// closed.
//
// Its fields describe agent task sessions. For a Poppy session it returns the
// agent's client ID URL as ID and an empty TaskID, and it does not expose the
// Poppy identity; read that from the session token with
// SessionClaims.AgentClaim.
func (s *Session) AgentActor() *AgentActor {
	var actor struct {
		Type string `json:"type"`
		AgentActor
	}
	_ = json.Unmarshal(s.Actor, &actor)
	if actor.Type != "agent" {
		return nil
	}
	return &actor.AgentActor
}

// AgentActor describes the agent of an agent session, as the session object
// reports it. AgentClaim reads the same agent from the session token.
type AgentActor struct {
	// ID is the actor's subject. It is the same value as AgentClaim.ID.
	ID string `json:"sub"`
	// TaskID is the ID of the agent task that created the session, or empty
	// when no agent task created it.
	TaskID string `json:"task_id"`
	// Restricted reports whether Clerk's Frontend API refuses the session's
	// writes. It is false when the actor omits it. Clerk enforces the
	// restriction; to refuse agent sessions, check Session.AgentActor() != nil.
	Restricted bool `json:"restricted"`
}

type SessionList struct {
	APIResource
	Sessions   []*Session `json:"data"`
	TotalCount int64      `json:"total_count"`
}
