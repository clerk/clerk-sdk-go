package clerk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIResponse(t *testing.T) {
	body := []byte(`{"foo":"bar"}`)
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
		Header: http.Header(map[string][]string{
			"Clerk-Trace-Id":  {"trace-id"},
			"x-custom-header": {"custom-header"},
		}),
	}
	res := NewAPIResponse(resp, body)
	assert.Equal(t, body, []byte(res.RawJSON))
	assert.Equal(t, "200 OK", res.Status)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "trace-id", res.TraceID)
	assert.Equal(t, resp.Header, res.Header)
}

func TestNewBackend(t *testing.T) {
	defaultSecretKey := "sk_test_123"
	SetKey(defaultSecretKey)
	withDefaults, ok := NewBackend(&BackendConfig{}).(*defaultBackend)
	require.True(t, ok)
	require.NotNil(t, withDefaults.HTTPClient)
	assert.Equal(t, defaultHTTPTimeout, withDefaults.HTTPClient.Timeout)
	assert.Equal(t, APIURL, withDefaults.URL)
	assert.Equal(t, defaultSecretKey, withDefaults.Key)

	u := "https://some.other.url"
	httpClient := &http.Client{}
	secretKey := defaultSecretKey + "diff"
	config := &BackendConfig{
		URL:        &u,
		HTTPClient: httpClient,
		Key:        &secretKey,
	}
	withOverrides, ok := NewBackend(config).(*defaultBackend)
	require.True(t, ok)
	assert.Equal(t, u, withOverrides.URL)
	assert.Equal(t, httpClient, withOverrides.HTTPClient)
	assert.Equal(t, secretKey, withOverrides.Key)
}

func TestGetBackend_DataRace(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b, ok := GetBackend().(*defaultBackend)
			require.True(t, ok)
			assert.Equal(t, APIURL, b.URL)
		}()
	}
	wg.Wait()
}

func TestAPIErrorResponse(t *testing.T) {
	// API error response that is valid JSON. The error
	// string is the raw JSON.
	resp := &APIErrorResponse{
		HTTPStatusCode: 200,
		TraceID:        "trace-id",
		Errors: []Error{
			{
				Code:        "error-code",
				Message:     "message",
				LongMessage: "long message",
			},
		},
	}
	expected := fmt.Sprintf(`{
	"status":%d,
	"clerk_trace_id":"%s",
	"errors":[{"code":"%s","message":"%s","long_message":"%s"}]
}`,
		resp.HTTPStatusCode,
		resp.TraceID,
		resp.Errors[0].Code,
		resp.Errors[0].Message,
		resp.Errors[0].LongMessage,
	)
	assert.JSONEq(t, expected, resp.Error())
}

// This is how you define a Clerk API resource which is ready to be
// used by the library.
type testResource struct {
	APIResource
	ID     string `json:"id"`
	Object string `json:"object"`
}

// This is how you define types which can be used as Clerk API
// request parameters.
type testResourceParams struct {
	APIParams
	Name string `json:"name"`
}

// This is how you define a Clerk API resource which can be used in
// API operations that read a list of resources.
type testResourceList struct {
	APIResource
	Resources  []testResource `json:"data"`
	TotalCount int64          `json:"total_count"`
}

// This is how you define a type which can be used as parameters
// to a Clerk API operation that lists resources.
type testResourceListParams struct {
	APIParams
	ListParams
	Name     string
	Appended string
}

// We need to implement the Params interface.
func (params testResourceListParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	q.Set("name", params.Name)
	q.Set("appended", params.Appended)
	return q
}

func TestBackendCall_RequestHeaders(t *testing.T) {
	ctx := context.Background()
	method := http.MethodPost
	path := "/resources"
	secretKey := "sk_test_123"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, method, r.Method)
		require.Equal(t, "/"+clerkAPIVersion+path, r.URL.Path)

		// The client sets the Authorization header correctly.
		assert.Equal(t, fmt.Sprintf("Bearer %s", secretKey), r.Header.Get("Authorization"))
		// The client sets the User-Agent header.
		assert.Equal(t, "Clerk/v1 SDK-Go/v2.0.0", r.Header.Get("User-Agent"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		// The client includes a custom header with the SDK version.
		assert.Equal(t, "go/v2.0.0", r.Header.Get("X-Clerk-SDK"))

		_, err := w.Write([]byte(`{}`))
		require.NoError(t, err)
	}))
	defer ts.Close()

	// Set up a mock backend which triggers requests to our test server above.
	SetBackend(NewBackend(&BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	// Simulate usage for an API operation on a testResource.
	// We need to initialize a request and use the Backend to send it.
	SetKey(secretKey)
	req := NewAPIRequest(method, path)
	err := GetBackend().Call(ctx, req, &testResource{})
	require.NoError(t, err)
}

// TestBackendCall_SuccessfulResponse_PostRequest tests that for POST
// requests (or other mutating operations) we serialize all parameters
// in the request body.
func TestBackendCall_SuccessfulResponse_PostRequest(t *testing.T) {
	ctx := context.Background()
	name := "the-name"
	rawJSON := `{"id":"res_123","object":"resource"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request parameters were passed correctly in
		// the request body.
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		assert.JSONEq(t, fmt.Sprintf(`{"name":"%s"}`, name), string(body))

		_, err = w.Write([]byte(rawJSON))
		require.NoError(t, err)
	}))
	defer ts.Close()

	// Set up a mock backend which triggers requests to our test server above.
	SetBackend(NewBackend(&BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	// Simulate usage for an API operation on a testResource.
	// We need to initialize a request and use the Backend to send it.
	resource := &testResource{}
	req := NewAPIRequest(http.MethodPost, "/resources")
	req.SetParams(&testResourceParams{Name: name})
	err := GetBackend().Call(ctx, req, resource)
	require.NoError(t, err)

	// The API response has been unmarshaled in the testResource struct.
	assert.Equal(t, "resource", resource.Object)
	assert.Equal(t, "res_123", resource.ID)
	// We stored the API response
	require.NotNil(t, resource.Response)
	assert.JSONEq(t, rawJSON, string(resource.Response.RawJSON))
}

// TestBackendCall_SuccessfulResponse_GetRequest tests that for GET
// requests which don't have a body, we serialize any parameters in
// the URL query string.
func TestBackendCall_SuccessfulResponse_GetRequest(t *testing.T) {
	ctx := context.Background()
	name := "the-name"
	limit := 1
	rawJSON := `{"data": [{"id":"res_123","object":"resource"}], "total_count": 1}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request parameters were set in the URL
		// query string.
		q := r.URL.Query()
		assert.Equal(t, name, q.Get("name"))
		assert.Equal(t, strconv.Itoa(limit), q.Get("limit"))
		// Optional query parameters are omitted.
		_, ok := q["offset"]
		assert.False(t, ok)
		// Existing query parameters are preserved
		assert.Equal(t, "still-here", q.Get("existing"))
		// Existing query parameters will be appended, not overriden
		assert.Equal(t, []string{"false", "true"}, q["appended"])

		_, err := w.Write([]byte(rawJSON))
		require.NoError(t, err)
	}))
	defer ts.Close()

	// Set up a mock backend which triggers requests to our test server above.
	SetBackend(NewBackend(&BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	// Simulate usage for an API operation on a testResourceList.
	// We need to initialize a request and use the Backend to send it.
	resource := &testResourceList{}
	req := NewAPIRequest(http.MethodGet, "/resources?existing=still-here&appended=false")
	req.SetParams(&testResourceListParams{
		ListParams: ListParams{
			Limit: Int64(int64(limit)),
		},
		Name:     name,
		Appended: "true",
	})
	err := GetBackend().Call(ctx, req, resource)
	require.NoError(t, err)

	// The API response has been unmarshaled correctly into a list of
	// testResource structs.
	assert.Equal(t, "resource", resource.Resources[0].Object)
	assert.Equal(t, "res_123", resource.Resources[0].ID)
	// We stored the API response
	require.NotNil(t, resource.Response)
	assert.JSONEq(t, rawJSON, string(resource.Response.RawJSON))
}

// TestBackendCall_ParseableError tests responses with a non-successful
// status code and a body that can be deserialized to an "expected"
// error response. These errors usually happen due to a client error
// and result in 4xx response statuses. The Clerk API responds with a
// familiar response body.
func TestBackendCall_ParseableError(t *testing.T) {
	errorJSON := `{
	"clerk_trace_id": "trace-id",
	"errors": [
		{
			"code": "error-code",
			"message": "error-message",
			"long_message": "long-error-message",
			"meta": {
				"param_name": "param-name"
			}
		}
	]
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err := w.Write([]byte(errorJSON))
		require.NoError(t, err)
	}))
	defer ts.Close()

	SetBackend(NewBackend(&BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))
	resource := &testResource{}
	err := GetBackend().Call(context.Background(), NewAPIRequest(http.MethodPost, "/resources"), resource)
	require.Error(t, err)

	// The error is an APIErrorResponse. We can assert on certain useful fields.
	apiErr, ok := err.(*APIErrorResponse)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.HTTPStatusCode)
	assert.Equal(t, "trace-id", apiErr.TraceID)

	// The response errors have been deserialized correctly.
	require.Equal(t, 1, len(apiErr.Errors))
	assert.Equal(t, "error-code", apiErr.Errors[0].Code)
	assert.Equal(t, "error-message", apiErr.Errors[0].Message)
	assert.Equal(t, "long-error-message", apiErr.Errors[0].LongMessage)
	assert.JSONEq(t, `{"param_name":"param-name"}`, string(apiErr.Errors[0].Meta))

	// We've stored the raw response as well.
	require.NotNil(t, apiErr.Response)
	assert.JSONEq(t, errorJSON, string(apiErr.Response.RawJSON))
}

// TestBackendCall_ParseableError tests responses with a non-successful
// status code and a body that can be deserialized to an unexpected
// error response. This might happen when the Clerk API encounters an
// unexpected server error and usually results in 5xx status codes.
func TestBackendCall_NonParseableError(t *testing.T) {
	errorResponse := `{invalid}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(errorResponse))
		require.NoError(t, err)
	}))
	defer ts.Close()

	SetBackend(NewBackend(&BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))
	resource := &testResource{}
	err := GetBackend().Call(context.Background(), NewAPIRequest(http.MethodPost, "/resources"), resource)
	require.Error(t, err)
	// The raw error is returned since we cannot unmarshal it to a
	// familiar API error response.
	assert.Equal(t, errorResponse, err.Error())
}
