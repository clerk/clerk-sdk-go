package actortoken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestActorTokenCreate(t *testing.T) {
	userID := "usr_123"
	id := "act_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"user_id":"%s"}`, userID)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","user_id":"%s"}`, id, userID)),
				Path:   "/v1/actor_tokens",
				Method: http.MethodPost,
			},
		},
	}))

	actorToken, err := Create(context.Background(), &CreateParams{
		UserID: clerk.String(userID),
	})
	require.NoError(t, err)
	require.Equal(t, id, actorToken.ID)
	require.Equal(t, userID, actorToken.UserID)
}

func TestActorTokenCreate_Error(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Status: http.StatusBadRequest,
				Out: json.RawMessage(`{
  "errors":[{
		"code":"create-error-code"
	}],
	"clerk_trace_id":"create-trace-id"
}`),
			},
		},
	}))

	_, err := Create(context.Background(), &CreateParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "create-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "create-error-code", apiErr.Errors[0].Code)
}

func TestActorTokenRevoke(t *testing.T) {
	id := "act_456"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"revoked"}`, id)),
				Path:   fmt.Sprintf("/v1/actor_tokens/%s/revoke", id),
				Method: http.MethodPost,
			},
		},
	}))

	actorToken, err := Revoke(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, actorToken.ID)
	require.Equal(t, "revoked", actorToken.Status)
}
