package agenttask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestAgentTaskCreate(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
	"on_behalf_of":{"user_id":"user_123"},
	"permissions":"*",
	"agent_name":"test-agent",
	"task_description":"Check the billing page",
	"redirect_url":"https://example.com/billing",
	"session_max_duration_in_seconds":600,
	"session_restricted":true,
	"session_organization_id":"org_123"
}`),
				Out: json.RawMessage(`{
	"object":"agent_task",
	"agent_id":"agt_123",
	"agent_task_id":"agttsk_123",
	"task_id":"agttsk_123",
	"url":"https://example.clerk.accounts.dev/v1/agent-tasks/accept?ticket=tkt_123",
	"session_restricted":true,
	"session_organization_id":"org_123"
}`),
				Path:   "/v1/agents/tasks",
				Method: http.MethodPost,
			},
		},
	}))

	task, err := Create(context.Background(), &CreateParams{
		OnBehalfOf:                  &OnBehalfOf{UserID: clerk.String("user_123")},
		Permissions:                 clerk.String("*"),
		AgentName:                   clerk.String("test-agent"),
		TaskDescription:             clerk.String("Check the billing page"),
		RedirectURL:                 clerk.String("https://example.com/billing"),
		SessionMaxDurationInSeconds: clerk.Int64(600),
		SessionRestricted:           clerk.Bool(true),
		OrganizationID:              clerk.String("org_123"),
	})
	require.NoError(t, err)
	require.Equal(t, "agent_task", task.Object)
	require.Equal(t, "agt_123", task.AgentID)
	require.Equal(t, "agttsk_123", task.AgentTaskID)
	require.NotNil(t, task.URL)
	require.Equal(t, "https://example.clerk.accounts.dev/v1/agent-tasks/accept?ticket=tkt_123", *task.URL)
	require.True(t, task.SessionRestricted)
	require.NotNil(t, task.OrganizationID)
	require.Equal(t, "org_123", *task.OrganizationID)
}

func TestAgentTaskCreate_Identifier(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
	"on_behalf_of":{"identifier":"user@example.com"},
	"permissions":"*",
	"agent_name":"test-agent",
	"task_description":"Check the billing page",
	"redirect_url":"https://example.com/billing"
}`),
				Out: json.RawMessage(`{
	"object":"agent_task",
	"agent_id":"agt_123",
	"agent_task_id":"agttsk_123",
	"task_id":"agttsk_123",
	"url":"https://example.clerk.accounts.dev/v1/agent-tasks/accept?ticket=tkt_123",
	"session_restricted":false,
	"session_organization_id":null
}`),
				Path:   "/v1/agents/tasks",
				Method: http.MethodPost,
			},
		},
	}))

	task, err := Create(context.Background(), &CreateParams{
		OnBehalfOf:      &OnBehalfOf{Identifier: clerk.String("user@example.com")},
		Permissions:     clerk.String("*"),
		AgentName:       clerk.String("test-agent"),
		TaskDescription: clerk.String("Check the billing page"),
		RedirectURL:     clerk.String("https://example.com/billing"),
	})
	require.NoError(t, err)
	require.Equal(t, "agttsk_123", task.AgentTaskID)
	require.False(t, task.SessionRestricted)
	require.Nil(t, task.OrganizationID)
}

func TestAgentTaskCreate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Status: http.StatusUnprocessableEntity,
				Out: json.RawMessage(`{
	"errors":[{"code":"form_param_missing"}],
	"clerk_trace_id":"create-trace-id"
}`),
			},
		},
	}))

	_, err := Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "form_param_missing", apiErr.Errors[0].Code)
}

func TestAgentTaskRevoke(t *testing.T) {
	id := "agttsk_456"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				Out: json.RawMessage(fmt.Sprintf(`{
	"object":"agent_task",
	"agent_id":"agt_123",
	"agent_task_id":"%s",
	"task_id":"%s",
	"session_restricted":true,
	"session_organization_id":"org_123"
}`, id, id)),
				Path:   fmt.Sprintf("/v1/agents/tasks/%s/revoke", id),
				Method: http.MethodPost,
			},
		},
	}))

	task, err := Revoke(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, task.AgentTaskID)
	require.Nil(t, task.URL)
	require.True(t, task.SessionRestricted)
	require.NotNil(t, task.OrganizationID)
	require.Equal(t, "org_123", *task.OrganizationID)
}

func TestAgentTaskRevoke_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Status: http.StatusNotFound,
				Out: json.RawMessage(`{
	"errors":[{"code":"resource_not_found"}],
	"clerk_trace_id":"revoke-trace-id"
}`),
			},
		},
	}))

	_, err := Revoke(context.Background(), "agttsk_missing")
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "revoke-trace-id", apiErr.TraceID)
	require.Equal(t, "resource_not_found", apiErr.Errors[0].Code)
}
