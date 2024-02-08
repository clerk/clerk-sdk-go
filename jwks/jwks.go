// Package jwks provides access to the JWKS endpoint.
package jwks

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

const path = "/jwks"

type GetParams struct {
	clerk.APIParams
}

// Get retrieves a JSON Web Key set.
func Get(ctx context.Context, params *GetParams) (*clerk.JSONWebKeySet, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)

	set := &clerk.JSONWebKeySet{}
	err := clerk.GetBackend().Call(ctx, req, set)
	return set, err
}
