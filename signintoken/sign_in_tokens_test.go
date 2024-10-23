package signintoken

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

func TestSignInTokenTokenCreate(t *testing.T) {
	userID := "usr_123"
	id := "sign_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"user_id":"%s"}`, userID)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","user_id":"%s"}`, id, userID)),
				Path:   "/v1/sign_in_tokens",
				Method: http.MethodPost,
			},
		},
	}))

	signInToken, err := Create(context.Background(), &CreateParams{
		UserID: clerk.String(userID),
	})
	require.NoError(t, err)
	require.Equal(t, id, signInToken.ID)
	require.Equal(t, userID, signInToken.UserID)
}

func TestSignInTokenRevoke(t *testing.T) {
	id := "sign_456"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","status":"revoked"}`, id)),
				Path:   fmt.Sprintf("/v1/sign_in_tokens/%s/revoke", id),
				Method: http.MethodPost,
			},
		},
	}))

	signInToken, err := Revoke(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, signInToken.ID)
	require.Equal(t, "revoked", signInToken.Status)
}
