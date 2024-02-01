// Package clerk provides a way to communicate with the Clerk API.
// Includes types for Clerk API requests, responses and all
// available resources.
package clerk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	sdkVersion      string = "v2.0.0"
	clerkAPIVersion string = "v1"
)

const (
	// APIURL is the base URL for the Clerk API.
	APIURL string = "https://api.clerk.com"
)

// The Clerk secret key. Configured on a package level.
var secretKey string

// SetKey sets the Clerk API key.
func SetKey(key string) {
	secretKey = key
}

// APIResource describes a Clerk API resource and contains fields and
// methods common to all resources.
type APIResource struct {
	Response *APIResponse `json:"-"`
}

// Read sets the response on the resource.
func (r *APIResource) Read(response *APIResponse) {
	r.Response = response
}

// APIParams implements functionality that's common to all types
// that can be used as API request parameters.
// It is recommended to embed this type to all types that will be
// used for API operation parameters.
type APIParams struct {
}

// Add can be used to set parameters to url.Values. The method
// is currently a no-op, but is defined so that all types that
// describe API operation parameters implement the Queryable
// interface.
func (params *APIParams) Add(q url.Values) {
}

// APIResponse describes responses coming from the Clerk API.
// Exposes some commonly used HTTP response fields along with
// the raw data in the response body.
type APIResponse struct {
	Header     http.Header
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200

	// TraceID is a unique identifier for tracing the origin of the
	// response.
	// Useful for debugging purposes.
	TraceID string
	// RawJSON contains the response body as raw bytes.
	RawJSON json.RawMessage
}

// Success returns true for API response status codes in the
// 200-399 range, false otherwise.
func (resp *APIResponse) Success() bool {
	return resp.StatusCode < 400
}

// NewAPIResponse creates an APIResponse from the passed http.Response
// and the raw response body.
func NewAPIResponse(resp *http.Response, body json.RawMessage) *APIResponse {
	return &APIResponse{
		Header:     resp.Header,
		TraceID:    resp.Header.Get("Clerk-Trace-Id"),
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		RawJSON:    body,
	}
}

// APIRequest describes requests to the Clerk API.
type APIRequest struct {
	Method string
	Path   string
	Params Queryable
}

// SetParams sets the APIRequest.Params.
func (req *APIRequest) SetParams(params Queryable) {
	req.Params = params
}

// NewAPIRequest creates an APIRequest with the provided HTTP method
// and path.
func NewAPIRequest(method, path string) *APIRequest {
	return &APIRequest{
		Method: method,
		Path:   path,
	}
}

// Backend is the primary interface for communicating with the Clerk
// API.
type Backend interface {
	// Call makes requests to the Clerk API.
	Call(context.Context, *APIRequest, ResponseReader) error
}

// ResponseReader reads Clerk API responses.
type ResponseReader interface {
	Read(*APIResponse)
}

// Queryable can add parameters to url.Values.
// Useful for constructing a request query string.
type Queryable interface {
	Add(url.Values)
}

// BackendConfig is used to configure a new Clerk Backend.
type BackendConfig struct {
	// HTTPClient is an HTTP client instance that will be used for
	// making API requests.
	// If it's not set a default HTTP client will be used.
	HTTPClient *http.Client
	// URL is the base URL to use for API endpoints.
	// If it's not set, the default value for the Backend will be used.
	URL *string
}

// NewBackend returns a default backend implementation with the
// provided configuration.
// Please note that the return type is an interface because the
// Backend is not supposed to be used directly.
func NewBackend(config *BackendConfig) Backend {
	if config.HTTPClient == nil {
		config.HTTPClient = defaultHTTPClient
	}
	if config.URL == nil {
		config.URL = String(APIURL)
	}
	return &defaultBackend{
		HTTPClient: config.HTTPClient,
		URL:        *config.URL,
	}
}

// GetBackend returns the library's supported backend for the Clerk
// API.
func GetBackend() Backend {
	var b Backend

	backend.mu.RLock()
	b = backend.Backend
	backend.mu.RUnlock()

	if b != nil {
		return b
	}

	b = NewBackend(&BackendConfig{})
	SetBackend(b)
	return b
}

// SetBackend sets the Backend that will be used to make requests
// to the Clerk API.
// Use this method if you need to override the default Backend
// configuration.
func SetBackend(b Backend) {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	backend.Backend = b
}

type defaultBackend struct {
	HTTPClient *http.Client
	URL        string
}

// Call sends requests to the Clerk API and handles the responses.
func (b *defaultBackend) Call(ctx context.Context, apiReq *APIRequest, setter ResponseReader) error {
	req, err := b.newRequest(ctx, apiReq)
	if err != nil {
		return err
	}

	return b.do(req, apiReq.Params, setter)
}

func (b *defaultBackend) newRequest(ctx context.Context, apiReq *APIRequest) (*http.Request, error) {
	path, err := url.JoinPath(b.URL, clerkAPIVersion, apiReq.Path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, apiReq.Method, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", fmt.Sprintf("Clerk/%s SDK-Go/%s", clerkAPIVersion, sdkVersion))
	req.Header.Add("X-Clerk-SDK", fmt.Sprintf("go/%s", sdkVersion))
	req = req.WithContext(ctx)

	return req, nil
}

func (b *defaultBackend) do(req *http.Request, params Queryable, setter ResponseReader) error {
	err := setRequestBody(req, params)
	if err != nil {
		return err
	}

	resp, err := b.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	resBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	apiResponse := NewAPIResponse(resp, resBody)
	// Looks like something went wrong. Handle the error.
	if !apiResponse.Success() {
		return handleError(apiResponse, resBody)
	}

	setter.Read(apiResponse)
	err = json.Unmarshal(resBody, setter)
	if err != nil {
		return err
	}

	return nil
}

// Sets the params in either the request body, or the querystring
// for GET requests.
func setRequestBody(req *http.Request, params Queryable) error {
	// GET requests don't have a body, but we will pass the params
	// in the query string.
	if req.Method == http.MethodGet && params != nil {
		q := req.URL.Query()
		params.Add(q)
		req.URL.RawQuery = q.Encode()
		return nil
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	return nil
}

// Error response handling
func handleError(resp *APIResponse, body []byte) error {
	apiError := &APIErrorResponse{
		HTTPStatusCode: resp.StatusCode,
	}
	apiError.Read(resp)
	err := json.Unmarshal(body, apiError)
	if err != nil || apiError.Errors == nil {
		// This is probably not an expected API error.
		// Return the raw server response.
		return errors.New(string(body))
	}
	return apiError
}

// The active Backend
var backend api

// This type is a container for a Backend. Guarantees thread-safe
// access to the current Backend.
type api struct {
	Backend Backend
	mu      sync.RWMutex
}

// defaultHTTPTimeout is the default timeout on the http.Client used
// by the library.
const defaultHTTPTimeout = 5 * time.Second

// The default HTTP client used for communication with the Clerk API.
var defaultHTTPClient = &http.Client{
	Timeout: defaultHTTPTimeout,
}

// APIErrorResponse is used for cases where requests to the Clerk
// API result in error responses.
type APIErrorResponse struct {
	APIResource

	Errors []Error `json:"errors"`

	HTTPStatusCode int    `json:"status,omitempty"`
	TraceID        string `json:"clerk_trace_id,omitempty"`
}

// Error returns the marshaled representation of the APIErrorResponse.
func (resp *APIErrorResponse) Error() string {
	ret, err := json.Marshal(resp)
	if err != nil {
		// This shouldn't happen, let's return the raw response
		return string(resp.Response.RawJSON)
	}
	return string(ret)
}

// Error is a representation of a single error that can occur in the
// Clerk API.
type Error struct {
	Code        string          `json:"code"`
	Message     string          `json:"message"`
	LongMessage string          `json:"long_message"`
	Meta        json.RawMessage `json:"meta,omitempty"`
}

// ListParams holds fields that are common for list API operations.
type ListParams struct {
	Limit  *int64 `json:"limit,omitempty"`
	Offset *int64 `json:"offset,omitempty"`
}

// Add sets list params to the passed in url.Values.
func (params ListParams) Add(q url.Values) {
	if params.Limit != nil {
		q.Set("limit", strconv.FormatInt(*params.Limit, 10))
	}
	if params.Offset != nil {
		q.Set("offset", strconv.FormatInt(*params.Offset, 10))
	}
}

// String returns a pointer to the provided string value.
func String(v string) *string {
	return &v
}

// Int64 returns a pointer to the provided int64 value.
func Int64(v int64) *int64 {
	return &v
}
