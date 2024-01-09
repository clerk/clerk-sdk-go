package clerk

import (
	"testing"

	"github.com/go-jose/go-jose/v3"
)

func TestClient_CheckAuthorization(t *testing.T) {
	c, _ := NewClient("token")
	token, pubKey := testGenerateTokenJWT(t, dummySessionClaims, "kid")
	client := c.(*client)
	client.jwksCache.set(testBuildJWKS(t, pubKey, jose.RS256, "kid"))

	// user has permission
	hasPermission, err := c.CheckAuthorization(
		token,
		CheckAuthorizationParams{Permission: "org:billing:manage"},
	)
	if err != nil {
		t.Error(err)
	}
	if !hasPermission {
		t.Errorf("Expected user to have permission: %s", "org:billing:manage")
	}

	// user does not have permission
	hasPermission, err = c.CheckAuthorization(
		token,
		CheckAuthorizationParams{Permission: "org:billing:create"},
	)
	if err != nil {
		t.Error(err)
	}
	if hasPermission {
		t.Errorf("Expected user to not have permission: %s", "org:billing:create")
	}

	// user has role
	hasRole, err := c.CheckAuthorization(
		token,
		CheckAuthorizationParams{Role: "org_role"},
	)
	if err != nil {
		t.Error(err)
	}
	if !hasRole {
		t.Errorf("Expected user to have role: %s", "org_role")
	}

	// user does not have role
	hasRole, err = c.CheckAuthorization(
		token,
		CheckAuthorizationParams{Role: "org_role_nonexistent"},
	)
	if err != nil {
		t.Error(err)
	}
	if hasRole {
		t.Errorf("Expected user to not have role: %s", "org_role_nonexistent")
	}
}

func TestSessionClaims_CheckPermission(t *testing.T) {
	// user has permission
	hasPermission := dummySessionClaims.CheckPermission("org:billing:manage")
	if !hasPermission {
		t.Errorf("Expected user to have permission: %s", "org:billing:manage")
	}

	// user does not have permission
	hasPermission = dummySessionClaims.CheckPermission("org:billing:create")
	if hasPermission {
		t.Errorf("Expected user to not have permission: %s", "org:billing:create")
	}
}

func TestSessionClaims_CheckRole(t *testing.T) {
	// user has role
	hasRole := dummySessionClaims.CheckRole("org_role")
	if !hasRole {
		t.Errorf("Expected user to have role: %s", "org_role")
	}

	// user does not have role
	hasRole = dummySessionClaims.CheckPermission("org_role_nonexistent")
	if hasRole {
		t.Errorf("Expected user to not have role: %s", "org_role_nonexistent")
	}
}
