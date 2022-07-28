//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func TestInstances(t *testing.T) {
	client := createClient()

	enabled := true
	supportEmail := "support@example.com"
	err := client.Instances().Update(clerk.UpdateInstanceParams{
		TestMode:                    &enabled,
		HIBP:                        &enabled,
		EnhancedEmailDeliverability: &enabled,
		SupportEmail:                &supportEmail,
	})
	if err != nil {
		t.Fatalf("Instances.Update returned error: %v", err)
	}
}
