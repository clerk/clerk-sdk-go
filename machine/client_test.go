package machine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestMachineClientCreate(t *testing.T) {
	t.Parallel()
	id := "machine_123"
	name := "Test Machine"
	instanceID := "instance_456"
	createdAt := int64(1640995200)
	updatedAt := int64(1640995200)
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"machine","instance_id":"%s","created_at":%d,"updated_at":%d}`, id, name, instanceID, createdAt, updatedAt)),
			Method: http.MethodPost,
			Path:   "/v1/machines",
		},
	}
	client := NewClient(config)
	machine, err := client.Create(context.Background(), &CreateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, machine.ID)
	require.Equal(t, name, machine.Name)
	require.Equal(t, instanceID, machine.InstanceID)
	require.Equal(t, createdAt, machine.CreatedAt)
	require.Equal(t, updatedAt, machine.UpdatedAt)
}

func TestMachineClientCreate_Error(t *testing.T) {
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
	_, err := client.Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestMachineClientGet(t *testing.T) {
	t.Parallel()
	id := "machine_123"
	name := "Test Machine"
	instanceID := "instance_456"
	createdAt := int64(1640995200)
	updatedAt := int64(1640995200)
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"machine","instance_id":"%s","created_at":%d,"updated_at":%d}`, id, name, instanceID, createdAt, updatedAt)),
			Method: http.MethodGet,
			Path:   "/v1/machines/" + id,
		},
	}
	client := NewClient(config)
	machine, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, machine.ID)
	require.Equal(t, name, machine.Name)
	require.Equal(t, instanceID, machine.InstanceID)
	require.Equal(t, createdAt, machine.CreatedAt)
	require.Equal(t, updatedAt, machine.UpdatedAt)
}

func TestMachineClientUpdate(t *testing.T) {
	t.Parallel()
	id := "machine_123"
	name := "Updated Machine"
	instanceID := "instance_456"
	createdAt := int64(1640995200)
	updatedAt := int64(1640995260)
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"name":"%s"}`, name)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"%s","object":"machine","instance_id":"%s","created_at":%d,"updated_at":%d}`, id, name, instanceID, createdAt, updatedAt)),
			Method: http.MethodPatch,
			Path:   "/v1/machines/" + id,
		},
	}
	client := NewClient(config)
	machine, err := client.Update(context.Background(), id, &UpdateParams{
		Name: clerk.String(name),
	})
	require.NoError(t, err)
	require.Equal(t, id, machine.ID)
	require.Equal(t, name, machine.Name)
	require.Equal(t, instanceID, machine.InstanceID)
	require.Equal(t, createdAt, machine.CreatedAt)
	require.Equal(t, updatedAt, machine.UpdatedAt)
}

func TestMachineClientUpdate_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"update-error-code"
	}],
	"clerk_trace_id":"update-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.Update(context.Background(), "machine_123", &UpdateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "update-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "update-error-code", apiErr.Errors[0].Code)
}

func TestMachineClientDelete(t *testing.T) {
	t.Parallel()
	id := "machine_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/machines/" + id,
		},
	}
	client := NewClient(config)
	deletedResource, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, deletedResource.ID)
	require.True(t, deletedResource.Deleted)
}

func TestMachineClientList(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
"data": [{"id":"machine_123","name":"Test Machine","object":"machine","instance_id":"instance_456","created_at":1640995200,"updated_at":1640995200}],
"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/machines",
			Query: &url.Values{
				"limit":  []string{"1"},
				"offset": []string{"2"},
				"query":  []string{"Test"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		Query: clerk.String("Test"),
	}
	params.Limit = clerk.Int64(1)
	params.Offset = clerk.Int64(2)
	machineList, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), machineList.TotalCount)
	require.Equal(t, 1, len(machineList.Machines))
	require.Equal(t, "machine_123", machineList.Machines[0].ID)
	require.Equal(t, "Test Machine", machineList.Machines[0].Name)
	require.Equal(t, "instance_456", machineList.Machines[0].InstanceID)
	require.Equal(t, int64(1640995200), machineList.Machines[0].CreatedAt)
	require.Equal(t, int64(1640995200), machineList.Machines[0].UpdatedAt)
}
