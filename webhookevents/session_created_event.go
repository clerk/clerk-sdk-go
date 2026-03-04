package webhookevents

import "github.com/clerk/clerk-sdk-go/v3"

type SessionCreatedEvent struct {
	clerk.Session
	User clerk.User `json:"user"`
}
