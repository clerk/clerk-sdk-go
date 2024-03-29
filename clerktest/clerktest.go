// Package clerktest provides utilities for testing.
package clerktest

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/require"
)

// RoundTripper can be used as a mock Transport for http.Clients.
// Set the RoundTripper's fields accordingly to determine the
// response or perform assertions on the http.Request properties.
type RoundTripper struct {
	T *testing.T
	// Status is the response Status code.
	Status int
	// Out is the response body.
	Out json.RawMessage
	// Set this field to assert on the request method.
	Method string
	// Set this field to assert that the request path matches.
	Path string
	// Set this field to assert that the request URL querystring matches.
	Query *url.Values
	// Set this field to assert that the request body matches.
	In json.RawMessage
}

// RoundTrip returns an http.Response based on the RoundTripper's fields.
// It will also perform assertions on the http.Request.
func (rt *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt.Status == 0 {
		rt.Status = http.StatusOK
	}
	if rt.Method != "" {
		require.Equal(rt.T, rt.Method, r.Method)
	}
	if rt.Path != "" {
		require.Equal(rt.T, rt.Path, r.URL.Path)
	}
	if rt.Query != nil {
		require.Equal(rt.T, rt.Query.Encode(), r.URL.Query().Encode())
	}
	if rt.In != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		require.JSONEq(rt.T, string(rt.In), string(body))
	}
	return &http.Response{
		StatusCode: rt.Status,
		Body:       io.NopCloser(bytes.NewReader(rt.Out)),
	}, nil
}

// GenerateJWT creates a JSON web token with the provided claims
// and key ID.
func GenerateJWT(t *testing.T, claims any, kid string) (string, crypto.PublicKey) {
	t.Helper()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	signerOpts := &jose.SignerOptions{}
	signerOpts.WithType("JWT")
	if kid != "" {
		signerOpts.WithHeader("kid", kid)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privKey}, signerOpts)
	require.NoError(t, err)

	builder := jwt.Signed(signer)
	builder = builder.Claims(claims)
	token, err := builder.CompactSerialize()
	require.NoError(t, err)

	return token, privKey.Public()
}

// Clock provides a test clock which can be manually advanced through time.
type Clock struct {
	mu sync.RWMutex
	// The current time of this test clock.
	time time.Time
}

// NewClockAt returns a Clock initialized at the given time.
func NewClockAt(t time.Time) *Clock {
	return &Clock{time: t}
}

// Now returns the clock's current time.
func (c *Clock) Now() time.Time {
	return c.time
}

// Advance moves the test clock to a new point in time.
func (c *Clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.time = c.time.Add(d)
}
