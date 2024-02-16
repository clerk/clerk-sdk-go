package jwt

import (
	"context"
	"testing"

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
