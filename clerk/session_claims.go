package clerk

import (
	"encoding/json"

	"github.com/go-jose/go-jose/v3/jwt"
)

type SessionClaims struct {
	jwt.Claims
	SessionID                     string          `json:"sid"`
	AuthorizedParty               string          `json:"azp"`
	ActiveOrganizationID          string          `json:"org_id"`
	ActiveOrganizationSlug        string          `json:"org_slug"`
	ActiveOrganizationRole        string          `json:"org_role"`
	ActiveOrganizationPermissions []string        `json:"org_permissions"`
	Actor                         json.RawMessage `json:"act,omitempty"`
}

// HasPermission checks if the user has the specific permission
// in their session claims.
func (s *SessionClaims) HasPermission(permission string) bool {
	for _, sessPermission := range s.ActiveOrganizationPermissions {
		if sessPermission == permission {
			return true
		}
	}
	return false
}

// HasRole checks if the user has the specific role
// in their session claims.
// Performing role checks is not considered a best-practice and
// developers should avoid it as much as possible.
// Usually, complex role checks can be refactored with a single permission check.
func (s *SessionClaims) HasRole(role string) bool {
	return s.ActiveOrganizationRole == role
}
