package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	clerkJSVersion := "42"
	err := client.Instances().Update(UpdateInstanceParams{
		TestMode:                    &enabled,
		HIBP:                        &enabled,
		EnhancedEmailDeliverability: &enabled,
		SupportEmail:                &supportEmail,
		ClerkJSVersion:              &clerkJSVersion,
	})

	if err != nil {
		t.Errorf("expected no error to be returned, found %v instead", err)
	}
}

func TestInstanceService_Update_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	enabled := true
	supportEmail := "support@clerk.dev"
	clerkJSVersion := "999"
	err := client.Instances().Update(UpdateInstanceParams{
		TestMode:                    &enabled,
		HIBP:                        &enabled,
		EnhancedEmailDeliverability: &enabled,
		SupportEmail:                &supportEmail,
		ClerkJSVersion:              &clerkJSVersion,
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestInstanceService_UpdateRestrictions_happyPath(t *testing.T) {
	token := "token"
	dummyRestrictionsResponseJSON := `{
		"allowlist": true,
		"blocklist": true,
		"block_email_subaddresses": true,
		"block_disposable_email_domains": true
	}`
	var restrictionsResponse InstanceRestrictionsResponse
	err := json.Unmarshal([]byte(dummyRestrictionsResponseJSON), &restrictionsResponse)
	assert.NoError(t, err)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/instance/restrictions", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPatch)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyRestrictionsResponseJSON)
	})

	enabled := true
	got, _ := client.Instances().UpdateRestrictions(UpdateRestrictionsParams{
		Allowlist:                   &enabled,
		Blocklist:                   &enabled,
		BlockEmailSubaddresses:      &enabled,
		BlockDisposableEmailDomains: &enabled,
		IgnoreDotsForGmailAddresses: &enabled,
	})

	assert.Equal(t, &restrictionsResponse, got)
}

func TestInstanceService_UpdateRestrictions_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	enabled := true
	_, err := client.Instances().UpdateRestrictions(UpdateRestrictionsParams{
		Allowlist: &enabled,
		Blocklist: &enabled,
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestInstanceService_UpdateOrganizationSettings_happyPath(t *testing.T) {
	token := "token"
	dummyOrganizationSettingsResponseJSON := `{
		"enabled": true,
		"max_allowed_memberships": 2,
		"max_allowed_roles": 10,
		"max_allowed_permissions": 50,
		"creator_role": "org:custom_admin",
		"admin_delete_enabled": true,
		"domains_enabled": true,
		"domains_enrollment_modes": [
			"manual_invitation",
			"automatic_invitation",
			"automatic_suggestion"
		],
		"domains_default_role": "org:custom_domains"
	}`
	var organizationSettingsResponse OrganizationSettingsResponse
	_ = json.Unmarshal([]byte(dummyOrganizationSettingsResponseJSON), &organizationSettingsResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/instance/organization_settings", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPatch)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyOrganizationSettingsResponseJSON)
	})

	enabled := true
	got, _ := client.Instances().UpdateOrganizationSettings(UpdateOrganizationSettingsParams{
		Enabled:              &enabled,
		CreatorRoleID:        stringToPtr("role_2XcSZn6swGCjX59Nk0XbGer22jb"),
		DomainsDefaultRoleID: stringToPtr("role_2XZCQwxfLbXOz2hoBXKFVRjwmGc"),
	})

	assert.Equal(t, &organizationSettingsResponse, got)
}

func TestInstanceService_UpdateOrganizationSettings_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	enabled := true
	_, err := client.Instances().UpdateOrganizationSettings(UpdateOrganizationSettingsParams{
		Enabled: &enabled,
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}
