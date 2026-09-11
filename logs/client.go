// Package logs provides the Logs API
package logs

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/logs"

// Client is used to invoke the Logs API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type ListParams struct {
	clerk.APIParams
	// Applies a limit to the number of results returned.
	// Can be used for paginating the results together with StartingAfter or EndingBefore.
	Limit *int64 `json:"limit,omitempty"`
	// A cursor for pagination. Provide the starting_after cursor from a previous response to fetch the next page.
	StartingAfter *string `json:"starting_after,omitempty"`
	// A cursor for pagination. Provide the ending_before cursor from a previous response to fetch the previous page.
	EndingBefore *string `json:"ending_before,omitempty"`
	// Filter logs by subject.
	Subject *string `json:"subject,omitempty"`
	// Filter logs by actor.
	Actor *string `json:"actor,omitempty"`
	// Filter logs by trace ID.
	TraceID *string `json:"trace_id,omitempty"`
	// Filter logs by event type (e.g., email_send).
	Type *string `json:"type,omitempty"`
	// Filter logs by client ID.
	ClientID *string `json:"client_id,omitempty"`
	// Filter logs by impersonator user ID.
	ImpersonatorUserID *string `json:"impersonator_user_id,omitempty"`
	// Filter logs by device IP address (exact match against
	// device_info_ip_address).
	IPAddress *string `json:"ip_address,omitempty"`
	// FilterMatch controls how Subject, Type, Actor, TraceID, ClientID,
	// ImpersonatorUserID, IPAddress, and PayloadFilters are combined when
	// more than one is supplied.
	//
	//   - LogFilterMatchAll (default, also when nil): every supplied
	//     filter must match (AND).
	//   - LogFilterMatchAny: an event matches if at least one of the
	//     supplied filters matches (OR).
	//
	// EventTimeAfter, EventTimeBefore, and EndUserFacingOnly are always
	// ANDed with the (possibly OR-ed) string filter group regardless of
	// FilterMatch.
	FilterMatch *clerk.LogFilterMatch `json:"filter_match,omitempty"`
	// Filter logs to events on or after this date (Unix timestamp in milliseconds).
	EventTimeAfter *int64 `json:"event_time_after,omitempty"`
	// Filter logs to events on or before this date (Unix timestamp in milliseconds).
	EventTimeBefore *int64 `json:"event_time_before,omitempty"`
	// When true, only returns events marked as end-user facing.
	// When false or omitted, returns all events.
	EndUserFacingOnly *bool `json:"end_user_facing_only,omitempty"`
	// PayloadFilters are exact-match filters on event payload fields,
	// keyed by the field's dot-path (nested fields use dots, e.g.
	// "captcha_attempt_payload.provider"). Serialized as
	// payload_filter[<path>]=<value>, one value per field.
	//
	// Requires Type. Fields are validated by the API against the payload
	// schema of that event type (discoverable via GetSchema); an unknown
	// field returns a 422 listing the allowed fields. Values must parse as
	// the field's type and match exactly; an event whose payload lacks the
	// field never matches. At most 10 filters per request. Payload filters
	// join the same FilterMatch group as the string filters.
	PayloadFilters map[string]string `json:"-"`
	// PayloadFields are payload field dot-paths to include per returned
	// log. Requires Type and is validated like PayloadFilters. Each
	// returned log gains a partial Payload map containing only the
	// requested fields. Not allowed together with LogFilterMatchAny.
	PayloadFields []string `json:"-"`
}

// ToQuery returns the params as url.Values.
func (params *ListParams) ToQuery() url.Values {
	q := url.Values{}
	if params.Limit != nil {
		q.Add("limit", strconv.FormatInt(*params.Limit, 10))
	}
	if params.StartingAfter != nil {
		q.Add("starting_after", *params.StartingAfter)
	}
	if params.EndingBefore != nil {
		q.Add("ending_before", *params.EndingBefore)
	}
	if params.Subject != nil {
		q.Add("subject", *params.Subject)
	}
	if params.Actor != nil {
		q.Add("actor", *params.Actor)
	}
	if params.TraceID != nil {
		q.Add("trace_id", *params.TraceID)
	}
	if params.Type != nil {
		q.Add("type", *params.Type)
	}
	if params.ClientID != nil {
		q.Add("client_id", *params.ClientID)
	}
	if params.ImpersonatorUserID != nil {
		q.Add("impersonator_user_id", *params.ImpersonatorUserID)
	}
	if params.IPAddress != nil {
		q.Add("ip_address", *params.IPAddress)
	}
	if params.FilterMatch != nil {
		q.Add("filter_match", string(*params.FilterMatch))
	}
	if params.EventTimeAfter != nil {
		q.Add("event_time_after", strconv.FormatInt(*params.EventTimeAfter, 10))
	}
	if params.EventTimeBefore != nil {
		q.Add("event_time_before", strconv.FormatInt(*params.EventTimeBefore, 10))
	}
	if params.EndUserFacingOnly != nil {
		q.Add("end_user_facing_only", strconv.FormatBool(*params.EndUserFacingOnly))
	}
	// Sorted so the query string is deterministic regardless of map order.
	for _, path := range slices.Sorted(maps.Keys(params.PayloadFilters)) {
		q.Add("payload_filter["+path+"]", params.PayloadFilters[path])
	}
	if len(params.PayloadFields) > 0 {
		q.Add("payload_fields", strings.Join(params.PayloadFields, ","))
	}
	return q
}

type GetParams struct {
	EventTimeMs int64
	EventID     string
}

// Get retrieves a single log by its composite key (event_time_ms:event_id).
// Unlike List, the returned object includes the full event payload.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.LogWithPayload, error) {
	compositeID := fmt.Sprintf("%d:%s", params.EventTimeMs, params.EventID)
	path, err := clerk.JoinPath(path, compositeID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.LogWithPayload{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// List returns a list of logs.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.LogList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.LogList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

type GetSchemaParams struct {
	clerk.APIParams
	// Type is the event type to describe: concrete ("sign_in.completed"),
	// trailing-* wildcard ("sign_in.*", whose fields are the intersection
	// across every matching type), or bare family name ("sign_in",
	// normalized to its wildcard form). Unknown types return a 422.
	Type string `json:"type"`
}

// ToQuery returns the params as url.Values.
func (params *GetSchemaParams) ToQuery() url.Values {
	q := url.Values{}
	q.Add("type", params.Type)
	return q
}

// GetSchema retrieves the payload schema of a log event type: the fields
// its payload carries, addressed by the dot-paths that List's PayloadFilters
// and PayloadFields accept.
func (c *Client) GetSchema(ctx context.Context, params *GetSchemaParams) (*clerk.LogSchema, error) {
	path, err := clerk.JoinPath(path, "schemas")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	resource := &clerk.LogSchema{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}
