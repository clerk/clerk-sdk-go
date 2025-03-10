package clerk

import (
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
