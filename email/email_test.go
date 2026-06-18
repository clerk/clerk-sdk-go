package email

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestEmailSend(t *testing.T) {
	id := "ema_123"
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"to": {"address": "admin@acme.com"},
					"from": {"address": "noreply@acme.com"},
					"reply_to": {"address": "support@acme.com"},
					"subject": "Hello",
					"html": "<p>hi</p>"
				}`),
				Out: json.RawMessage(`{
					"id": "ema_123",
					"object": "email",
					"from_email_name": "noreply",
					"to_email_address": "admin@acme.com",
					"subject": "Hello",
					"body": "<p>hi</p>",
					"status": "queued",
					"delivered_by_clerk": true
				}`),
				Path:   "/v1/email",
				Method: http.MethodPost,
			},
		},
	}))

	email, err := Send(context.Background(), &SendParams{
		To:      Mailbox{Address: "admin@acme.com"},
		From:    Mailbox{Address: "noreply@acme.com"},
		ReplyTo: &Mailbox{Address: "support@acme.com"},
		Subject: "Hello",
		HTML:    "<p>hi</p>",
	})
	require.NoError(t, err)
	require.Equal(t, id, email.ID)
	require.Equal(t, "email", email.Object)
	require.Equal(t, "admin@acme.com", email.ToEmailAddress)
	require.Equal(t, "queued", email.Status)
	require.True(t, email.DeliveredByClerk)
}
