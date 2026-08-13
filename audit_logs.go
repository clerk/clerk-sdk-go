package clerk

type AuditLog struct {
	APIResource
	Object       string                `json:"object"`
	ID           string                `json:"id"`
	EventTime    int64                 `json:"event_time"`
	Actor        string                `json:"actor"`
	ActorType    string                `json:"actor_type"`
	Impersonator *ImpersonatorResponse `json:"impersonator"`
	Subject      string                `json:"subject"`
	Type         string                `json:"type"`
	Source       *string               `json:"source"`
	ClientID     *string               `json:"client_id"`
	TraceID      string                `json:"trace_id"`
	SpanID       string                `json:"span_id"`
	ParentSpanID *string               `json:"parent_span_id"`
	SessionID    *string               `json:"session_id"`
	EventContext AuditLogContext       `json:"event_context"`
}

// AuditLogWithPayload is returned by the single-event endpoint and includes
// the full event payload. List responses use AuditLog which omits payload.
type AuditLogWithPayload struct {
	AuditLog
	Payload map[string]any `json:"payload"`
}

type AuditLogList struct {
	APIResource
	AuditLogs []*AuditLog               `json:"data"`
	Cursor    *ExtendedPaginationCursor `json:"cursor"`
}

// NextPageStatus is the tri-state value for has_next_page in cursor pagination.
// It is "true" when there is a next page, "false" when there is not, and
// "unknown" when the query was bounded by a time window and more results may
// exist beyond that window.
type NextPageStatus string

const (
	NextPageTrue    NextPageStatus = "true"
	NextPageFalse   NextPageStatus = "false"
	NextPageUnknown NextPageStatus = "unknown"
)

// AuditLogFilterMatch controls how the string filters on the audit logs list
// endpoint (subject, type, actor, trace_id, client_id, impersonator_user_id)
// are combined when more than one is supplied.
//
//   - AuditLogFilterMatchAll (default): every supplied filter must match.
//     Filters are joined with AND.
//   - AuditLogFilterMatchAny: an event matches if at least one of the
//     supplied filters matches. Filters are joined with OR.
//
// Time bounds (event_time_after, event_time_before) and end_user_facing_only
// are always ANDed with the (possibly OR-ed) string filter group regardless
// of the value.
type AuditLogFilterMatch string

const (
	AuditLogFilterMatchAll AuditLogFilterMatch = "all"
	AuditLogFilterMatchAny AuditLogFilterMatch = "any"
)

type ExtendedPaginationCursor struct {
	StartingAfter *string `json:"starting_after"`
	EndingBefore  *string `json:"ending_before"`
	// Deprecated: Use NextPageStatus instead. HasNextPage is kept for backwards
	// compatibility and is derived as NextPageStatus != NextPageFalse.
	HasNextPage    bool           `json:"has_next_page"`
	NextPageStatus NextPageStatus `json:"next_page_status"`
	// RetentionLimitReached is true when next_page_status is "false" and the limit
	// was the plan's retention period, not the caller's event_time_after bound.
	// When true, extending event_time_after further back will not yield more results.
	RetentionLimitReached bool `json:"retention_limit_reached"`
}

// PaginationCursor contains the cursors for pagination.
type PaginationCursor struct {
	StartingAfter *string `json:"starting_after"`
	EndingBefore  *string `json:"ending_before"`
	HasNextPage   bool    `json:"has_next_page"`
}

type AuditLogContext struct {
	Environment *EnvironmentContext `json:"environment"`
	DeviceInfo  *DeviceContext      `json:"device"`
}

type EnvironmentContext struct {
	Type          *string             `json:"type"`
	Application   *ApplicationContext `json:"application"`
	Domain        *DomainContext      `json:"domain"`
	PrimaryDomain *DomainContext      `json:"primary_domain"`
}

type ApplicationContext struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type DomainContext struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type DeviceContext struct {
	IPAddress      *string          `json:"ip_address"`
	UserAgent      *string          `json:"user_agent"`
	Browser        *BrowserContext  `json:"browser"`
	DeviceType     *string          `json:"device_type"`
	IsMobile       *bool            `json:"is_mobile"`
	Location       *LocationContext `json:"location"`
	ClerkJSVersion *string          `json:"clerk_js_version"`
	IsNative       *bool            `json:"is_native"`
}

type BrowserContext struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

type LocationContext struct {
	City    *string `json:"city"`
	Country *string `json:"country"`
}

type ImpersonatorResponse struct {
	UserID *string `json:"user_id"`
}
