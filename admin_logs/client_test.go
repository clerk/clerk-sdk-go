package admin_logs

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
	expectedCursor := "MDE5NDAwZjctYzZlNC03ZjAwLTgwMDAtMDAwMDAwMDAwMDAxfDE3MDUzMTU4MDAwMDA="

	response := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":             "019400f7-c6e4-7f00-8000-000000000001",
				"object":         "admin_log",
				"type":           "instance_key.created",
				"source":         "bapi",
				"event_time":     1705315800000,
				"workspace":      "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK",
				"instance":       "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC",
				"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
				"subject":        "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1",
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
			Path:   "/v1/admin_logs",
			Query: &url.Values{
				"limit":                []string{"10"},
				"instance":             []string{"ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC"},
				"subject":              []string{"ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1"},
				"client_id":            []string{"client_123"},
				"impersonator_user_id": []string{"user_2xPNClBrCHGhpOITVJlhdhBfGS8"},
			},
		},
	}
	client := NewClient(config)
	params := &ListParams{
		Limit:              clerk.Int64(10),
		Instance:           clerk.String("ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC"),
		Subject:            clerk.String("ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1"),
		ClientID:           clerk.String("client_123"),
		ImpersonatorUserID: clerk.String("user_2xPNClBrCHGhpOITVJlhdhBfGS8"),
	}
	list, err := client.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, 1, len(list.AdminLogs))

	adminLog := list.AdminLogs[0]
	require.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", adminLog.ID)
	require.Equal(t, "admin_log", adminLog.Object)
	require.Equal(t, "instance_key.created", adminLog.Type)
	require.NotNil(t, adminLog.Source)
	require.Equal(t, "bapi", *adminLog.Source)
	// Workspace and Instance are unique to admin logs (audit logs has only
	// instance scoping baked into the URL); make sure both round-trip.
	require.Equal(t, "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK", adminLog.Workspace)
	require.Equal(t, "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC", adminLog.Instance)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", adminLog.Actor)
	require.Equal(t, "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1", adminLog.Subject)
	require.NotNil(t, adminLog.ClientID)
	require.Equal(t, "client_123", *adminLog.ClientID)
	require.Equal(t, "00000000000000000000000000000001", adminLog.TraceID)
	require.Equal(t, "0000000000000001", adminLog.SpanID)
	require.Equal(t, "0000000000000001", *adminLog.SessionID)
	require.Nil(t, adminLog.ParentSpanID)

	// EventContext reuses AuditLogContext on the SDK side; confirm the nested
	// fields decode against that type.
	require.NotNil(t, adminLog.EventContext.Environment)
	require.NotNil(t, adminLog.EventContext.Environment.Type)
	require.Equal(t, "production", *adminLog.EventContext.Environment.Type)
	require.NotNil(t, adminLog.EventContext.Environment.Application)
	require.Equal(t, "app_123", *adminLog.EventContext.Environment.Application.ID)
	require.Equal(t, "My Application", *adminLog.EventContext.Environment.Application.Name)
	require.NotNil(t, adminLog.EventContext.Environment.Domain)
	require.Equal(t, "domain_123", *adminLog.EventContext.Environment.Domain.ID)
	require.Equal(t, "example.com", *adminLog.EventContext.Environment.Domain.Name)
	require.NotNil(t, adminLog.EventContext.Environment.PrimaryDomain)
	require.Equal(t, "domain_456", *adminLog.EventContext.Environment.PrimaryDomain.ID)
	require.Equal(t, "primary.example.com", *adminLog.EventContext.Environment.PrimaryDomain.Name)

	require.NotNil(t, adminLog.EventContext.DeviceInfo)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.IPAddress)
	require.Equal(t, "192.168.1.1", *adminLog.EventContext.DeviceInfo.IPAddress)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.UserAgent)
	require.Equal(t, "Mozilla/5.0", *adminLog.EventContext.DeviceInfo.UserAgent)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.Browser)
	require.Equal(t, "Chrome", *adminLog.EventContext.DeviceInfo.Browser.Name)
	require.Equal(t, "120.0.0", *adminLog.EventContext.DeviceInfo.Browser.Version)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.DeviceType)
	require.Equal(t, "desktop", *adminLog.EventContext.DeviceInfo.DeviceType)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.IsMobile)
	require.False(t, *adminLog.EventContext.DeviceInfo.IsMobile)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.ClerkJSVersion)
	require.Equal(t, "5.0.0", *adminLog.EventContext.DeviceInfo.ClerkJSVersion)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.IsNative)
	require.False(t, *adminLog.EventContext.DeviceInfo.IsNative)
	require.NotNil(t, adminLog.EventContext.DeviceInfo.Location)
	require.Equal(t, "New York", *adminLog.EventContext.DeviceInfo.Location.City)
	require.Equal(t, "US", *adminLog.EventContext.DeviceInfo.Location.Country)

	require.NotNil(t, adminLog.Impersonator)
	require.NotNil(t, adminLog.Impersonator.UserID)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *adminLog.Impersonator.UserID)

	require.NotNil(t, list.Cursor)
	require.Equal(t, expectedCursor, *list.Cursor.StartingAfter)
	require.Equal(t, expectedCursor, *list.Cursor.EndingBefore)
	require.Equal(t, clerk.NextPageTrue, list.Cursor.NextPageStatus)
	require.False(t, list.Cursor.RetentionLimitReached)
}

// TestList_TypeWildcard verifies the trailing-* type filter round-trips as
// a query parameter. The BAPI handler treats this as a LIKE prefix match.
func TestList_TypeWildcard(t *testing.T) {
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
			Path:   "/v1/admin_logs",
			Query: &url.Values{
				"type": []string{"instance_key.*"},
			},
		},
	}
	client := NewClient(config)
	_, err := client.List(context.Background(), &ListParams{
		Type: clerk.String("instance_key.*"),
	})
	require.NoError(t, err)
}

// TestList_FilterMatchAny verifies the OR-mode filter_match serializes to
// the expected wire value and exercises a multi-axis filter combination.
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
			Path:   "/v1/admin_logs",
			Query: &url.Values{
				"instance":     []string{"ins_alpha"},
				"actor":        []string{"actor_b"},
				"trace_id":     []string{"00000000000000000000000000000001"},
				"filter_match": []string{"any"},
			},
		},
	}
	client := NewClient(config)

	filterMatch := clerk.AdminLogFilterMatchAny
	params := &ListParams{
		Instance:    clerk.String("ins_alpha"),
		Actor:       clerk.String("actor_b"),
		TraceID:     clerk.String("00000000000000000000000000000001"),
		FilterMatch: &filterMatch,
	}
	_, err := client.List(context.Background(), params)
	require.NoError(t, err)
}

// TestList_FilterMatchAll verifies the explicit AND-mode value round-trips.
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
			Path:   "/v1/admin_logs",
			Query: &url.Values{
				"subject":      []string{"ikey_alpha"},
				"type":         []string{"instance_key.created"},
				"filter_match": []string{"all"},
			},
		},
	}
	client := NewClient(config)

	filterMatch := clerk.AdminLogFilterMatchAll
	params := &ListParams{
		Subject:     clerk.String("ikey_alpha"),
		Type:        clerk.String("instance_key.created"),
		FilterMatch: &filterMatch,
	}
	_, err := client.List(context.Background(), params)
	require.NoError(t, err)
}

// TestList_FilterMatchOmitted verifies a nil FilterMatch is not emitted to
// the query string, preserving the API's default AND behavior.
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
			Path:   "/v1/admin_logs",
			Query: &url.Values{
				"subject": []string{"ikey_alpha"},
			},
		},
	}
	client := NewClient(config)
	_, err := client.List(context.Background(), &ListParams{
		Subject: clerk.String("ikey_alpha"),
	})
	require.NoError(t, err)
}

// TestListParams_ToQuery_FilterMatch unit-tests ListParams.ToQuery directly
// to lock the query parameter name and serialized value.
func TestListParams_ToQuery_FilterMatch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		filterMatch *clerk.AdminLogFilterMatch
		want        []string
	}{
		{name: "nil omits the param", filterMatch: nil, want: nil},
		{
			name:        "any serializes as \"any\"",
			filterMatch: filterMatchPtr(clerk.AdminLogFilterMatchAny),
			want:        []string{"any"},
		},
		{
			name:        "all serializes as \"all\"",
			filterMatch: filterMatchPtr(clerk.AdminLogFilterMatchAll),
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

func filterMatchPtr(v clerk.AdminLogFilterMatch) *clerk.AdminLogFilterMatch {
	return &v
}

// TestListParams_ToQuery_Instance locks the wire name of the instance
// filter — the parameter unique to admin logs.
func TestListParams_ToQuery_Instance(t *testing.T) {
	t.Parallel()

	params := &ListParams{Instance: clerk.String("ins_alpha")}
	require.Equal(t, []string{"ins_alpha"}, params.ToQuery()["instance"])

	// Confirm nil pointer leaves the param off the wire entirely.
	emptyParams := &ListParams{}
	require.Nil(t, emptyParams.ToQuery()["instance"])
}

// TestList_TimeBoundsAndEndUserFacing verifies that the time-window and
// end_user_facing_only params serialize correctly and that they are
// independent of the filter_match group (always ANDed).
func TestList_TimeBoundsAndEndUserFacing(t *testing.T) {
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
			Path:   "/v1/admin_logs",
			Query: &url.Values{
				"event_time_after":     []string{"1705315800000"},
				"event_time_before":    []string{"1705325800000"},
				"end_user_facing_only": []string{"true"},
			},
		},
	}
	client := NewClient(config)
	_, err := client.List(context.Background(), &ListParams{
		EventTimeAfter:    clerk.Int64(1705315800000),
		EventTimeBefore:   clerk.Int64(1705325800000),
		EndUserFacingOnly: clerk.Bool(true),
	})
	require.NoError(t, err)
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
			Path:   "/v1/admin_logs",
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
		"object":         "admin_log",
		"type":           "instance_key.created",
		"source":         "bapi",
		"event_time":     1705315800000,
		"workspace":      "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK",
		"instance":       "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC",
		"actor":          "user_2xPNClBrCHGhpOITVJlhdhBfGS7",
		"subject":        "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1",
		"client_id":      "client_123",
		"trace_id":       "00000000000000000000000000000001",
		"span_id":        "0000000000000001",
		"parent_span_id": nil,
		"session_id":     "0000000000000001",
		"payload":        map[string]interface{}{"key_name": "production_key"},
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
			Path:   "/v1/admin_logs/1705315800000:019400f7-c6e4-7f00-8000-000000000001",
		},
	}
	client := NewClient(config)
	adminLog, err := client.Get(context.Background(), &GetParams{
		EventTimeMs: 1705315800000,
		EventID:     "019400f7-c6e4-7f00-8000-000000000001",
	})
	require.NoError(t, err)

	require.Equal(t, "019400f7-c6e4-7f00-8000-000000000001", adminLog.ID)
	require.Equal(t, "admin_log", adminLog.Object)
	require.Equal(t, "instance_key.created", adminLog.Type)
	require.NotNil(t, adminLog.Source)
	require.Equal(t, "bapi", *adminLog.Source)
	require.Equal(t, "org_2xPNCmYKPnPKaF3h0Ll8qas0WORK", adminLog.Workspace)
	require.Equal(t, "ins_2xPNCmYKPnPKaF3h0Ll8qas0INSTNC", adminLog.Instance)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS7", adminLog.Actor)
	require.Equal(t, "ikey_2xPNCmYKPnPKaF3h0Ll8qas0KEY1", adminLog.Subject)
	require.NotNil(t, adminLog.ClientID)
	require.Equal(t, "client_123", *adminLog.ClientID)
	require.Equal(t, "00000000000000000000000000000001", adminLog.TraceID)
	require.Equal(t, "0000000000000001", adminLog.SpanID)
	require.Equal(t, "0000000000000001", *adminLog.SessionID)
	require.Nil(t, adminLog.ParentSpanID)

	// AdminLogWithPayload is the only place the payload field appears.
	require.NotNil(t, adminLog.Payload)
	require.Equal(t, "production_key", adminLog.Payload["key_name"])

	require.NotNil(t, adminLog.EventContext.Environment)
	require.NotNil(t, adminLog.EventContext.Environment.Type)
	require.Equal(t, "production", *adminLog.EventContext.Environment.Type)

	require.NotNil(t, adminLog.Impersonator)
	require.NotNil(t, adminLog.Impersonator.UserID)
	require.Equal(t, "user_2xPNClBrCHGhpOITVJlhdhBfGS8", *adminLog.Impersonator.UserID)
}
