package waitlistentry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/waitlist_entries"

// Client is used to invoke the Waitlist Entries API.
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
	clerk.ListParams
	OrderBy  *string  `json:"order_by,omitempty"`
	Query    *string  `json:"query,omitempty"`
	Statuses []string `json:"status,omitempty"`
}

// ToQuery returns query string values from the params.
func (params *ListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.OrderBy != nil {
		q.Set("order_by", *params.OrderBy)
	}
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	for _, status := range params.Statuses {
		q.Add("status", status)
	}
	return q
}

// List returns all waitlist entries.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.WaitlistEntriesList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)
	list := &clerk.WaitlistEntriesList{}
	err := c.Backend.Call(ctx, req, list)
	return list, err
}

type CreateParams struct {
	clerk.APIParams
	EmailAddress string `json:"email_address"`
	Notify       *bool  `json:"notify,omitempty"`
}

// Create adds a new waitlist entry.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*clerk.WaitlistEntry, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	invitation := &clerk.WaitlistEntry{}
	err := c.Backend.Call(ctx, req, invitation)
	return invitation, err
}

type BulkCreateParams struct {
	clerk.APIParams
	WaitlistEntries []*CreateParams
}

func (b BulkCreateParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.WaitlistEntries)
}

type bulkCreateResponse struct {
	clerk.APIResource
	WaitlistEntries []*clerk.WaitlistEntry
}

func (b *bulkCreateResponse) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &b.WaitlistEntries)
}

// BulkCreate creates multiple waitlist entries.
func (c *Client) BulkCreate(ctx context.Context, params *BulkCreateParams) (*clerk.WaitlistEntries, error) {
	path, err := clerk.JoinPath(path, "bulk")
	if err != nil {
		return nil, err
	}

	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)

	res := &bulkCreateResponse{}
	if err := c.Backend.Call(ctx, req, res); err != nil {
		return nil, err
	}

	return &clerk.WaitlistEntries{
		APIResource:     res.APIResource,
		WaitlistEntries: res.WaitlistEntries,
	}, nil
}
