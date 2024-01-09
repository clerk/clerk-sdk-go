package clerk

import "slices"

type CheckAuthorizationParams struct {
	Permission string
	Role       string
}

// CheckAuthorization verifies if the user has the given permission or role.
// Performing role checks is not considered a best-practice and
// developers should avoid it as much as possible.
// Usually, complex role checks can be refactored with a single permission check.
func (c *client) CheckAuthorization(token string, params CheckAuthorizationParams) (bool, error) {
	claims, err := c.VerifyToken(token)
	if err != nil {
		return false, err
	}

	permission := params.Permission
	role := params.Role

	if permission != "" && slices.Contains(claims.ActiveOrganizationPermissions, permission) {
		return true, nil
	}

	if claims.ActiveOrganizationRole == role {
		return true, nil
	}

	return false, nil
}

// CheckPermission checks if the user has the specific permission
// in their session claims.
func (s *SessionClaims) CheckPermission(permission string) bool {
	return slices.Contains(s.ActiveOrganizationPermissions, permission)
}

// CheckRole checks if the user has the specific role
// in their session claims.
// Performing role checks is not considered a best-practice and
// developers should avoid it as much as possible.
// Usually, complex role checks can be refactored with a single permission check.
func (s *SessionClaims) CheckRole(role string) bool {
	return s.ActiveOrganizationRole == role
}
