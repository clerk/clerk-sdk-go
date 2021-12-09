package clerk

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWithSessionV2_nonBrowserRequest(t *testing.T) {
	c, mux, serverUrl, teardown := setup("test_dummy")
	defer teardown()

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should be signed out
		_, ok := r.Context().Value(ActiveSessionClaims).(*SessionClaims)
		if ok {
			t.Error("Expected session claims not to be present in request context")
		}
	})

	mux.Handle("/dummy", WithSessionV2(c)(dummyHandler))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/dummy", serverUrl), nil)
	req.Header.Set("User-Agent", "curl/7.64.1")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}
