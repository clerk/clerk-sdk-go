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
				"actor_type":     "user",
				"subject":        "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
				"client_id":      "client_123",
				"trace_id":       "00000000000000000000000000000001",
				"span_id":        "0000000000000001",
				"parent_span_id": nil,
				"session_id":     "0000000000000001",
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
	require.Equal(t, "user", auditLog.ActorType)
	require.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", auditLog.Subject)
	require.NotNil(t, auditLog.ClientID)
	require.Equal(t, "client_123", *auditLog.ClientID)
	require.Equal(t, "00000000000000000000000000000001", auditLog.TraceID)
	require.Equal(t, "0000000000000001", auditLog.SpanID)
	require.Equal(t, "0000000000000001", *auditLog.SessionID)
	require.Nil(t, auditLog.ParentSpanID)

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
	require.Equal(t, clerk.NextPageTrue, list.Cursor.NextPageStatus)
	require.False(t, list.Cursor.RetentionLimitReached)
}

// TestList_FilterMatchAny verifies that the filter_match query parameter is
// serialized correctly when set to the OR-mode value, and exercises a
// multi-axis filter combination similar to a real "show me everything for
// this user OR this trace" query.
func TestList_FilterMatchAny(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{},
		"cursor": map[string]interface{}{
			"starting_after":   nil,
			"ending_before":    nil,
			"has_next_page":    false,
			"next_page_status": "false",
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
				"subject":      []string{"user_alpha"},
				"actor":        []string{"actor_b"},
				"trace_id":     []string{"00000000000000000000000000000001"},
				"filter_match": []string{"any"},
			},
		},
	}
	client := NewClient(config)

	filterMatch := clerk.AuditLogFilterMatchAny
	params := &ListParams{
		Subject:     clerk.String("user_alpha"),
		Actor:       clerk.String("actor_b"),
		TraceID:     clerk.String("00000000000000000000000000000001"),
		FilterMatch: &filterMatch,
	}
	_, err := client.List(context.Background(), params)
	require.NoError(t, err)
}

// TestList_FilterMatchAll verifies the explicit AND-mode value also
// round-trips as a query parameter.
func TestList_FilterMatchAll(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{},
		"cursor": map[string]interface{}{
			"starting_after":   nil,
			"ending_before":    nil,
			"has_next_page":    false,
			"next_page_status": "false",
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
				"subject":      []string{"user_alpha"},
				"type":         []string{"user.created"},
				"filter_match": []string{"all"},
			},
		},
	}
	client := NewClient(config)

	filterMatch := clerk.AuditLogFilterMatchAll
	params := &ListParams{
		Subject:     clerk.String("user_alpha"),
		Type:        clerk.String("user.created"),
		FilterMatch: &filterMatch,
	}
	_, err := client.List(context.Background(), params)
	require.NoError(t, err)
}

// TestList_FilterMatchOmitted verifies that a nil FilterMatch is not
// emitted to the query string at all (preserving the API's default
// behavior).
func TestList_FilterMatchOmitted(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{},
		"cursor": map[string]interface{}{
			"starting_after":   nil,
			"ending_before":    nil,
			"has_next_page":    false,
			"next_page_status": "false",
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
				"subject": []string{"user_alpha"},
			},
		},
	}
	client := NewClient(config)

	params := &ListParams{
		Subject: clerk.String("user_alpha"),
	}
	_, err := client.List(context.Background(), params)
	require.NoError(t, err)
}

// TestListParams_ToQuery_FilterMatch unit-tests ListParams.ToQuery
// directly to lock the query parameter name and serialized value.
func TestListParams_ToQuery_FilterMatch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		filterMatch *clerk.AuditLogFilterMatch
		want        []string
	}{
		{name: "nil omits the param", filterMatch: nil, want: nil},
		{
			name:        "any serializes as \"any\"",
			filterMatch: filterMatchPtr(clerk.AuditLogFilterMatchAny),
			want:        []string{"any"},
		},
		{
			name:        "all serializes as \"all\"",
			filterMatch: filterMatchPtr(clerk.AuditLogFilterMatchAll),
			want:        []string{"all"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			params := &ListParams{FilterMatch: tc.filterMatch}
			got := params.ToQuery()["filter_match"]
			require.Equal(t, tc.want, got)
		})
	}
}

func filterMatchPtr(v clerk.AuditLogFilterMatch) *clerk.AuditLogFilterMatch {
	return &v
}

// TestListParams_ToQuery_IPAddress locks the wire name for the IP filter
// so the SDK stays in sync with the ip_address query parameter accepted by
// the audit_logs list endpoint.
func TestListParams_ToQuery_IPAddress(t *testing.T) {
	t.Parallel()

	params := &ListParams{IPAddress: clerk.String("203.0.113.10")}
	require.Equal(t, []string{"203.0.113.10"}, params.ToQuery()["ip_address"])

	require.Empty(t, (&ListParams{}).ToQuery()["ip_address"], "nil IPAddress should be omitted")
}

// TestList_ActorTypeOmitted verifies that a response without actor_type
// decodes to the zero value rather than failing, so the SDK keeps working
// against API responses predating the field.
func TestList_ActorTypeOmitted(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":         "019400f7-c6e4-7f00-8000-000000000001",
				"object":     "audit_log",
				"type":       "user.created",
				"event_time": 1705315800000,
				"actor":      "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"subject":    "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
			},
		},
		"cursor": map[string]interface{}{
			"starting_after":   nil,
			"ending_before":    nil,
			"has_next_page":    false,
			"next_page_status": "false",
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
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Equal(t, 1, len(list.AuditLogs))
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", list.AuditLogs[0].Actor)
	require.Empty(t, list.AuditLogs[0].ActorType)
}

func TestList_RetentionLimitReached(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"data": []map[string]interface{}{},
		"cursor": map[string]interface{}{
			"starting_after":          nil,
			"ending_before":           nil,
			"has_next_page":           false,
			"next_page_status":        "false",
			"retention_limit_reached": true,
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
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.NotNil(t, list.Cursor)
	require.Equal(t, clerk.NextPageFalse, list.Cursor.NextPageStatus)
	require.True(t, list.Cursor.RetentionLimitReached)
}

func TestGet(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"id":             "019400f7-c6e4-7f00-8000-000000000001",
		"object":         "audit_log",
		"type":           "user.created",
		"source":         "bapi",
		"event_time":     1705315800000,
		"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
		"actor_type":     "instance_key",
		"subject":        "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w",
		"client_id":      "client_123",
		"trace_id":       "00000000000000000000000000000001",
		"span_id":        "0000000000000001",
		"parent_span_id": nil,
		"session_id":     "0000000000000001",
		"payload":        map[string]interface{}{"user_id": "user_123"},
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
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    responseJSON,
			Method: http.MethodGet,
			Path:   "/v1/audit_logs/1705315800000:019400f7-c6e4-7f00-8000-000000000001",
		},
	}
	client := NewClient(config)
	auditLog, err := client.Get(context.Background(), &GetParams{
		EventTimeMs: 1705315800000,
		EventID:     "019400f7-c6e4-7f00-8000-000000000001",
	})
	require.NoError(t, err)

	require.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", auditLog.ID)
	require.Equal(t, "audit_log", auditLog.Object)
	require.Equal(t, "user.created", auditLog.Type)
	require.NotNil(t, auditLog.Source)
	require.Equal(t, "bapi", *auditLog.Source)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", auditLog.Actor)
	require.Equal(t, "instance_key", auditLog.ActorType)
	require.Equal(t, "user_2xPNCmYKPnPKaF3h0Ll8qas0I0w", auditLog.Subject)
	require.NotNil(t, auditLog.ClientID)
	require.Equal(t, "client_123", *auditLog.ClientID)
	require.Equal(t, "00000000000000000000000000000001", auditLog.TraceID)
	require.Equal(t, "0000000000000001", auditLog.SpanID)
	require.Equal(t, "0000000000000001", *auditLog.SessionID)
	require.Nil(t, auditLog.ParentSpanID)

	require.NotNil(t, auditLog.Payload)
	require.Equal(t, "user_123", auditLog.Payload["user_id"])

	require.NotNil(t, auditLog.EventContext.Environment)
	require.NotNil(t, auditLog.EventContext.Environment.Type)
	require.Equal(t, "production", *auditLog.EventContext.Environment.Type)

	require.NotNil(t, auditLog.Impersonator)
	require.NotNil(t, auditLog.Impersonator.UserID)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *auditLog.Impersonator.UserID)
}
