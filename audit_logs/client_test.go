package audit_logs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
)

func TestList(t *testing.T) {
	t.Parallel()

	// Cursor is base64 encoded: eventID|eventTimeMillis
	// e.g., "019400f7-c6e4-7f00-8000-000000000001|1705315800000" -> "MDE5NDAwZjctYzZlNC03ZjAwLTgwMDAtMDAwMDAwMDAwMDAxfDE3MDUzMTU4MDAwMDA="
	expectedCursor := "MDE5NDAwZjctYzZlNC03ZjAwLTgwMDAtMDAwMDAwMDAwMDAxfDE3MDUzMTU4MDAwMDA="

	response := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":               "019400f7-c6e4-7f00-8000-000000000001",
				"object":           "audit_log",
				"type":             "user.created",
				"event_time":       1705315800000,
				"subject_instance": "ins_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"actor":            "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"subject":          "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
				"trace_id":         "00000000000000000000000000000001",
				"span_id":          "0000000000000001",
				"parent_span_id":   nil,
				"payload":          map[string]interface{}{},
			},
		},
		"cursor": map[string]interface{}{
			"starting_after": expectedCursor,
			"ending_before":  expectedCursor,
			"has_next_page":  true,
		},
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    responseJSON,
			Method: http.MethodGet,
			Path:   "/v1/audit_logs",
			Query: &url.Values{
				"limit":            []string{"10"},
				"subject_instance": []string{"ins_2xPNClBrCHGhpOITVJlhdhBfGS7"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		Limit:           clerk.Int64(10),
		SubjectInstance: clerk.String("ins_2xPNClBrCHGhpOITVJlhdhBfGS7"),
	}
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, 1, len(list.AuditLogs))

	auditLog := list.AuditLogs[0]
	require.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", auditLog.ID)
	require.Equal(t, "audit_log", auditLog.Object)
	require.Equal(t, "user.created", auditLog.Type)
	require.Equal(t, "ins_2xPNClBrCHGhpOITVJlhdhBfGS7", auditLog.SubjectInstance)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", auditLog.Actor)
	require.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", auditLog.Subject)
	require.Equal(t, "00000000000000000000000000000001", auditLog.TraceID)
	require.Equal(t, "0000000000000001", auditLog.SpanID)
	require.Nil(t, auditLog.ParentSpanID)
	require.NotNil(t, auditLog.Payload)

	require.NotNil(t, list.Cursor)
	require.Equal(t, expectedCursor, *list.Cursor.StartingAfter)
	require.Equal(t, expectedCursor, *list.Cursor.EndingBefore)
	require.True(t, list.Cursor.HasNextPage)
}
