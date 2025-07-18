package machinescope

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

func TestMachineScopeClientCreateScope(t *testing.T) {
	t.Parallel()
	machineID := "machine_123"
	otherMachineID := "machine_456"
	fromMachineID := "machine_123"
	toMachineID := "machine_456"
	createdAt := int64(1640995200)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"to_machine_id":"%s"}`, otherMachineID)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"machine_scope","from_machine_id":"%s","to_machine_id":"%s","created_at":%d}`, fromMachineID, toMachineID, createdAt)),
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/v1/machines/%s/scopes", machineID),
		},
	}
	client := NewClient(config)
	machineScope, err := client.CreateScope(context.Background(), machineID, &CreateScopeParams{
		ToMachineID: otherMachineID,
	})
	require.NoError(t, err)
	require.Equal(t, "machine_scope", machineScope.Object)
	require.Equal(t, fromMachineID, machineScope.FromMachineID)
	require.Equal(t, toMachineID, machineScope.ToMachineID)
	require.Equal(t, createdAt, machineScope.CreatedAt)
}

func TestMachineScopeClientCreateScope_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"create-error-code"
	}],
	"clerk_trace_id":"create-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.CreateScope(context.Background(), "machine_123", &CreateScopeParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestMachineScopeClientDeleteScope(t *testing.T) {
	t.Parallel()
	machineID := "machine_123"
	otherMachineID := "machine_456"
	fromMachineID := "machine_123"
	toMachineID := "machine_456"

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"machine_scope","from_machine_id":"%s","to_machine_id":"%s","deleted":true}`, fromMachineID, toMachineID)),
			Method: http.MethodDelete,
			Path:   fmt.Sprintf("/v1/machines/%s/scopes/%s", machineID, otherMachineID),
		},
	}
	client := NewClient(config)
	deletedMachineScope, err := client.DeleteScope(context.Background(), machineID, otherMachineID)
	require.NoError(t, err)
	require.Equal(t, "machine_scope", deletedMachineScope.Object)
	require.Equal(t, fromMachineID, deletedMachineScope.FromMachineID)
	require.Equal(t, toMachineID, deletedMachineScope.ToMachineID)
	require.True(t, deletedMachineScope.Deleted)
}

func TestMachineScopeClientDeleteScope_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"delete-error-code"
	}],
	"clerk_trace_id":"delete-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.DeleteScope(context.Background(), "machine_123", "machine_456")
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "delete-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "delete-error-code", apiErr.Errors[0].Code)
}
