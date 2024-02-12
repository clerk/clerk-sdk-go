// Package jwks provides access to the JWKS endpoint.
package jwks

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/jwks"

// Client is used to invoke the JWKS endpoint.
type Client struct {
	Backend clerk.Backend
	// Cache is a store for JSON Web Key sets attached to this Client.
	// The clerk.JSONWebKeySet that is returned by [Get] is not likely
	// to change often. We can take advantage of the fact and add
	// caching capabilities to a Client which is intended to call
	// Get often.
	// An example of frequent calls to Get is the
	// [http.WithHeaderAuthorization] middleware. In an HTTP server
	// the middleware needs the JSON Web Key Set to verify the JSON
	// Web Token for every HTTP request. Using a Client with a Cache
	// improves latency by caching Get responses.
	//
	//   client := jwks.NewClient(&jwks.Config{
	//     EnableCache: true,
	//   })
	//   http.WithHeaderAuthorization(http.JWKSClient(client))
	//
	Cache Cache
}

type ClientConfig struct {
	clerk.BackendConfig
	// EnableCache can initialize a cache for the Client when set to
	// true. The Cache will be used by the SDK's HTTP middleware to
	// improve performance. The JWKS endpoint is not likely to change
	// often. Enabling a caching layer allows consumers of the Client
	// to take advantage of its Cache to store the JSON Web Key Set
	// from the API.
	// Please note that the Client does not use the Cache by default,
	// and it's up to the consumer to decide how and what to cache.
	// See [jwt.VerifyToken] for an example of how the Cache can be
	// used.
	EnableCache bool
}

func NewClient(config *ClientConfig) *Client {
	client := &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
	if config.EnableCache {
		client.Cache = &defaultCache{}
	}
	return client
}

type GetParams struct {
	clerk.APIParams
}

// Get retrieves a JSON Web Key set.
func (c *Client) Get(ctx context.Context, params *GetParams) (*clerk.JSONWebKeySet, error) {
	req := clerk.NewAPIRequest(http.MethodGet, path)
	req.SetParams(params)

	set := &clerk.JSONWebKeySet{}
	err := clerk.GetBackend().Call(ctx, req, set)
	return set, err
}

// Cache provides a store for JWKS values.
type Cache interface {
	// Get retrieves an entry from the cache.
	Get() *clerk.JSONWebKeySet
	// Set sets an entry in the cache.
	Set(*clerk.JSONWebKeySet, time.Time)
}

// Caching store for JSON Web Key Sets.
type defaultCache struct {
	mu        sync.RWMutex
	value     *clerk.JSONWebKeySet
	expiresAt time.Time
}

// Get returns the *clerk.JSONWebKeySet that's stored in the cache,
// unless the cache has expired.
func (c *defaultCache) Get() *clerk.JSONWebKeySet {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.expiresAt.Before(time.Now().UTC()) {
		return nil
	}
	return c.value
}

// Set adds a new entry with the provided value in the cache.
// An expiration date will be set for the entry.
func (c *defaultCache) Set(value *clerk.JSONWebKeySet, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = value
	c.expiresAt = expiresAt
}

var cacheInit sync.Once

// A "singleton" cache for the package.
var cache Cache

// GetCache returns the library's default cache singleton.
// Please note that the returned Cache is a package-level variable.
// Using the package with more than one Clerk API secret keys might
// require to use different Clients with their own Cache.
func GetCache() Cache {
	cacheInit.Do(func() {
		cache = &defaultCache{}
	})
	return cache
}
