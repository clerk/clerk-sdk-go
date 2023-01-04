package clerk

import (
	"context"
	"net/http"
	"strings"
)

const (
	ActiveSession = iota
	ActiveSessionClaims

// TODO: we should use a type alias instead of int, so as to avoid collisions
// with other packages
)

func WithSession(client Client) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token, isAuthV2 := isAuthV2Request(r, client); isAuthV2 {
				// Validate using session token
				claims, err := client.VerifyToken(token)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte(err.Error()))
					return
				}

				ctx := context.WithValue(r.Context(), ActiveSessionClaims, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				// Validate using session verify request
				session, err := client.Verification().Verify(r)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(err.Error()))
					return
				}

				ctx := context.WithValue(r.Context(), ActiveSession, session)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func isAuthV2Request(r *http.Request, client Client) (string, bool) {
	// Try with token from header
	headerToken := strings.TrimSpace(r.Header.Get("Authorization"))
	headerToken = strings.TrimPrefix(headerToken, "Bearer ")

	claims, err := client.DecodeToken(headerToken)
	if err == nil {
		return headerToken, isValidIssuer(claims.Issuer)
	}

	// Verification from header token failed, try with token from cookie
	cookieSession, err := r.Cookie(CookieSession)
	if err != nil {
		return "", false
	}

	claims, err = client.DecodeToken(cookieSession.Value)
	if err != nil {
		return "", false
	}

	return cookieSession.Value, isValidIssuer(claims.Issuer)
}
