// Package audit_logs provides the Audit Logs API
package audit_logs

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/clerk/clerk-sdk-go/v3"
)

//go:generate go run ../cmd/gen/main.go

const path = "/audit_logs"

// Client is used to invoke the Audit Logs API.
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
	// A cursor for pagination. Provide the cursor from a previous response to fetch the next page.
	StartingAfter *string `json:"starting_after,omitempty"`
	// A cursor for pagination. Provide the cursor from a previous response to fetch the previous page.
	EndingBefore *string `json:"ending_before,omitempty"`
	// Filter audit logs by subject (user ID or organization ID).
	Subject *string `json:"subject,omitempty"`
	// Filter audit logs by event type (e.g., email_send).
	Type *string `json:"type,omitempty"`
	// Filter audit logs to events on or after this date (Unix timestamp in milliseconds).
	EventTimeAfter *int64 `json:"event_time_after,omitempty"`
	// Filter audit logs to events on or before this date (Unix timestamp in milliseconds).
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
	if params.Subject != nil {
		q.Add("subject", *params.Subject)
	}
	if params.Type != nil {
		q.Add("type", *params.Type)
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

// List returns a list of audit logs.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.AuditLogList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.AuditLogList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}
