package clerk

import (
	"context"
	"encoding/json"
	"time"

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
	// Standard IANA JWT claims
	RegisteredClaims
	// Clerk specific JWT claims
	Claims

	// Custom can hold any custom claims that might be found in a JWT.
	Custom any `json:"-"`
}

func (s *SessionClaims) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.RegisteredClaims)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &s.Claims)
	return err
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

// NeedsReverification checks if a session needs to be reverified based on the
// provided policy.
func (s *SessionClaims) NeedsReverification(policy SessionReverificationPolicy) bool {
	firstFactorAgeMinutes, secondFactorAgeMinutes := s.FactorVerificationAge[0], s.FactorVerificationAge[1]

	firstFactorNeedsReverification := firstFactorAgeMinutes == factorVerificationAgeNotEnabled ||
		policy.AfterMinutes < firstFactorAgeMinutes
	isSecondFactorEnabled := secondFactorAgeMinutes != factorVerificationAgeNotEnabled
	secondFactorNeedsReverification := isSecondFactorEnabled &&
		policy.AfterMinutes < secondFactorAgeMinutes

	switch policy.Level {
	case SessionReverificationLevelFirstFactor:
		return firstFactorNeedsReverification
	case SessionReverificationLevelSecondFactor:
		if !isSecondFactorEnabled {
			return firstFactorNeedsReverification
		}
		return secondFactorNeedsReverification
	case SessionReverificationLevelMultiFactor:
		if !isSecondFactorEnabled {
			return firstFactorNeedsReverification
		}
		return firstFactorNeedsReverification || secondFactorNeedsReverification
	default:
		return true
	}
}

// RegisteredClaims holds public claim values (as specified in RFC 7519).
type RegisteredClaims struct {
	Issuer    string   `json:"iss,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	Expiry    *int64   `json:"exp,omitempty"`
	NotBefore *int64   `json:"nbf,omitempty"`
	IssuedAt  *int64   `json:"iat,omitempty"`
	ID        string   `json:"jti,omitempty"`
	raw       jwt.Claims
}

func (c *RegisteredClaims) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &c.raw)
	if err != nil {
		return err
	}
	c.Issuer = c.raw.Issuer
	c.Subject = c.raw.Subject
	c.Audience = c.raw.Audience
	c.ID = c.raw.ID
	if c.raw.Expiry != nil {
		c.Expiry = Int64(c.raw.Expiry.Time().Unix())
	}
	if c.raw.NotBefore != nil {
		c.NotBefore = Int64(c.raw.NotBefore.Time().Unix())
	}
	if c.raw.IssuedAt != nil {
		c.IssuedAt = Int64(c.raw.IssuedAt.Time().Unix())
	}
	return nil
}

// ValidateWithLeeway checks expiration and issuance claims against
// an expected time.
// You may pass a zero value to check the time values with no leeway,
// but it is not recommended.
// The leeway gives some extra time to the token from the server's
// point of view. That is, if the token is expired, ValidateWithLeeway
// will still accept the token for 'leeway' amount of time.
func (c *RegisteredClaims) ValidateWithLeeway(expected time.Time, leeway time.Duration) error {
	return c.raw.ValidateWithLeeway(jwt.Expected{Time: expected}, leeway)
}

// Claims represents private JWT claims that are defined and used
// by Clerk.
type Claims struct {
	Version                       int             `json:"v"`
	SessionID                     string          `json:"sid"`
	AuthorizedParty               string          `json:"azp"`
	ActiveOrganizationID          string          `json:"org_id"`
	ActiveOrganizationSlug        string          `json:"org_slug"`
	ActiveOrganizationRole        string          `json:"org_role"`
	ActiveOrganizationPermissions []string        `json:"org_permissions"`
	Actor                         json.RawMessage `json:"act,omitempty"`
	FactorVerificationAge         [2]int64        `json:"fva"`
}

// factorVerificationAgeNotEnabled is the value Clerk uses in either 'fva' slot
// for a factor group that isn't enabled.
const factorVerificationAgeNotEnabled int64 = -1

// UnmarshalJSON decodes Clerk's private claims, handling omitted 'fva'.
//
// This is necessary because an omitted 'fva' otherwise marshals to [0, 0]. The
// zero value here would satisfy every reverification check, which is not intended.
// Plenty of Clerk-signed JWTs carry no 'fva' claim at all (JWT template
// tokens, OIDC ID tokens and machine tokens, none of which describe a session
// step-up), but this should not entail that we just completed reverification.
func (c *Claims) UnmarshalJSON(data []byte) error {
	// The alias has the same fields but no methods, so this doesn't recurse.
	type alias Claims
	if err := json.Unmarshal(data, (*alias)(c)); err != nil {
		return err
	}

	// Decode 'fva' again, into a slice this time. By now c.FactorVerificationAge
	// has been zero-filled, and a [2]int64 can't tell us whether that came from
	// the token or from Go. A slice can: nil means the claim wasn't there.
	var raw struct {
		FactorVerificationAge []int64 `json:"fva"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.FactorVerificationAge = [2]int64{factorVerificationAgeNotEnabled, factorVerificationAgeNotEnabled}
	if len(raw.FactorVerificationAge) == 2 &&
		isValidFactorVerificationAge(raw.FactorVerificationAge[0]) &&
		isValidFactorVerificationAge(raw.FactorVerificationAge[1]) {
		c.FactorVerificationAge = [2]int64(raw.FactorVerificationAge)
	}
	return nil
}

// isValidFactorVerificationAge reports whether age is a usable 'fva' entry:
// either an age in minutes, or the sentinel for a factor that is not set.
func isValidFactorVerificationAge(age int64) bool {
	return age == factorVerificationAgeNotEnabled || age >= 0
}

// UnverifiedToken holds the result of a JWT decoding without any
// verification.
// The UnverifiedToken includes registered and custom claims, as
// well as the KeyID (kid) header.
type UnverifiedToken struct {
	RegisteredClaims
	// Any headers not recognized get unmarshalled
	// from JSON in a generic manner and placed in this map.
	Extra map[string]any
	KeyID string
}
