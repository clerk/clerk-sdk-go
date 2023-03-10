package clerk

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/square/go-jose.v2/jwt"
)

var urlSchemeRe = regexp.MustCompile(`(^\w+:|^)\/\/`)

// RequireSessionV2 will hijack the request and return an HTTP status 403
// if the session is not authenticated.
func RequireSessionV2(client Client, verifyTokenOptions ...VerifyTokenOption) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ActiveSessionClaims).(*SessionClaims)
			if !ok || claims == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})

		return WithSessionV2(client, verifyTokenOptions...)(f)
	}
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
func WithSessionV2(client Client, verifyTokenOptions ...VerifyTokenOption) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ****************************************************
			//                                                    *
			//                HEADER AUTHENTICATION               *
			//                                                    *
			// ****************************************************
			_, authorizationHeaderExists := r.Header["Authorization"]

			if authorizationHeaderExists {
				headerToken := strings.TrimSpace(r.Header.Get("Authorization"))
				headerToken = strings.TrimPrefix(headerToken, "Bearer ")

				_, err := client.DecodeToken(headerToken)
				if err != nil {
					// signed out
					next.ServeHTTP(w, r)
					return
				}

				claims, err := client.VerifyToken(headerToken, verifyTokenOptions...)
				if err == nil { // signed in
					ctx := context.WithValue(r.Context(), ActiveSessionClaims, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				// Clerk.js should refresh the token and retry
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// In development or staging environments only, based on the request User Agent, detect non-browser
			// requests (e.g. scripts). If there is no Authorization header, consider the user as signed out
			// and prevent interstitial rendering
			if isDevelopmentOrStaging(client) && !strings.HasPrefix(r.UserAgent(), "Mozilla/") {
				// signed out
				next.ServeHTTP(w, r)
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
			cookieToken, _ := r.Cookie("__session")
			clientUat, _ := r.Cookie("__client_uat")

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

			claims, err := client.VerifyToken(cookieToken.Value, verifyTokenOptions...)

			if err == nil {
				if claims.IssuedAt != nil && clientUatTs <= int64(*claims.IssuedAt) {
					ctx := context.WithValue(r.Context(), ActiveSessionClaims, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				renderInterstitial(client, w)
				return
			}

			if errors.Is(err, jwt.ErrExpired) || errors.Is(err, jwt.ErrIssuedInTheFuture) {
				renderInterstitial(client, w)
				return
			}

			// signed out
			next.ServeHTTP(w, r)
			return
		})
	}
}

func isCrossOrigin(r *http.Request) bool {
	// origin contains scheme+host and optionally port (ommitted if 80 or 443)
	// ref. https://www.rfc-editor.org/rfc/rfc6454#section-6.1
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	origin = urlSchemeRe.ReplaceAllString(origin, "") // strip scheme
	if origin == "" {
		return false
	}

	// parse request's host and port, taking into account reverse proxies
	u := &url.URL{Host: r.Host}
	host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = u.Hostname()
	}
	port := strings.TrimSpace(r.Header.Get("X-Forwarded-Port"))
	if port == "" {
		port = u.Port()
	}

	if port != "" && port != "80" && port != "443" {
		host = net.JoinHostPort(host, port)
	}

	return origin != host
}

func isDevelopmentOrStaging(c Client) bool {
	return strings.HasPrefix(c.APIKey(), "test_") || strings.HasPrefix(c.APIKey(), "sk_test_")
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
