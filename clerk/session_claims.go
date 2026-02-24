package clerk

import (
	"encoding/json"
	"strings"

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
	Features                      string          `json:"fea,omitempty"`
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

// HasFeature checks if the user has access to the specified feature
// in their session claims.
// The feature can be unscoped ("dark_mode") to match any scope,
// or scoped ("org:dark_mode", "user:dark_mode") to match a specific scope.
// An unrecognized scope prefix returns false.
func (s *SessionClaims) HasFeature(feature string) bool {
	if s.Features == "" {
		return false
	}

	orgFeatures, userFeatures := splitFeaturesByScope(s.Features)

	colonIdx := strings.IndexByte(feature, ':')
	if colonIdx == -1 {
		// Unscoped: check both org and user features
		for _, f := range orgFeatures {
			if f == feature {
				return true
			}
		}
		for _, f := range userFeatures {
			if f == feature {
				return true
			}
		}
		return false
	}

	scope := feature[:colonIdx]
	id := feature[colonIdx+1:]

	switch scope {
	case "o", "org", "organization":
		for _, f := range orgFeatures {
			if f == id {
				return true
			}
		}
	case "u", "user":
		for _, f := range userFeatures {
			if f == id {
				return true
			}
		}
	}
	return false
}

// splitFeaturesByScope parses the comma-separated fea JWT claim
// into org-scoped and user-scoped feature slices.
// Claim format: "o:dark_mode,u:new_billing,uo:feature_flags"
// Malformed entries (missing colon) are silently skipped.
func splitFeaturesByScope(fea string) (org []string, user []string) {
	parts := strings.Split(fea, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		colonIdx := strings.IndexByte(part, ':')
		if colonIdx == -1 {
			continue
		}
		scope := part[:colonIdx]
		value := part[colonIdx+1:]

		switch scope {
		case "o":
			org = append(org, value)
		case "u":
			user = append(user, value)
		case "ou", "uo":
			org = append(org, value)
			user = append(user, value)
		}
	}
	return
}
