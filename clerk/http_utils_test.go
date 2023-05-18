package clerk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setup sets up a test HTTP server along with a `clerk.Client` that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup(token string) (client Client, mux *http.ServeMux, serverURL *url.URL, teardown func()) {
	versionPath := "/v1"

	mux = http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(versionPath+"/", http.StripPrefix(versionPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	baseURL, _ := url.Parse(server.URL + versionPath + "/")
	client, _ = NewClient(token, WithBaseURL(baseURL.String()))

	return client, mux, baseURL, server.Close
}

func testHttpMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testQuery(t *testing.T, r *http.Request, want url.Values) {
	t.Helper()

	query := r.URL.Query()

	for k := range want {
		if query.Get(k) == "" {
			t.Errorf("Request query doesn't match: have %v, want %v", query, want)
		}
	}
}
