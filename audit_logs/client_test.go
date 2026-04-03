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
				"id":             "019400f7-c6e4-7f00-8000-000000000001",
				"object":         "audit_log",
				"type":           "user.created",
				"source":         "bapi",
				"event_time":     1705315800000,
				"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"subject":        "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
				"client_id":      "client_123",
				"trace_id":       "00000000000000000000000000000001",
				"span_id":        "0000000000000001",
				"parent_span_id": nil,
				"payload":        map[string]interface{}{},
				"impersonator": map[string]interface{}{
					"user_id": "user_2xPNClBrCHGhpOITVJlhdhBfGS8",
				},
				"event_context": map[string]interface{}{
					"environment": map[string]interface{}{
						"type": "production",
						"application": map[string]interface{}{
							"id":   "app_123",
							"name": "My Application",
						},
						"domain": map[string]interface{}{
							"id":   "domain_123",
							"name": "example.com",
						},
						"primary_domain": map[string]interface{}{
							"id":   "domain_456",
							"name": "primary.example.com",
						},
					},
					"device": map[string]interface{}{
						"ip_address": "192.168.1.1",
						"user_agent": "Mozilla/5.0",
						"browser": map[string]interface{}{
							"name":    "Chrome",
							"version": "120.0.0",
						},
						"device_type":      "desktop",
						"is_mobile":        false,
						"clerk_js_version": "5.0.0",
						"is_native":        false,
						"location": map[string]interface{}{
							"city":    "New York",
							"country": "US",
						},
					},
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

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    responseJSON,
			Method: http.MethodGet,
			Path:   "/v1/audit_logs",
			Query: &url.Values{
				"limit":                []string{"10"},
				"subject":              []string{"user_2xPNCmYKPnPKaF3h0Ll8qas0I0w"},
				"client_id":            []string{"client_123"},
				"impersonator_user_id": []string{"user_2xPNClBrCHGhpOITVJlhdhBfGS8"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		Limit:              clerk.Int64(10),
		Subject:            clerk.String("user_2xPNCmYKPnPKaF3h0Ll8qas0I0w"),
		ClientID:           clerk.String("client_123"),
		ImpersonatorUserID: clerk.String("user_2xPNClBrCHGhpOITVJlhdhBfGS8"),
	}
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, 1, len(list.AuditLogs))

	auditLog := list.AuditLogs[0]
	require.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", auditLog.ID)
	require.Equal(t, "audit_log", auditLog.Object)
	require.Equal(t, "user.created", auditLog.Type)
	require.NotNil(t, auditLog.Source)
	require.Equal(t, "bapi", *auditLog.Source)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", auditLog.Actor)
	require.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", auditLog.Subject)
	require.NotNil(t, auditLog.ClientID)
	require.Equal(t, "client_123", *auditLog.ClientID)
	require.Equal(t, "00000000000000000000000000000001", auditLog.TraceID)
	require.Equal(t, "0000000000000001", auditLog.SpanID)
	require.Nil(t, auditLog.ParentSpanID)
	require.NotNil(t, auditLog.Payload)

	// Verify EventContext fields
	require.NotNil(t, auditLog.EventContext.Environment)
	require.NotNil(t, auditLog.EventContext.Environment.Type)
	require.Equal(t, "production", *auditLog.EventContext.Environment.Type)
	require.NotNil(t, auditLog.EventContext.Environment.Application)
	require.Equal(t, "app_123", *auditLog.EventContext.Environment.Application.ID)
	require.Equal(t, "My Application", *auditLog.EventContext.Environment.Application.Name)
	require.NotNil(t, auditLog.EventContext.Environment.Domain)
	require.Equal(t, "domain_123", *auditLog.EventContext.Environment.Domain.ID)
	require.Equal(t, "example.com", *auditLog.EventContext.Environment.Domain.Name)
	require.NotNil(t, auditLog.EventContext.Environment.PrimaryDomain)
	require.Equal(t, "domain_456", *auditLog.EventContext.Environment.PrimaryDomain.ID)
	require.Equal(t, "primary.example.com", *auditLog.EventContext.Environment.PrimaryDomain.Name)

	require.NotNil(t, auditLog.EventContext.DeviceInfo)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.IPAddress)
	require.Equal(t, "192.168.1.1", *auditLog.EventContext.DeviceInfo.IPAddress)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.UserAgent)
	require.Equal(t, "Mozilla/5.0", *auditLog.EventContext.DeviceInfo.UserAgent)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.Browser)
	require.Equal(t, "Chrome", *auditLog.EventContext.DeviceInfo.Browser.Name)
	require.Equal(t, "120.0.0", *auditLog.EventContext.DeviceInfo.Browser.Version)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.DeviceType)
	require.Equal(t, "desktop", *auditLog.EventContext.DeviceInfo.DeviceType)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.IsMobile)
	require.False(t, *auditLog.EventContext.DeviceInfo.IsMobile)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.ClerkJSVersion)
	require.Equal(t, "5.0.0", *auditLog.EventContext.DeviceInfo.ClerkJSVersion)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.IsNative)
	require.False(t, *auditLog.EventContext.DeviceInfo.IsNative)
	require.NotNil(t, auditLog.EventContext.DeviceInfo.Location)
	require.Equal(t, "New York", *auditLog.EventContext.DeviceInfo.Location.City)
	require.Equal(t, "US", *auditLog.EventContext.DeviceInfo.Location.Country)

	// Verify Impersonator field
	require.NotNil(t, auditLog.Impersonator)
	require.NotNil(t, auditLog.Impersonator.UserID)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *auditLog.Impersonator.UserID)

	require.NotNil(t, list.Cursor)
	require.Equal(t, expectedCursor, *list.Cursor.StartingAfter)
	require.Equal(t, expectedCursor, *list.Cursor.EndingBefore)
	require.True(t, list.Cursor.HasNextPage)
	require.Equal(t, clerk.NextPageTrue, list.Cursor.NextPageStatus)
}
