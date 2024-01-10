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
