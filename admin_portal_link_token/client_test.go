package admin_portal_link_token

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate_minimal(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"object":               "admin_portal_link_token",
		"id":                   "aplt_3beecc9c60adb5f9b850e91a83beecc9c",
		"admin_portal_link_id": "apl_3beecc9c60adb5f9b850e91a8d2",
		"instance_id":          "ins_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
		"organization_id":      nil,
		"it_contact_id":        nil,
		"scopes":               nil,
		"token":                "aplt_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		"revoked":              false,
		"revocation_reason":    nil,
		"expired":              false,
		"expiration":           1716883200000,
		"created_at":           1640995200000,
		"updated_at":           1640995200000,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/admin_portal_link_tokens",
		},
	}

	client := NewClient(config)
	token, err := client.Create(context.Background(), &CreateParams{})
	require.NoError(t, err)
	assert.Equal(t, "admin_portal_link_token", token.Object)
	assert.Equal(t, "aplt_3beecc9c60adb5f9b850e91a83beecc9c", token.ID)
	assert.Equal(t, "apl_3beecc9c60adb5f9b850e91a8d2", token.AdminPortalLinkID)
	assert.Equal(t, "ins_2xhFjEI5X2qWRvtV13BzSj8H6Dk", token.InstanceID)
	assert.Nil(t, token.OrganizationID)
	assert.Nil(t, token.ITContactID)
	assert.Nil(t, token.Scopes)
	assert.Equal(t, "aplt_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", token.Token)
	assert.False(t, token.Revoked)
}

func TestCreate_withOrgAndContact(t *testing.T) {
	t.Parallel()

	orgID := "org_2xhFjEI5X2qWRvtV13BzSj8H6Dk"
	itContactID := "usr_abc123"

	response := map[string]interface{}{
		"object":               "admin_portal_link_token",
		"id":                   "aplt_3beecc9c60adb5f9b850e91a83beecc9c",
		"admin_portal_link_id": "apl_3beecc9c60adb5f9b850e91a8d2",
		"instance_id":          "ins_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
		"organization_id":      orgID,
		"it_contact_id":        itContactID,
		"scopes":               []string{"admin_portal:read"},
		"token":                "aplt_YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
		"revoked":              false,
		"revocation_reason":    nil,
		"expired":              false,
		"expiration":           1716883200000,
		"created_at":           1640995200000,
		"updated_at":           1640995200000,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"organization_id":"org_2xhFjEI5X2qWRvtV13BzSj8H6Dk","it_contact_id":"usr_abc123","scopes":["admin_portal:read"],"seconds_until_expiration":600}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/admin_portal_link_tokens",
		},
	}

	client := NewClient(config)
	params := &CreateParams{
		OrganizationID:         clerk.String(orgID),
		ITContactID:            clerk.String(itContactID),
		Scopes:                 []string{"admin_portal:read"},
		SecondsUntilExpiration: clerk.Int64(600),
	}

	token, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, token.OrganizationID)
	assert.Equal(t, orgID, *token.OrganizationID)
	require.NotNil(t, token.ITContactID)
	assert.Equal(t, itContactID, *token.ITContactID)
	assert.Equal(t, []string{"admin_portal:read"}, token.Scopes)
}

func TestRevoke(t *testing.T) {
	t.Parallel()

	tokenID := "aplt_3beecc9c60adb5f9b850e91a83beecc9c"
	reason := "Revoked by user"

	response := map[string]interface{}{
		"object":               "admin_portal_link_token",
		"id":                   tokenID,
		"admin_portal_link_id": "apl_3beecc9c60adb5f9b850e91a8d2",
		"instance_id":          "ins_2xhFjEI5X2qWRvtV13BzSj8H6Dk",
		"organization_id":      nil,
		"it_contact_id":        nil,
		"scopes":               nil,
		"revoked":              true,
		"revocation_reason":    reason,
		"expired":              false,
		"expiration":           1716883200000,
		"created_at":           1640995200000,
		"updated_at":           1640995200000,
	}

	responseJSON, _ := json.Marshal(response)

	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"revocation_reason":"Revoked by user"}`),
			Out:    json.RawMessage(responseJSON),
			Method: http.MethodPost,
			Path:   "/v1/admin_portal_link_tokens/" + tokenID + "/revoke",
		},
	}

	client := NewClient(config)
	params := &RevokeParams{
		RevocationReason: clerk.String(reason),
	}

	token, err := client.Revoke(context.Background(), tokenID, params)
	require.NoError(t, err)
	assert.Equal(t, tokenID, token.ID)
	assert.True(t, token.Revoked)
	require.NotNil(t, token.RevocationReason)
	assert.Equal(t, reason, *token.RevocationReason)
}
