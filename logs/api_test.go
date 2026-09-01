package logs

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
				"object":         "log",
				"type":           "user.created",
				"event_time":     1705315800000,
				"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"actor_type":     "user",
				"subject":        "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
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
				Path:   "/v1/logs",
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

	logs, err := List(context.Background(), params)
	require.NoError(t, err)
	assert.Len(t, logs.Logs, 1)

	log := logs.Logs[0]
	assert.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", log.ID)
	assert.Equal(t, "log", log.Object)
	assert.Equal(t, "user.created", log.Type)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", log.Actor)
	assert.Equal(t, "user", log.ActorType)
	assert.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", log.Subject)
	assert.Equal(t, "00000000000000000000000000000001", log.TraceID)
	assert.Equal(t, "0000000000000001", log.SpanID)
	assert.Equal(t, "0000000000000001", *log.SessionID)
	assert.Nil(t, log.ParentSpanID)
	assert.NotNil(t, log.Impersonator)
	assert.NotNil(t, log.Impersonator.UserID)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *log.Impersonator.UserID)

	assert.NotNil(t, logs.Cursor)
	assert.Equal(t, expectedCursor, *logs.Cursor.StartingAfter)
	assert.Equal(t, clerk.NextPageTrue, logs.Cursor.NextPageStatus)
	assert.False(t, logs.Cursor.RetentionLimitReached)
}

func TestPackageGet(t *testing.T) {
	response := map[string]interface{}{
		"id":             "019400f7-c6e4-7f00-8000-000000000001",
		"object":         "log",
		"type":           "user.created",
		"event_time":     1705315800000,
		"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
		"actor_type":     "instance_key",
		"subject":        "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
		"trace_id":       "00000000000000000000000000000001",
		"span_id":        "0000000000000001",
		"parent_span_id": nil,
		"session_id":     "0000000000000001",
		"payload":        map[string]interface{}{"user_id": "user_123"},
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
				Path:   "/v1/logs/1705315800000:019400f7-c6e4-7f00-8000-000000000001",
			},
		},
	}

	clerk.SetBackend(clerk.NewBackend(backend))

	log, err := Get(context.Background(), &GetParams{
		EventTimeMs: 1705315800000,
		EventID:     "019400f7-c6e4-7f00-8000-000000000001",
	})
	require.NoError(t, err)

	assert.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", log.ID)
	assert.Equal(t, "log", log.Object)
	assert.Equal(t, "user.created", log.Type)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", log.Actor)
	assert.Equal(t, "instance_key", log.ActorType)
	assert.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", log.Subject)
	assert.Equal(t, "00000000000000000000000000000001", log.TraceID)
	assert.Equal(t, "0000000000000001", log.SpanID)
	assert.Equal(t, "0000000000000001", *log.SessionID)
	assert.Nil(t, log.ParentSpanID)
	assert.NotNil(t, log.Payload)
	assert.Equal(t, "user_123", log.Payload["user_id"])
	assert.NotNil(t, log.Impersonator)
	assert.NotNil(t, log.Impersonator.UserID)
	assert.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *log.Impersonator.UserID)
}
