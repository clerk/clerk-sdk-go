package clerk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setup sets up a test HTTP server along with a `clerk.Client` that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup(token string) (client *client, mux *http.ServeMux, serverURL string, teardown func()) {
	versionPath := "/v1"

	mux = http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(versionPath+"/", http.StripPrefix(versionPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	url, _ := url.Parse(server.URL + versionPath + "/")
	client, _ = NewClientWithBaseUrl(token, url.String())

	return client, mux, server.URL, server.Close
}

func testHttpMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}
