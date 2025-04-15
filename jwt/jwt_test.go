package jwt

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/clerk/clerk-sdk-go/v3/jwks"
	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerify_InvalidParams(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"
	token, pubKey := clerktest.GenerateJWT(t, map[string]any{"iss": "https://clerk.com"}, kid)

	// Verifying with wrong public key for the key.
	_, err := Verify(ctx, &VerifyParams{
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
		"fva":             [2]int64{10, -1},
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
	require.Equal(t, 1, claims.Version)
	require.Equal(t, "sess_123", claims.SessionID)
	require.Equal(t, "user_123", claims.Subject)
	require.Equal(t, "org_123", claims.ActiveOrganizationID)
	require.Equal(t, "acmeinc", claims.ActiveOrganizationSlug)
	require.Equal(t, "org:admin", claims.ActiveOrganizationRole)
	require.Equal(t, 1, len(claims.ActiveOrganizationPermissions))
	require.Equal(t, "org:create", claims.ActiveOrganizationPermissions[0])
	require.NotNil(t, claims.NotBefore)
	require.Equal(t, int64(10), claims.FactorVerificationAge[0])
	require.Equal(t, int64(-1), claims.FactorVerificationAge[1])
}

func TestVerify_PublicClaimsVersion2(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	kid := "kid"

	exp := time.Now().Add(10 * time.Hour).Unix()
	nbf := time.Now().Add(-10 * time.Hour).Unix()
	// Generate a JWT for the following custom claims.
	tokenClaims := map[string]any{
		"v":   2,
		"fea": "o:repositories,o:projects,o:impersonation",
		"o": map[string]string{
			"id":  "org_123",
			"slg": "acmeinc",
			"rol": "admin",
			"per": "read,manage",
			"fpm": "1,2",
		},
		"fva": [2]int64{10, -1},
		"sid": "sess_123",
		"sub": "user_123",
		"iss": "https://clerk.com",
		"nbf": nbf,
		"exp": exp,
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
	require.Equal(t, []string{"org:repositories:read", "org:projects:manage"}, claims.ActiveOrganizationPermissions)
	require.NotNil(t, claims.NotBefore)
	require.Equal(t, int64(10), claims.FactorVerificationAge[0])
	require.Equal(t, int64(-1), claims.FactorVerificationAge[1])
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

// TestVerify_UsesTheJWKSClient tests that when verifying a JWT if
// you don't provide the JWK, the Verify method will make a request
// to GET /jwks to fetch the JWK set.
func TestVerify_UsesTheJWKSClient(t *testing.T) {
	t.Parallel()
	kid := "kid"
	totalJWKSRequests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			require.Equal(t, "custom client was used", r.Header.Get("X-Clerk-Application"))
			// Count the number of requests to the JWKS endpoint
			totalJWKSRequests++
			_, err := w.Write([]byte(
				fmt.Sprintf(
					`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
					kid,
				),
			))
			require.NoError(t, err)
			return
		}
	}))
	defer ts.Close()

	config := &clerk.ClientConfig{}
	config.HTTPClient = ts.Client()
	config.URL = &ts.URL
	config.CustomRequestHeaders = &clerk.CustomRequestHeaders{
		Application: "custom client was used",
	}
	jwksClient := jwks.NewClient(config)

	tokenClaims := map[string]any{
		"sid": "sess_123",
		"sub": "user_123",
		"iss": "https://clerk.com",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
	_, _ = Verify(context.Background(), &VerifyParams{
		Token:      token,
		JWKSClient: jwksClient,
	})
	// A request was made to fetch the JWKS
	require.Equal(t, 1, totalJWKSRequests)
}

// TestVerify_DefaultJWKSClient tests that when verifying a JWT if
// you don't provide the JWK and you dont' provide a jwks.Client,
// the Verify method will initialize a new jwks.Client with the
// default Backend and use it to fetch the JWK set.
func TestVerify_DefaultJWKSClient(t *testing.T) {
	kid := "kid"
	totalJWKSRequests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			// Count the number of requests to the JWKS endpoint
			totalJWKSRequests++
			_, err := w.Write([]byte(
				fmt.Sprintf(
					`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
					kid,
				),
			))
			require.NoError(t, err)
			return
		}
	}))
	defer ts.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	tokenClaims := map[string]any{
		"sid": "sess_123",
		"sub": "user_123",
		"iss": "https://clerk.com",
	}
	token, _ := clerktest.GenerateJWT(t, tokenClaims, kid)
	_, _ = Verify(context.Background(), &VerifyParams{
		Token: token,
	})
	// A request was made to fetch the JWKS
	require.Equal(t, 1, totalJWKSRequests)
}

func TestVerify_SessionTokenJWKv2(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name                   string
		Claims                 map[string]any
		ExpectedOrgID          string
		ExpectedOrgRole        string
		ExpectedOrgPermissions []string
	}{
		{
			Name: "only feature claims",
			Claims: map[string]any{
				"fea": "u:impersonation,u:repositories",
			},
			ExpectedOrgRole:        "",
			ExpectedOrgPermissions: nil,
		},
		{
			Name: "only org claims",
			Claims: map[string]any{
				"fea": "",
				"o": map[string]any{
					"id":  "org_xxx",
					"slg": "foobar",
					"rol": "admin",
				},
			},
			ExpectedOrgRole:        "org:admin",
			ExpectedOrgPermissions: []string{},
		},
		{
			Name: "org claims with only features",
			Claims: map[string]any{
				"fea": "o:impersonation,o:repositories",
				"o": map[string]any{
					"id":  "org_xxx",
					"slg": "foobar",
					"rol": "admin",
				},
			},
			ExpectedOrgRole:        "org:admin",
			ExpectedOrgPermissions: []string{},
		},
		{
			Name: "org claims with features and permissions",
			Claims: map[string]any{
				"fea": "o:repositories,o:projects,o:webhooks,o:impersonation",
				"o": map[string]any{
					"id":  "org_xxx",
					"slg": "foobar",
					"rol": "admin",
					"per": "read,manage",
					"fpm": "1,2,3",
				},
			},
			ExpectedOrgRole: "org:admin",
			ExpectedOrgPermissions: []string{
				"org:repositories:read",
				"org:projects:manage",
				"org:webhooks:read",
				"org:webhooks:manage",
			},
		},
		{
			Name: "org claims with features and permissions and feature permission mappings",
			Claims: map[string]any{
				"fea": "o:repositories,o:projects,o:webhooks,o:impersonation",
				"o": map[string]any{
					"id":  "org_xxx",
					"slg": "foobar",
					"rol": "admin",
					"per": "read,create,update,delete,revoke",
					"fpm": "7,21",
				},
			},
			ExpectedOrgRole: "org:admin",
			ExpectedOrgPermissions: []string{
				"org:repositories:read",
				"org:repositories:create",
				"org:repositories:update",
				"org:projects:read",
				"org:projects:update",
				"org:projects:revoke",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			kid := "kid"

			claims := tc.Claims
			claims["iss"] = "https://clerk.com"
			claims["v"] = 2
			token, pubKey := clerktest.GenerateJWT(t, claims, kid)

			verifiedClaims, err := Verify(ctx, &VerifyParams{
				Token: token,
				JWK: &clerk.JSONWebKey{
					Key:       pubKey,
					KeyID:     kid,
					Algorithm: string(jose.RS256),
					Use:       "sig",
				},
			})
			require.NoError(t, err)

			assert.Equal(t, tc.ExpectedOrgRole, verifiedClaims.ActiveOrganizationRole)
			assert.Equal(t, tc.ExpectedOrgPermissions, verifiedClaims.ActiveOrganizationPermissions)
		})
	}
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

func TestGetJSONWebKey_DefaultJWKSClient(t *testing.T) {
	kid := "kid"
	totalJWKSRequests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			// Count the number of requests to the JWKS endpoint
			totalJWKSRequests++
			_, err := w.Write([]byte(
				fmt.Sprintf(
					`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
					kid,
				),
			))
			require.NoError(t, err)
			return
		}
	}))
	defer ts.Close()

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: ts.Client(),
		URL:        &ts.URL,
	}))

	_, _ = GetJSONWebKey(context.Background(), &GetJSONWebKeyParams{
		KeyID: kid,
	})
	// A request was made to fetch the JWKS
	require.Equal(t, 1, totalJWKSRequests)
}

func TestGetJSONWebKey_UsesTheJWKSClient(t *testing.T) {
	t.Parallel()
	kid := "kid"
	totalJWKSRequests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			require.Equal(t, "custom client was used", r.Header.Get("X-Clerk-Application"))
			// Count the number of requests to the JWKS endpoint
			totalJWKSRequests++
			_, err := w.Write([]byte(
				fmt.Sprintf(
					`{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"ypsS9Iq26F71B3lPjT_IMtglDXo8Dko9h5UBmrvkWo6pdH_4zmMjeghozaHY1aQf1dHUBLsov_XvG_t-1yf7tFfO_ImC1JqSQwdSjrXZp3oMNFHwdwAknvtlBg3sBxJ8nM1WaCWaTlb2JhEmczIji15UG6V0M2cAp2VK_brcylQROaJLC2zVa4usGi4AHzAHaRUTv6XB9bGYMvkM-ZniuXgp9dPurisIIWg25DGrTaH-kg8LPaqGwa54eLEnvfAe0ZH_MvA4_bn_u_iDkQ9ZI_CD1vwf0EDnzLgd9ZG1khGsqmXY_4WiLRGsPqZe90HzaBJma9sAxXB4qj_aNnwD5w","e":"AQAB"}]}`,
					kid,
				),
			))
			require.NoError(t, err)
			return
		}
	}))
	defer ts.Close()

	config := &clerk.ClientConfig{}
	config.HTTPClient = ts.Client()
	config.URL = &ts.URL
	config.CustomRequestHeaders = &clerk.CustomRequestHeaders{
		Application: "custom client was used",
	}
	jwksClient := jwks.NewClient(config)

	_, _ = GetJSONWebKey(context.Background(), &GetJSONWebKeyParams{
		KeyID:      kid,
		JWKSClient: jwksClient,
	})
	// A request was made to fetch the JWKS
	require.Equal(t, 1, totalJWKSRequests)
}

func TestToBits(t *testing.T) {
	t.Parallel()
	testCases := map[int][]int{
		0: {},
		1: {1},
		2: {0, 1},
		3: {1, 1},
		4: {0, 0, 1},
	}
	for input, expected := range testCases {
		assert.Equal(t, expected, toBits(input))
	}
}
