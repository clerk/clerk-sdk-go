package clerk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionClaimsHasRole(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		active string
		role   string
		want   bool
	}{
		{
			active: "active",
			role:   "non-active",
			want:   false,
		},
		{
			active: "active",
			role:   "active",
			want:   true,
		},
	} {
		claims := SessionClaims{}
		claims.ActiveOrganizationRole = tc.active
		require.Equal(t, claims.HasRole(tc.role), tc.want)
	}
}

func TestSessionClaimsHasPermission(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		active     []string
		permission string
		want       bool
	}{
		{
			active:     []string{"active"},
			permission: "non-active",
			want:       false,
		},
		{
			active:     []string{"active", "non-active"},
			permission: "active",
			want:       true,
		},
		{
			active:     []string{},
			permission: "active",
			want:       false,
		},
	} {
		claims := SessionClaims{}
		claims.ActiveOrganizationPermissions = tc.active
		require.Equal(t, claims.HasPermission(tc.permission), tc.want)
	}
}

// TestSessionClaimsUnmarshalJSONFactorVerificationAge tests that a token which
// carries no usable 'fva' claim decodes to the "not set" sentinel instead of
// Go's zero value, which NeedsReverification would read as "verified just now".
func TestSessionClaimsUnmarshalJSONFactorVerificationAge(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload string
		want    [2]int64
	}{
		{
			name:    "claim present",
			payload: `{"sid":"sess_123","fva":[10,5]}`,
			want:    [2]int64{10, 5},
		},
		{
			name:    "claim present, second factor not enrolled",
			payload: `{"sid":"sess_123","fva":[10,-1]}`,
			want:    [2]int64{10, -1},
		},
		{
			name:    "claim present, both factors just verified",
			payload: `{"sid":"sess_123","fva":[0,0]}`,
			want:    [2]int64{0, 0},
		},
		{
			// JWT template tokens, OIDC ID tokens and machine tokens.
			name:    "claim absent",
			payload: `{"sid":"sess_123"}`,
			want:    [2]int64{-1, -1},
		},
		{
			name:    "claim null",
			payload: `{"sid":"sess_123","fva":null}`,
			want:    [2]int64{-1, -1},
		},
		{
			name:    "claim empty",
			payload: `{"sid":"sess_123","fva":[]}`,
			want:    [2]int64{-1, -1},
		},
		{
			// Go would zero-fill the missing element to [5, 0], marking the
			// second factor enrolled and freshly verified.
			name:    "claim too short",
			payload: `{"sid":"sess_123","fva":[5]}`,
			want:    [2]int64{-1, -1},
		},
		{
			name:    "claim too long",
			payload: `{"sid":"sess_123","fva":[5,10,15]}`,
			want:    [2]int64{-1, -1},
		},
		{
			// Negative ages other than the sentinel compare as fresher than
			// any policy threshold.
			name:    "claim with out of range age",
			payload: `{"sid":"sess_123","fva":[-5,-5]}`,
			want:    [2]int64{-1, -1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			claims := &SessionClaims{}
			require.NoError(t, json.Unmarshal([]byte(tc.payload), claims))
			require.Equal(t, tc.want, claims.FactorVerificationAge)
			// The rest of the claims still decode.
			require.Equal(t, "sess_123", claims.SessionID)
		})
	}
}

// TestNeedsReverificationWithoutFactorVerificationAge tests that every policy
// fails closed for a token that carries no 'fva' claim.
func TestNeedsReverificationWithoutFactorVerificationAge(t *testing.T) {
	t.Parallel()
	for name, policy := range map[string]SessionReverificationPolicy{
		"strict mfa":   SessionReverificationStrictMFA,
		"strict":       SessionReverificationStrict,
		"moderate":     SessionReverificationModerate,
		"lax":          SessionReverificationLax,
		"first factor": {AfterMinutes: 10, Level: SessionReverificationLevelFirstFactor},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			claims := &SessionClaims{}
			require.NoError(t, json.Unmarshal([]byte(`{"sid":"sess_123"}`), claims))
			require.True(t, claims.NeedsReverification(policy))
		})
	}
}

func TestNeedsReverification(t *testing.T) {
	for _, tc := range []struct {
		name           string
		factorAges     [2]int64
		policyLevel    SessionReverificationLevel
		policyMinutes  int64
		expectedResult bool
	}{
		// FirstFactor
		{
			name:           "first factor - both valid",
			factorAges:     [2]int64{5, -1},
			policyLevel:    SessionReverificationLevelFirstFactor,
			policyMinutes:  10,
			expectedResult: false,
		},
		{
			name:           "first factor - exactly at the threshold",
			factorAges:     [2]int64{10, -1},
			policyLevel:    SessionReverificationLevelFirstFactor,
			policyMinutes:  10,
			expectedResult: false,
		},
		{
			name:           "first factor needs reverification",
			factorAges:     [2]int64{15, -1},
			policyLevel:    SessionReverificationLevelFirstFactor,
			policyMinutes:  10,
			expectedResult: true,
		},
		{
			name:           "first factor - invalid state",
			factorAges:     [2]int64{-1, -1},
			policyLevel:    SessionReverificationLevelFirstFactor,
			policyMinutes:  10,
			expectedResult: true,
		},

		// SecondFactor
		{
			name:           "second factor both valid",
			factorAges:     [2]int64{5, 5},
			policyLevel:    SessionReverificationLevelSecondFactor,
			policyMinutes:  10,
			expectedResult: false,
		},
		{
			name:           "second factor needs reverification",
			factorAges:     [2]int64{5, 15},
			policyLevel:    SessionReverificationLevelSecondFactor,
			policyMinutes:  10,
			expectedResult: true,
		},
		{
			name:           "second factor valid, but first not",
			factorAges:     [2]int64{15, 0},
			policyLevel:    SessionReverificationLevelSecondFactor,
			policyMinutes:  10,
			expectedResult: false,
		},
		{
			name:           "second factor not enabled",
			factorAges:     [2]int64{15, -1},
			policyLevel:    SessionReverificationLevelSecondFactor,
			policyMinutes:  10,
			expectedResult: true,
		},
		{
			name:           "second factor - invalid state",
			factorAges:     [2]int64{-1, -1},
			policyLevel:    SessionReverificationLevelSecondFactor,
			policyMinutes:  10,
			expectedResult: true,
		},

		// MultiFactor
		{
			name:           "multi factor - both valid",
			factorAges:     [2]int64{5, 5},
			policyLevel:    SessionReverificationLevelMultiFactor,
			policyMinutes:  10,
			expectedResult: false,
		},
		{
			name:           "multi factor - first needs reverification, second valid",
			factorAges:     [2]int64{15, 5},
			policyLevel:    SessionReverificationLevelMultiFactor,
			policyMinutes:  10,
			expectedResult: true,
		},
		{
			name:           "multi factor - first valid, second needs reverification",
			factorAges:     [2]int64{5, 15},
			policyLevel:    SessionReverificationLevelMultiFactor,
			policyMinutes:  10,
			expectedResult: true,
		},
		{
			name:           "multi factor - both need reverification",
			factorAges:     [2]int64{15, 15},
			policyLevel:    SessionReverificationLevelMultiFactor,
			policyMinutes:  10,
			expectedResult: true,
		},
		{
			name:           "multi factor - second not enabled",
			factorAges:     [2]int64{5, -1},
			policyLevel:    SessionReverificationLevelMultiFactor,
			policyMinutes:  10,
			expectedResult: false,
		},
		{
			name:           "multi factor - invalid state",
			factorAges:     [2]int64{-1, -1},
			policyLevel:    SessionReverificationLevelMultiFactor,
			policyMinutes:  10,
			expectedResult: true,
		},

		// Edge cases
		{
			name:           "invalid policy level",
			factorAges:     [2]int64{5, 5},
			policyLevel:    SessionReverificationLevel("InvalidLevel"),
			policyMinutes:  10,
			expectedResult: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sessionClaims := &SessionClaims{
				Claims: Claims{
					FactorVerificationAge: tc.factorAges,
				},
			}
			policy := SessionReverificationPolicy{
				Level:        tc.policyLevel,
				AfterMinutes: tc.policyMinutes,
			}

			result := sessionClaims.NeedsReverification(policy)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
