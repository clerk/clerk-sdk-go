package clerk

// AgentTask is a task that an agent performs on behalf of a user. Visiting
// its URL creates a session for that user.
type AgentTask struct {
	APIResource
	Object string `json:"object"`
	// AgentID is the ID of the agent, the same for every task created with
	// the same agent name in an instance.
	AgentID     string `json:"agent_id"`
	AgentTaskID string `json:"agent_task_id"`
	// URL signs the user in when visited. Only the response to a create
	// request has it.
	URL *string `json:"url,omitempty"`
	// SessionRestricted reports whether the session that the task creates is
	// restricted.
	SessionRestricted bool `json:"session_restricted"`
	// OrganizationID is the ID of the organization the session starts in, or
	// nil for the user's personal account.
	OrganizationID *string `json:"session_organization_id"`
}
