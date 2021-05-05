// +build integration

package integration

import (
	"errors"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"testing"
)

func TestWebhooks(t *testing.T) {
	client := createClient()

	diahookResponse, err := client.Webhooks().CreateDiahook()
	if err != nil {
		checkIfCanRecoverFromDiahookCreation(t, err)
	} else if diahookResponse.DiahookURL == "" {
		t.Fatalf("Webhooks.CreateDiahook returned empty url")
	}

	diahookResponse, err = client.Webhooks().RefreshDiahookURL()
	if err != nil {
		t.Fatalf("Webhooks.RefreshDiahookURL returned unexpected unexpected error %v", err)
	} else if diahookResponse.DiahookURL == "" {
		t.Fatalf("was expecting a Diahook url, found none instead")
	}

	err = client.Webhooks().DeleteDiahook()
	if err != nil {
		t.Fatalf("Webhooks.DeleteDiahook returned unexpected error %v", err)
	}
}

// checkIfCanRecoverFromDiahookCreation checks whether the error from creating
// a new Diahook app was that there was already a Diahook for the given instance.
// If it was, then we can continue, otherwise we fail.
func checkIfCanRecoverFromDiahookCreation(t *testing.T, err error) {
	var errorResponse *clerk.ErrorResponse
	if !errors.As(err, &errorResponse) {
		t.Fatalf("unexpected error found while creating diahook: %v", err)
	}
	if len(errorResponse.Errors) != 1 {
		t.Fatalf("was only expecting at most one error, found %d instead: %v", len(errorResponse.Errors), errorResponse.Errors)
	}
	if errorResponse.Errors[0].Code != "diahook_app_exists" {
		t.Fatalf("was expecting a diahook_app_exists error, found %s instead", errorResponse.Errors[0].Code)
	}
}
