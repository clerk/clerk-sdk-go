// Package clerktest provides utilities for testing.
package clerktest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

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
