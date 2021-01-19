package clerk

import (
	"context"
	"net/http"
)

const (
	ActiveSession = iota
)

func Middleware(client Client) func(handler http.Handler) http.Handler {
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
