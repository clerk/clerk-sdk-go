// Package template provides the Templates API.
package template

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/templates"

// Client is used to invoke the Templates API.
type Client struct {
	Backend clerk.Backend
}

type ClientConfig struct {
	clerk.BackendConfig
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type GetParams struct {
	TemplateType clerk.TemplateType `json:"-"`
	Slug         string             `json:"-"`
}

// Get retrieves details for a template.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.Template, error) {
	path, err := clerk.JoinPath(path, string(params.TemplateType), params.Slug)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	resource := &clerk.Template{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type UpdateParams struct {
	clerk.APIParams
	Name             *string            `json:"name,omitempty"`
	Subject          *string            `json:"subject,omitempty"`
	Markup           *string            `json:"markup,omitempty"`
	Body             *string            `json:"body,omitempty"`
	FromEmailName    *string            `json:"from_email_name,omitempty"`
	DeliveredByClerk *bool              `json:"delivered_by_clerk,omitempty"`
	TemplateType     clerk.TemplateType `json:"-"`
	Slug             string             `json:"-"`
}

// Update updates an existing template or creates a new one with the
// provided params.
func (c *Client) Update(ctx context.Context, params *UpdateParams) (*clerk.Template, error) {
	path, err := clerk.JoinPath(path, string(params.TemplateType), params.Slug)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPut, path)
	req.SetParams(params)
	resource := &clerk.Template{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type DeleteParams struct {
	TemplateType clerk.TemplateType `json:"-"`
	Slug         string             `json:"-"`
}

// Delete deletes a custom user template.
func (c *Client) Delete(ctx context.Context, params *DeleteParams) (*clerk.DeletedResource, error) {
	path, err := clerk.JoinPath(path, string(params.TemplateType), params.Slug)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	resource := &clerk.DeletedResource{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type RevertParams struct {
	TemplateType clerk.TemplateType `json:"-"`
	Slug         string             `json:"-"`
}

// Revert reverts a template to its default state.
func (c *Client) Revert(ctx context.Context, params *RevertParams) (*clerk.Template, error) {
	path, err := clerk.JoinPath(path, string(params.TemplateType), params.Slug, "/revert")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	resource := &clerk.Template{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ToggleDeliveryParams struct {
	clerk.APIParams
	DeliveredByClerk *bool              `json:"delivered_by_clerk,omitempty"`
	TemplateType     clerk.TemplateType `json:"-"`
	Slug             string             `json:"-"`
}

// ToggleDelivery sets the delivery by Clerk for a template.
func (c *Client) ToggleDelivery(ctx context.Context, params *ToggleDeliveryParams) (*clerk.Template, error) {
	path, err := clerk.JoinPath(path, string(params.TemplateType), params.Slug, "/toggle_delivery")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.Template{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type PreviewParams struct {
	clerk.APIParams
	Subject       *string            `json:"subject,omitempty"`
	Body          *string            `json:"body,omitempty"`
	FromEmailName *string            `json:"from_email_name,omitempty"`
	TemplateType  clerk.TemplateType `json:"-"`
	Slug          string             `json:"-"`
}

// Preview returns a preview of a template.
func (c *Client) Preview(ctx context.Context, params *PreviewParams) (*clerk.TemplatePreview, error) {
	path, err := clerk.JoinPath(path, string(params.TemplateType), params.Slug, "/preview")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)
	resource := &clerk.TemplatePreview{}
	err = c.Backend.Call(ctx, req, resource)
	return resource, err
}

type ListParams struct {
	TemplateType clerk.TemplateType `json:"-"`
}

// List returns a list of templates of a given type.
func (c *Client) List(ctx context.Context, params *ListParams) (*clerk.TemplateList, error) {
	path, err := clerk.JoinPath(path, fmt.Sprintf("%s?paginated=true", params.TemplateType))
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, path)
	list := &clerk.TemplateList{}
	err = c.Backend.Call(ctx, req, list)
	return list, err
}
