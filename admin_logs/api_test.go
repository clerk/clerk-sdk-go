package admin_logs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
)

func TestPackageList(t *testing.T) {
	// Cursor is base64 encoded: eventID|eventTimeMillis
	expectedCursor := "MDE5NDAwZjctYzZlNC03ZjAwLTgwMDAtMDAwMDAwMDAwMDAxfDE3MDUzMTU4MDAwMDA="

	response := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":             "019400f7-c6e4-7f00-8000-000000000001",
				"object":         "admin_log",
				"type":           "instance_key.created",
				"event_time":     1705315800000,
				"workspace":      "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK",
				"instance":       "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC",
				"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"subject":        "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1",
				"trace_id":       "00000000000000000000000000000001",
				"span_id":        "0000000000000001",
				"parent_span_id": nil,
				"session_id":     "0000000000000001",
				"impersonator": map[string]interface{}{
					"user_id": "user_2xPNClBrCHGhpOITVJlhdhBfGS8",
				},
			},
		},
		"cursor": map[string]interface{}{
			"starting_after":   expectedCursor,
			"ending_before":    expectedCursor,
			"has_next_page":    true,
			"next_page_status": "true",
		},
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodGet,
				Path:   "/v1/admin_logs",
				Query: &url.Values{
					"limit":    []string{"10"},
					"instance": []string{"ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC"},
					"type":     []string{"instance_key.created"},
				},
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &ListParams{
		Limit:    clerk.Int64(10),
		Instance: clerk.String("ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC"),
		Type:     clerk.String("instance_key.created"),
	}

	adminLogs, err := List(context.Background(), params)
	require.NoError(t, err)
	assert.Len(t, adminLogs.AdminLogs, 1)

	adminLog := adminLogs.AdminLogs[0]
	assert.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", adminLog.ID)
	assert.Equal(t, "admin_log", adminLog.Object)
	assert.Equal(t, "instance_key.created", adminLog.Type)
	assert.Equal(t, "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK", adminLog.Workspace)
	assert.Equal(t, "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC", adminLog.Instance)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", adminLog.Actor)
	assert.Equal(t, "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1", adminLog.Subject)
	assert.Equal(t, "00000000000000000000000000000001", adminLog.TraceID)
	assert.Equal(t, "0000000000000001", adminLog.SpanID)
	assert.Equal(t, "0000000000000001", *adminLog.SessionID)
	assert.Nil(t, adminLog.ParentSpanID)
	assert.NotNil(t, adminLog.Impersonator)
	assert.NotNil(t, adminLog.Impersonator.UserID)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *adminLog.Impersonator.UserID)

	assert.NotNil(t, adminLogs.Cursor)
	assert.Equal(t, expectedCursor, *adminLogs.Cursor.StartingAfter)
	assert.Equal(t, clerk.NextPageTrue, adminLogs.Cursor.NextPageStatus)
	assert.False(t, adminLogs.Cursor.RetentionLimitReached)
}

func TestPackageGet(t *testing.T) {
	response := map[string]interface{}{
		"id":             "019400f7-c6e4-7f00-8000-000000000001",
		"object":         "admin_log",
		"type":           "instance_key.created",
		"event_time":     1705315800000,
		"workspace":      "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK",
		"instance":       "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC",
		"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
		"subject":        "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1",
		"trace_id":       "00000000000000000000000000000001",
		"span_id":        "0000000000000001",
		"parent_span_id": nil,
		"session_id":     "0000000000000001",
		"payload":        map[string]interface{}{"key_name": "production_key"},
		"impersonator": map[string]interface{}{
			"user_id": "user_2xPNClBrCHGhpOITVJlhdhBfGS8",
		},
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodGet,
				Path:   "/v1/admin_logs/1705315800000:019400f7-c6e4-7f00-8000-000000000001",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	adminLog, err := Get(context.Background(), &GetParams{
		EventTimeMs: 1705315800000,
		EventID:     "019400f7-c6e4-7f00-8000-000000000001",
	})
	require.NoError(t, err)

	assert.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", adminLog.ID)
	assert.Equal(t, "admin_log", adminLog.Object)
	assert.Equal(t, "instance_key.created", adminLog.Type)
	assert.Equal(t, "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK", adminLog.Workspace)
	assert.Equal(t, "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC", adminLog.Instance)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", adminLog.Actor)
	assert.Equal(t, "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1", adminLog.Subject)
	assert.Equal(t, "00000000000000000000000000000001", adminLog.TraceID)
	assert.Equal(t, "0000000000000001", adminLog.SpanID)
	assert.Equal(t, "0000000000000001", *adminLog.SessionID)
	assert.Nil(t, adminLog.ParentSpanID)
	assert.NotNil(t, adminLog.Payload)
	assert.Equal(t, "production_key", adminLog.Payload["key_name"])
	assert.NotNil(t, adminLog.Impersonator)
	assert.NotNil(t, adminLog.Impersonator.UserID)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *adminLog.Impersonator.UserID)
}
