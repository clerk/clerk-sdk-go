package clerk

type AuditLog struct {
	APIResource
	Object          string         `json:"object"`
	ID              string         `json:"id"`
	EventTime       int64          `json:"event_time"`
	SubjectInstance string         `json:"subject_instance"`
	Actor           string         `json:"actor"`
	Subject         string         `json:"subject"`
	Type            string         `json:"type"`
	TraceID         string         `json:"trace_id"`
	SpanID          string         `json:"span_id"`
	ParentSpanID    *string        `json:"parent_span_id"`
	Payload         map[string]any `json:"payload"`
}

type AuditLogList struct {
	APIResource
	AuditLogs []*AuditLog              `json:"data"`
	Cursor    *CursorPaginationCursors `json:"cursor"`
}

// CursorPaginationCursors contains the cursors for pagination.
type CursorPaginationCursors struct {
	StartingAfter *string `json:"starting_after"`
	EndingBefore  *string `json:"ending_before"`
}
