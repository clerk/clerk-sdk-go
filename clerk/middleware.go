package clerk

import (
	"context"
	"net/http"
	"strings"
)

const (
	ActiveSession = iota
	ActiveSessionClaims
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
	headerToken := r.Header.Get("Authorization")

	claims, err := client.DecodeToken(headerToken)
	if err == nil {
		return headerToken, strings.HasPrefix(claims.Issuer, "https://clerk.")
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

	return cookieSession.Value, strings.HasPrefix(claims.Issuer, "https://clerk.")
}
