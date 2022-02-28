//go:build integration
// +build integration

package integration

import (
	"testing"
)

func TestSessions(t *testing.T) {
	client := createClient()

	sessions, err := client.Sessions().ListAll()
	if err != nil {
		t.Fatalf("Sessions.ListAll returned error: %v", err)
	}
	if sessions == nil {
		t.Fatalf("Sessions.ListAll returned nil")
	}

	for _, session := range sessions {
		sessionId := session.ID
		session, err := client.Sessions().Read(sessionId)
		if err != nil {
			t.Fatalf("Sessions.Read returned error: %v", err)
		}
		if session == nil {
			t.Fatalf("Sessions.Read returned nil")
		}
	}
}
