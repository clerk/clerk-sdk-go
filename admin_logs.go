package clerk

// AdminLog is a single admin-log event returned by the list endpoint.
// Admin-log events describe actions taken by a workspace member (a "C1"
// actor, e.g. a Clerk Dashboard user) on a workspace's resources, such as
// creating an instance key or rotating a secret. They are scoped to the
// workspace that owns the authenticated instance's application; events for
// sibling instances inside the same workspace are visible too, which is why
// AdminLog carries both Workspace and Instance as top-level fields.
type AdminLog struct {
	APIResource
	Object       string                `json:"object"`
	ID           string                `json:"id"`
	EventTime    int64                 `json:"event_time"`
	Workspace    string                `json:"workspace"`
	Instance     string                `json:"instance"`
	Actor        string                `json:"actor"`
	Impersonator *ImpersonatorResponse `json:"impersonator"`
	Subject      string                `json:"subject"`
	Type         string                `json:"type"`
	Source       *string               `json:"source"`
	ClientID     *string               `json:"client_id"`
	TraceID      string                `json:"trace_id"`
	SpanID       string                `json:"span_id"`
	ParentSpanID *string               `json:"parent_span_id"`
	SessionID    *string               `json:"session_id"`
	// EventContext shares the wire shape of [AuditLog.EventContext]; the
	// underlying response struct is reused on the server side, so the SDK
	// reuses the same type rather than duplicating it.
	EventContext AuditLogContext `json:"event_context"`
}

// AdminLogWithPayload is returned by the single-event endpoint and includes
// the full event payload. List responses use [AdminLog] which omits payload.
type AdminLogWithPayload struct {
	AdminLog
	Payload map[string]any `json:"payload"`
}

type AdminLogList struct {
	APIResource
	AdminLogs []*AdminLog               `json:"data"`
	Cursor    *ExtendedPaginationCursor `json:"cursor"`
}

// AdminLogFilterMatch controls how the string filters on the admin logs list
// endpoint (instance, subject, type, actor, trace_id, client_id,
// impersonator_user_id) are combined when more than one is supplied.
//
//   - [AdminLogFilterMatchAll] (default): every supplied filter must match.
//     Filters are joined with AND.
//   - [AdminLogFilterMatchAny]: an event matches if at least one of the
//     supplied filters matches. Filters are joined with OR.
//
// Time bounds (event_time_after, event_time_before) and end_user_facing_only
// are always ANDed with the (possibly OR-ed) string filter group regardless
// of the value.
//
// The wire values match [AuditLogFilterMatch] ("all" / "any"); the type is
// kept distinct so the constants are namespaced to their endpoint.
type AdminLogFilterMatch string

const (
	AdminLogFilterMatchAll AdminLogFilterMatch = "all"
	AdminLogFilterMatchAny AdminLogFilterMatch = "any"
)
