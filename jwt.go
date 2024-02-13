package clerk

import (
	"context"
	"encoding/json"

	"github.com/go-jose/go-jose/v3/jwt"
)

type key string

const clerkActiveSessionClaims = key("clerkActiveSessionClaims")

// ContextWithSessionClaims returns a new context which includes the
// active session claims.
func ContextWithSessionClaims(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, clerkActiveSessionClaims, value)
}

// SessionClaimsFromContext returns the active session claims from
// the context.
func SessionClaimsFromContext(ctx context.Context) (*SessionClaims, bool) {
	claims, ok := ctx.Value(clerkActiveSessionClaims).(*SessionClaims)
	return claims, ok
}

// SessionClaims represents Clerk specific JWT claims.
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

// HasPermission checks if the session claims contain the provided
// organization permission.
// Use this helper to check if a user has the specific permission in
// the active organization.
func (s *SessionClaims) HasPermission(permission string) bool {
	for _, sessPermission := range s.ActiveOrganizationPermissions {
		if sessPermission == permission {
			return true
		}
	}
	return false
}

// HasRole checks if the session claims contain the provided
// organization role.
// However, the HasPermission helper is the recommended way to
// check for permissions. Complex role checks can usually be
// translated to a single permission check.
// For example, checks for an "admin" role that can modify a resource
// can be replaced by checks for a "modify" permission.
func (s *SessionClaims) HasRole(role string) bool {
	return s.ActiveOrganizationRole == role
}

// Claims holds generic JWT claims.
type Claims struct {
	jwt.Claims
	// Any headers not recognized get unmarshalled
	// from JSON in a generic manner and placed in this map.
	Extra map[string]any
}
