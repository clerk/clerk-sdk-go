// Package http provides HTTP utilities and handler middleware.
package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

// RequireHeaderAuthorization will respond with HTTP 403 Forbidden if
// the Authorization header does not contain a valid session token.
func RequireHeaderAuthorization(opts ...AuthorizationOption) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WithHeaderAuthorization(opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok || claims == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

// WithHeaderAuthorization checks the Authorization request header
// for a valid Clerk authorization JWT. The token is parsed and verified
// and the active session claims are written to the http.Request context.
// The middleware uses Bearer authentication, so the Authorization header
// is expected to have the following format:
// Authorization: Bearer <token>
func WithHeaderAuthorization(opts ...AuthorizationOption) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := strings.TrimSpace(r.Header.Get("Authorization"))
			if authorization == "" {
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(authorization, "Bearer ")
			_, err := jwt.Decode(r.Context(), &jwt.DecodeParams{Token: token})
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			params := &AuthorizationParams{}
			for _, opt := range opts {
				err = opt(params)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			params.Token = token
			claims, err := jwt.Verify(r.Context(), &params.VerifyParams)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Token was verified. Add the session claims to the request context
			newCtx := clerk.ContextWithSessionClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

type AuthorizationParams struct {
	jwt.VerifyParams
}

// AuthorizationOption is a functional parameter for configuring
// authorization options.
type AuthorizationOption func(*AuthorizationParams) error

// AuthorizedParty sets the authorized parties that will be checked
// against the azp JWT claim.
func AuthorizedParty(parties ...string) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.SetAuthorizedParties(parties...)
		return nil
	}
}

// CustomClaims allows to pass a type (e.g. struct), which will be populated with the token claims based on json tags.
// You must pass a pointer for this option to work.
func CustomClaims(claims any) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.CustomClaims = claims
		return nil
	}
}

// Leeway allows to set a custom leeway when comparing time values
// for JWT verification.
// The leeway gives some extra time to the token. That is, if the
// token is expired, it will still be accepted for 'leeway' amount
// of time.
// This option accomodates for clock skew.
func Leeway(leeway time.Duration) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.Leeway = leeway
		return nil
	}
}

// ProxyURL can be used to set the URL that proxies the Clerk Frontend
// API. Useful for proxy based setups.
// See https://clerk.com/docs/advanced-usage/using-proxies
func ProxyURL(proxyURL string) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.ProxyURL = clerk.String(proxyURL)
		return nil
	}
}

// Satellite can be used to signify that the authorization happens
// on a satellite domain.
// See https://clerk.com/docs/advanced-usage/satellite-domains
func Satellite(isSatellite bool) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		params.IsSatellite = isSatellite
		return nil
	}
}

// JSONWebKey allows to provide a custom JSON Web Key (JWK) based on
// which the authorization JWT will be verified.
// When verifying the authorization JWT without a custom key, the JWK
// will be fetched from the Clerk API and cached for one hour, then
// the JWK will be fetched again from the Clerk API.
// Passing a custom JSON Web Key means that no request to fetch JSON
// web keys will be made. It's the caller's responsibility to refresh
// the JWK when keys are rolled.
func JSONWebKey(key string) AuthorizationOption {
	return func(params *AuthorizationParams) error {
		// From the Clerk docs: "Note that the JWT Verification key is not in
		// PEM format, the header and footer are missing, in order to be shorter
		// and single-line for easier setup."
		if !strings.HasPrefix(key, "-----BEGIN") {
			key = "-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----"
		}
		jwk, err := clerk.JSONWebKeyFromPEM(key)
		if err != nil {
			return err
		}
		params.JWK = jwk
		return nil
	}
}
