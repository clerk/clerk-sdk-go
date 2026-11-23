package clerk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionAgentActor(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		actor string
		want  *AgentActor
	}{
		{
			name:  "agent task session",
			actor: `{"sub":"agt_123","type":"agent","task_id":"agttsk_123"}`,
			want:  &AgentActor{ID: "agt_123", TaskID: "agttsk_123"},
		},
		{
			name:  "restricted agent task session",
			actor: `{"sub":"agt_123","type":"agent","task_id":"agttsk_123","restricted":true}`,
			want:  &AgentActor{ID: "agt_123", TaskID: "agttsk_123", Restricted: true},
		},
		{
			name:  "agent session without a task",
			actor: `{"sub":"agt_123","type":"agent","extra":"value","restricted":true}`,
			want:  &AgentActor{ID: "agt_123", Restricted: true},
		},
		{
			name:  "poppy session",
			actor: `{"sub":"https://agent.example.com/client.json","type":"agent","restricted":true,"source":"poppy","source_id":"ps_123","external_user_id":"pu_123"}`,
			want:  &AgentActor{ID: "https://agent.example.com/client.json", Restricted: true},
		},
		{
			name:  "restricted flag that is not a boolean",
			actor: `{"sub":"agt_123","type":"agent","task_id":"agttsk_123","restricted":"true"}`,
			want:  &AgentActor{ID: "agt_123", TaskID: "agttsk_123"},
		},
		{
			name:  "agent actor with unreadable fields",
			actor: `{"sub":1,"type":"agent","task_id":["agttsk_123"],"restricted":true}`,
			want:  &AgentActor{Restricted: true},
		},
		{
			name:  "agent actor with only a type",
			actor: `{"type":"agent"}`,
			want:  &AgentActor{},
		},
		{
			name:  "impersonation session",
			actor: `{"sub":"user_456"}`,
		},
		{
			name:  "other actor type",
			actor: `{"sub":"user_456","type":"user"}`,
		},
		{
			name:  "type is not a string",
			actor: `{"sub":"agt_123","type":1}`,
		},
		{
			name:  "actor is not an object",
			actor: `"agent"`,
		},
		{
			name:  "actor null",
			actor: `null`,
		},
		{
			name: "actor absent",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			session := &Session{Actor: json.RawMessage(tc.actor)}
			require.Equal(t, tc.want, session.AgentActor())
		})
	}
}

func TestSessionAgentActorFromResponse(t *testing.T) {
	t.Parallel()
	var session Session
	payload := `{"object":"session","id":"sess_123","user_id":"user_123","status":"active","actor":{"sub":"agt_123","type":"agent","task_id":"agttsk_123","restricted":true}}`
	require.NoError(t, json.Unmarshal([]byte(payload), &session))
	require.Equal(t, &AgentActor{ID: "agt_123", TaskID: "agttsk_123", Restricted: true}, session.AgentActor())
	require.JSONEq(t, `{"sub":"agt_123","type":"agent","task_id":"agttsk_123","restricted":true}`, string(session.Actor))
}
