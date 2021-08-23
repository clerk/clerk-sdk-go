package clerk

import (
	"context"
	"errors"
	"net/http"
)

const (
	ActiveSession = iota
	ActiveClaims
)

func WithSession(client Client) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			session, err := client.Verification().Verify(request)
			if err != nil {
				writer.WriteHeader(400)
				writer.Write([]byte(err.Error()))
				return
			}

			updatedRequest := addSessionToContext(request, session)
			next.ServeHTTP(writer, updatedRequest)
		})
	}
}

func addSessionToContext(request *http.Request, session *Session) *http.Request {
	ctx := request.Context()
	updated := context.WithValue(ctx, ActiveSession, session)
	return request.WithContext(updated)
}

func WithNetworkless(client Client) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := retrieveTokenFromRequest(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			claims, err := client.VerifyToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			ctx := context.WithValue(r.Context(), ActiveClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func retrieveTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token, nil
	}

	for _, cookie := range r.Cookies() {
		if cookie.Name == "__session" {
			return cookie.Value, nil
		}
	}

	return "", errors.New("no token found in request header or cookie")
}
