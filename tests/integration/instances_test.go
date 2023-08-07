//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/stretchr/testify/assert"
)

func TestInstances(t *testing.T) {
	client := createClient()

	enabled := true
	developmentOrigin := "http://localhost:3000"
	supportEmail := "support@example.com"
	err := client.Instances().Update(clerk.UpdateInstanceParams{
		TestMode:                    &enabled,
		HIBP:                        &enabled,
		EnhancedEmailDeliverability: &enabled,
		SupportEmail:                &supportEmail,
		DevelopmentOrigin:           &developmentOrigin,
	})
	if err != nil {
		t.Fatalf("Instances.Update returned error: %v", err)
	}
}

func TestInstanceRestrictions(t *testing.T) {
	client := createClient()

	enabled := true
	restrictionsResponse, err := client.Instances().UpdateRestrictions(clerk.UpdateRestrictionsParams{
		Allowlist: &enabled,
		Blocklist: &enabled,
	})
	assert.Nil(t, err)
	assert.True(t, restrictionsResponse.Allowlist)
	assert.True(t, restrictionsResponse.Blocklist)
	assert.False(t, restrictionsResponse.BlockEmailSubaddresses)
}

func TestInstanceOrganizationSettings(t *testing.T) {
	client := createClient()

	enabled := true
	allowedMemberships := 0
	organizationSettingsResponse, err := client.Instances().UpdateOrganizationSettings(clerk.UpdateOrganizationSettingsParams{
		Enabled:               &enabled,
		MaxAllowedMemberships: &allowedMemberships,
	})
	assert.Nil(t, err)
	assert.True(t, organizationSettingsResponse.Enabled)
}
