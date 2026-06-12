// Package admin_logs provides the Admin Logs API
package admin_logs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/admin_logs"

// Client is used to invoke the Admin Logs API.
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
	// Filter admin logs by instance. Admin events span sibling instances inside
	// a workspace (unlike audit logs, where every event already belongs to a
	// single instance), so the list endpoint exposes this as a first-class
	// filter and the SDK mirrors it here.
	Instance *string `json:"instance,omitempty"`
	// Filter admin logs by subject.
	Subject *string `json:"subject,omitempty"`
	// Filter admin logs by actor.
	Actor *string `json:"actor,omitempty"`
	// Filter admin logs by trace ID.
	TraceID *string `json:"trace_id,omitempty"`
	// Filter admin logs by event type (e.g., instance_key.created).
	// A trailing "*" is treated as a prefix match (e.g. "instance_key.*").
	Type *string `json:"type,omitempty"`
	// Filter admin logs by client ID.
	ClientID *string `json:"client_id,omitempty"`
	// Filter admin logs by impersonator user ID.
	ImpersonatorUserID *string `json:"impersonator_user_id,omitempty"`
	// FilterMatch controls how Instance, Subject, Type, Actor, TraceID,
	// ClientID, and ImpersonatorUserID are combined when more than one is
	// supplied.
	//
	//   - AdminLogFilterMatchAll (default, also when nil): every supplied
	//     filter must match (AND).
	//   - AdminLogFilterMatchAny: an event matches if at least one of the
	//     supplied filters matches (OR).
	//
	// EventTimeAfter, EventTimeBefore, and EndUserFacingOnly are always
	// ANDed with the (possibly OR-ed) string filter group regardless of
	// FilterMatch.
	FilterMatch *clerk.AdminLogFilterMatch `json:"filter_match,omitempty"`
	// Filter admin logs to events on or after this date (Unix timestamp in milliseconds).
	EventTimeAfter *int64 `json:"event_time_after,omitempty"`
	// Filter admin logs to events on or before this date (Unix timestamp in milliseconds).
	EventTimeBefore *int64 `json:"event_time_before,omitempty"`
	// When true, only returns events marked as end-user facing.
	// When false or omitted, returns all events.
	EndUserFacingOnly *bool `json:"end_user_facing_only,omitempty"`
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
	if params.Instance != nil {
		q.Add("instance", *params.Instance)
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
	return q
}

type GetParams struct {
	EventTimeMs int64
	EventID     string
}

// Get retrieves a single admin log by its composite key (event_time_ms:event_id).
// Unlike List, the returned object includes the full event payload.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.AdminLogWithPayload, error) {
	compositeID := fmt.Sprintf("%d:%s", params.EventTimeMs, params.EventID)
	path, err := clerk.JoinPath(path, compositeID)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.AdminLogWithPayload{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

// List returns a list of admin logs.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.AdminLogList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.AdminLogList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
