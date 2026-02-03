package clerk

type AuditLog struct {
	APIResource
	Object       string          `json:"object"`
	ID           string          `json:"id"`
	EventTime    int64           `json:"event_time"`
	Actor        string          `json:"actor"`
	Subject      string          `json:"subject"`
	Type         string          `json:"type"`
	ClientID     *string         `json:"client_id"`
	TraceID      string          `json:"trace_id"`
	SpanID       string          `json:"span_id"`
	ParentSpanID *string         `json:"parent_span_id"`
	Payload      map[string]any  `json:"payload"`
	EventContext AuditLogContext `json:"event_context"`
}

type AuditLogList struct {
	APIResource
	AuditLogs []*AuditLog       `json:"data"`
	Cursor    *PaginationCursor `json:"cursor"`
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
