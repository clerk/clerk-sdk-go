// Package domain provides the Domains API.
package domain

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

const path = "/domains"

type CreateParams struct {
	clerk.APIParams
	Name        *string `json:"name,omitempty"`
	ProxyURL    *string `json:"proxy_url,omitempty"`
	IsSatellite *bool   `json:"is_satellite,omitempty"`
}

// Create creates a new domain.
func Create(ctx context.Context, params *CreateParams) (*clerk.Domain, error) {
	req := clerk.NewAPIRequest(http.MethodPost, path)
	req.SetParams(params)

	domain := &clerk.Domain{}
	err := clerk.GetBackend().Call(ctx, req, domain)
	return domain, err
}

type UpdateParams struct {
	clerk.APIParams
	Name     *string `json:"name,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
}

// Update updates a domain's properties.
func Update(ctx context.Context, id string, params *UpdateParams) (*clerk.Domain, error) {
	path, err := url.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPatch, path)
	req.SetParams(params)

	domain := &clerk.Domain{}
	err = clerk.GetBackend().Call(ctx, req, domain)
	return domain, err
}

// Delete removes a domain.
func Delete(ctx context.Context, id string) (*clerk.DeletedResource, error) {
	path, err := url.JoinPath(path, id)
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodDelete, path)
	domain := &clerk.DeletedResource{}
	err = clerk.GetBackend().Call(ctx, req, domain)
	return domain, err
}

type ListParams struct {
	clerk.APIParams
}

// List returns a list of domains
func List(ctx context.Context, params *ListParams) (*clerk.DomainList, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	list := &clerk.DomainList{}
	err := clerk.GetBackend().Call(ctx, req, list)
	return list, err
}
