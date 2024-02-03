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
		claims := SessionClaims{
			ActiveOrganizationRole: tc.active,
		}
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
		claims := SessionClaims{
			ActiveOrganizationPermissions: tc.active,
		}
		require.Equal(t, claims.HasPermission(tc.permission), tc.want)
	}
}
