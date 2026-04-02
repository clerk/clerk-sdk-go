package audit_logs

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
				"object":         "audit_log",
				"type":           "user.created",
				"event_time":     1705315800000,
				"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"subject":        "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
				"trace_id":       "00000000000000000000000000000001",
				"span_id":        "0000000000000001",
				"parent_span_id": nil,
				"payload":        map[string]interface{}{},
				"impersonator": map[string]interface{}{
					"user_id": "user_2xPNClBrCHGhpOITVJlhdhBfGS8",
				},
			},
		},
		"cursor": map[string]interface{}{
			"starting_after": expectedCursor,
			"ending_before":  expectedCursor,
			"has_next_page":  "true",
		},
	}

	responseJSON, _ := json.Marshal(response)

	backend := &clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(responseJSON),
				Method: http.MethodGet,
				Path:   "/v1/audit_logs",
				Query: &url.Values{
					"limit":   []string{"10"},
					"subject": []string{"user_2xPNCmYKPnPKaF3h0Ll8qas0I0w"},
					"type":    []string{"user.created"},
				},
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	params := &ListParams{
		Limit:   clerk.Int64(10),
		Subject: clerk.String("user_2xPNCmYKPnPKaF3h0Ll8qas0I0w"),
		Type:    clerk.String("user.created"),
	}

	auditLogs, err := List(context.Background(), params)
	require.NoError(t, err)
	assert.Len(t, auditLogs.AuditLogs, 1)

	auditLog := auditLogs.AuditLogs[0]
	assert.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", auditLog.ID)
	assert.Equal(t, "audit_log", auditLog.Object)
	assert.Equal(t, "user.created", auditLog.Type)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", auditLog.Actor)
	assert.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", auditLog.Subject)
	assert.Equal(t, "00000000000000000000000000000001", auditLog.TraceID)
	assert.Equal(t, "0000000000000001", auditLog.SpanID)
	assert.Nil(t, auditLog.ParentSpanID)
	assert.NotNil(t, auditLog.Payload)
	assert.NotNil(t, auditLog.Impersonator)
	assert.NotNil(t, auditLog.Impersonator.UserID)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *auditLog.Impersonator.UserID)

	assert.NotNil(t, auditLogs.Cursor)
	assert.Equal(t, expectedCursor, *auditLogs.Cursor.StartingAfter)
	assert.Equal(t, clerk.NextPageTrue, auditLogs.Cursor.HasNextPage)
}
