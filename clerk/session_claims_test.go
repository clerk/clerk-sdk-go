package clerk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionClaims_HasPermissiont(t *testing.T) {
	// user has permission
	hasPermission := dummySessionClaims.HasPermission("org:billing:manage")
	assert.True(t, hasPermission)

	// user has second permission
	hasPermission = dummySessionClaims.HasPermission("org:report:view")
	assert.True(t, hasPermission)

	// user does not have permission
	hasPermission = dummySessionClaims.HasPermission("org:billing:create")
	assert.False(t, hasPermission)
}

func TestSessionClaims_HasRole(t *testing.T) {
	// user has role
	hasRole := dummySessionClaims.HasRole("org_role")
	assert.True(t, hasRole)

	// user does not have role
	hasRole = dummySessionClaims.HasRole("org_role_nonexistent")
	assert.False(t, hasRole)
}

func TestSessionClaims_HasFeature(t *testing.T) {
	claims := SessionClaims{
		Features: "o:dark_mode,u:new_billing,uo:feature_flags",
	}

	// unscoped matches any scope
	assert.True(t, claims.HasFeature("dark_mode"))
	assert.True(t, claims.HasFeature("new_billing"))
	assert.True(t, claims.HasFeature("feature_flags"))

	// org-scoped queries (all aliases)
	assert.True(t, claims.HasFeature("o:dark_mode"))
	assert.True(t, claims.HasFeature("org:dark_mode"))
	assert.True(t, claims.HasFeature("organization:dark_mode"))

	// user-scoped queries (all aliases)
	assert.True(t, claims.HasFeature("u:new_billing"))
	assert.True(t, claims.HasFeature("user:new_billing"))

	// uo/ou features match both scopes
	assert.True(t, claims.HasFeature("org:feature_flags"))
	assert.True(t, claims.HasFeature("user:feature_flags"))

	// org-only feature does NOT match user scope
	assert.False(t, claims.HasFeature("user:dark_mode"))

	// user-only feature does NOT match org scope
	assert.False(t, claims.HasFeature("org:new_billing"))

	// non-existent feature
	assert.False(t, claims.HasFeature("nonexistent"))

	// invalid scope returns false
	assert.False(t, claims.HasFeature("bad:dark_mode"))

	// empty features claim
	emptyClaims := SessionClaims{}
	assert.False(t, emptyClaims.HasFeature("anything"))

	// malformed entry is skipped, valid entries still work
	malformedClaims := SessionClaims{
		Features: "broken,o:valid_feature",
	}
	assert.True(t, malformedClaims.HasFeature("valid_feature"))
	assert.False(t, malformedClaims.HasFeature("broken"))
}
