package clerk

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// RequireSessionV2 will hijack the request and return an HTTP status 403
// if the session is not authenticated.
func RequireSessionV2(client Client, next http.Handler) http.Handler {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ActiveSessionClaims).(*SessionClaims)
		if !ok || claims == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})

	return WithSessionV2(client)(f)
}

// SessionFromContext returns the session's (if any) claims, as parsed from the
// token.
func SessionFromContext(ctx context.Context) (*SessionClaims, bool) {
	c, ok := ctx.Value(ActiveSessionClaims).(*SessionClaims)
	return c, ok
}

// WithSessionV2 is the new middleware that supports Auth v2. If the session is
// authenticated, it adds the corresponding session claims found in the JWT to
// request's context.
func WithSessionV2(client Client) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerToken := strings.TrimSpace(r.Header.Get("authorization"))
			cookieToken, _ := r.Cookie("__session")
			clientUat, _ := r.Cookie("__client_uat")

			// ****************************************************
			//                                                    *
			//                HEADER AUTHENTICATION               *
			//                                                    *
			// ****************************************************
			if headerToken != "" {
				_, err := client.DecodeToken(headerToken)
				if err != nil {
					// signed out
					next.ServeHTTP(w, r)
					return
				}

				claims, err := client.VerifyToken(headerToken)
				if err == nil { // signed in
					ctx := context.WithValue(r.Context(), ActiveSessionClaims, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				// Clerk.js should refresh the token and retry
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// in cross-origin requests the use of Authorization
			// header is mandatory
			if isCrossOrigin(r) {
				// signed out
				next.ServeHTTP(w, r)
				return
			}

			// ****************************************************
			//                                                    *
			//                COOKIE AUTHENTICATION               *
			//                                                    *
			// ****************************************************
			if isDevelopmentOrStaging(client) && (r.Referer() == "" || isCrossOrigin(r)) {
				renderInterstitial(client, w)
				return
			}

			if isProduction(client) && clientUat == nil {
				next.ServeHTTP(w, r)
				return
			}

			if clientUat != nil && clientUat.Value == "0" {
				next.ServeHTTP(w, r)
				return
			}

			if clientUat == nil {
				renderInterstitial(client, w)
				return
			}

			var clientUatTs int64
			ts, err := strconv.ParseInt(clientUat.Value, 10, 64)
			if err == nil {
				clientUatTs = ts
			}

			if cookieToken == nil {
				renderInterstitial(client, w)
				return
			}

			claims, err := client.VerifyToken(cookieToken.Value)
			if err == nil && claims.IssuedAt != nil && clientUatTs <= int64(*claims.IssuedAt) {
				ctx := context.WithValue(r.Context(), ActiveSessionClaims, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			renderInterstitial(client, w)
		})
	}
}

func isCrossOrigin(r *http.Request) bool {
	origin := r.Header.Get("origin")

	// r.Host may contain host:port, but we only want the host
	host := strings.Split(r.Host, ":")[0]

	return origin != "" && origin != host
}

func isDevelopmentOrStaging(c Client) bool {
	return strings.HasPrefix(c.APIKey(), "test_")
}

func isProduction(c Client) bool {
	return !isDevelopmentOrStaging(c)
}

func renderInterstitial(c Client, w http.ResponseWriter) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(401)
	resp, _ := c.Interstitial()
	w.Write(resp)
}
