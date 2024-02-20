package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/require"
)

func TestVerify_InvalidParams(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"
	token, pubKey := clerktest.GenerateJWT(t, map[string]any{"iss": "https://clerk.com"}, kid)

	// Verifying without providing a key returns an error.
	_, err := Verify(ctx, &VerifyParams{
		Token: token,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing json web key")

	// Verifying with wrong public key for the key.
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK: &clerk.JSONWebKey{
			Key:       nil,
			KeyID:     kid,
			Algorithm: string(jose.EdDSA),
			Use:       "sig",
		},
	})
	require.Error(t, err)

	// Verifying with wrong algorithm for the key.
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK: &clerk.JSONWebKey{
			Key:       pubKey,
			KeyID:     kid,
			Algorithm: string(jose.EdDSA),
			Use:       "sig",
		},
	})
	require.Error(t, err)

	// Verify with correct JSON web key.
	validKey := &clerk.JSONWebKey{
		Key:       pubKey,
		KeyID:     kid,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK:   validKey,
	})
	require.NoError(t, err)

	// Try an invalid token.
	_, err = Verify(ctx, &VerifyParams{
		Token: "this-is-not-a-token",
		JWK:   validKey,
	})
	require.Error(t, err)

	// Generate a token with an invalid issuer
	token, pubKey = clerktest.GenerateJWT(t, map[string]any{"iss": "https://whatever.com"}, kid)
	// Cannot verify if token has invalid issuer
	validKey = &clerk.JSONWebKey{
		Key:       pubKey,
		KeyID:     kid,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK:   validKey,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "issuer")
	// Satellite domains don't validate the issuer
	_, err = Verify(ctx, &VerifyParams{
		Token:       token,
		JWK:         validKey,
		IsSatellite: true,
	})
	require.NoError(t, err)
	// Issuer must match the proxy
	_, err = Verify(ctx, &VerifyParams{
		Token:    token,
		JWK:      validKey,
		ProxyURL: clerk.String("https://whatever.com"),
	})
	require.NoError(t, err)
	_, err = Verify(ctx, &VerifyParams{
		Token:    token,
		JWK:      validKey,
		ProxyURL: clerk.String("https://another.com/proxy"),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "issuer")

	// Generate a token with the 'azp' claim.
	token, pubKey = clerktest.GenerateJWT(
		t,
		map[string]any{
			"iss": "https://clerk.com",
			"azp": "whatever.com",
		},
		kid,
	)
	// Cannot verify if 'azp' does not match
	validKey = &clerk.JSONWebKey{
		Key:       pubKey,
		KeyID:     kid,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK:   validKey,
		AuthorizedPartyHandler: func(azp string) bool {
			return azp == "clerk.com"
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "authorized party")
}

func TestVerify_PublicClaims(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"

	exp := time.Now().Add(10 * time.Hour).Unix()
	nbf := time.Now().Add(-10 * time.Hour).Unix()
	// Generate a JWT for the following custom claims.
	tokenClaims := map[string]any{
		"orgs": map[string]string{
			"org_123": "org:admin",
			"org_456": "org:member",
		},
		"org_id":          "org_123",
		"org_role":        "org:admin",
		"org_permissions": []string{"org:create"},
		"org_slug":        "acmeinc",
		"sid":             "sess_123",
		"sub":             "user_123",
		"iss":             "https://clerk.com",
		"nbf":             nbf,
		"exp":             exp,
	}
	token, pubKey := clerktest.GenerateJWT(t, tokenClaims, kid)
	claims, err := Verify(ctx, &VerifyParams{
		Token: token,
		JWK: &clerk.JSONWebKey{
			Key:       pubKey,
			KeyID:     kid,
			Algorithm: string(jose.RS256),
			Use:       "sig",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "sess_123", claims.SessionID)
	require.Equal(t, "user_123", claims.Subject)
	require.Equal(t, "org_123", claims.ActiveOrganizationID)
	require.Equal(t, "acmeinc", claims.ActiveOrganizationSlug)
	require.Equal(t, "org:admin", claims.ActiveOrganizationRole)
	require.Equal(t, 1, len(claims.ActiveOrganizationPermissions))
	require.Equal(t, "org:create", claims.ActiveOrganizationPermissions[0])
	require.NotNil(t, claims.NotBefore)
}

// TestVerify_TimeValues tests that Verify validates that the token's
// not before (nbf) and expiry (exp) claims are respected.
func TestVerify_TimeValues(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"

	// Generate a JWT that's expired.
	exp := time.Now().Add(-1 * time.Minute).Unix()
	tokenClaims := map[string]any{
		"iss": "https://clerk.com",
		"exp": exp,
	}
	token, pubKey := clerktest.GenerateJWT(t, tokenClaims, kid)
	_, err := Verify(ctx, &VerifyParams{
		Token: token,
		JWK: &clerk.JSONWebKey{
			Key:       pubKey,
			KeyID:     kid,
			Algorithm: string(jose.RS256),
			Use:       "sig",
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "exp")

	// Generate a JWT that should be used after a date in the future.
	nbf := time.Now().Add(10 * time.Hour).Unix()
	tokenClaims = map[string]any{
		"iss": "https://clerk.com",
		"nbf": nbf,
	}
	token, pubKey = clerktest.GenerateJWT(t, tokenClaims, kid)
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK: &clerk.JSONWebKey{
			Key:       pubKey,
			KeyID:     kid,
			Algorithm: string(jose.RS256),
			Use:       "sig",
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "nbf")
}

type testCustomClaims struct {
	Domain      string `json:"domain"`
	Environment string `json:"environment"`
}

func TestVerify_CustomClaims(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"
	// Generate a JWT for the following custom claims.
	tokenClaims := map[string]any{
		"domain":      "clerk.com",
		"environment": "production",
		"sub":         "user_123",
		"iss":         "https://clerk.com",
	}
	token, pubKey := clerktest.GenerateJWT(t, tokenClaims, kid)

	customClaimsConstructor := func(_ context.Context) any {
		return &testCustomClaims{}
	}
	claims, err := Verify(ctx, &VerifyParams{
		Token: token,
		JWK: &clerk.JSONWebKey{
			Key:       pubKey,
			KeyID:     kid,
			Algorithm: string(jose.RS256),
			Use:       "sig",
		},
		CustomClaimsConstructor: customClaimsConstructor,
	})
	require.NoError(t, err)
	customClaims, ok := claims.Custom.(*testCustomClaims)
	require.True(t, ok)
	require.Equal(t, "user_123", claims.Subject)
	require.Equal(t, "clerk.com", customClaims.Domain)
	require.Equal(t, "production", customClaims.Environment)
}

func TestDecode_KeyID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"

	// Generate a JWT for the following custom claims.
	tokenClaims := map[string]any{
		"org_slug": "acmeinc",
		"sid":      "sess_123",
		"sub":      "user_123",
		"iss":      "https://clerk.com",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
	claims, err := Decode(ctx, &DecodeParams{
		Token: token,
	})
	require.NoError(t, err)
	require.Equal(t, "sess_123", claims.Extra["sid"])
	require.Equal(t, "acmeinc", claims.Extra["org_slug"])
	require.Equal(t, "user_123", claims.Subject)
	require.Equal(t, "https://clerk.com", claims.Issuer)
	require.Equal(t, kid, claims.KeyID)
}

// TestDecode_IsUnsafe tests that Decode does not do any validations
// on the token, like checking if it's expired.
func TestDecode_IsUnsafe(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"

	// Generate a JWT that's expired and should be used in the future.
	// What matters is that time values are invalid.
	exp := time.Now().Add(-1 * time.Minute).Unix()
	nbf := time.Now().Add(20 * time.Minute).Unix()
	tokenClaims := map[string]any{
		"iss": "https://clerk.com",
		"exp": exp,
		"nbf": nbf,
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
	claims, err := Decode(ctx, &DecodeParams{
		Token: token,
	})
	require.NoError(t, err)
	require.Equal(t, "https://clerk.com", claims.Issuer)
}
