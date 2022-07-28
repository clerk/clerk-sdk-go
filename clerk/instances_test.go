package clerk

import (
	"net/http"
	"testing"
)

func TestInstanceService_Update_happyPath(t *testing.T) {
	token := "token"
	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/instance", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPatch)
		testHeader(t, req, "Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusNoContent)
	})

	enabled := true
	supportEmail := "support@clerk.dev"
	err := client.Instances().Update(UpdateInstanceParams{
		TestMode:     &enabled,
		HIBP:         &enabled,
		SupportEmail: &supportEmail,
	})

	if err != nil {
		t.Errorf("expected no error to be return, found %v instead", err)
	}
}

func TestInstanceService_Update_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	enabled := true
	supportEmail := "support@clerk.dev"
	err := client.Instances().Update(UpdateInstanceParams{
		TestMode:     &enabled,
		HIBP:         &enabled,
		SupportEmail: &supportEmail,
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}
