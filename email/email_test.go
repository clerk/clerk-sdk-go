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
		To:      Recipient{Address: "admin@acme.com"},
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

// TestEmailSendToUserID asserts the user_id recipient form serializes to
// {"to":{"user_id":...}} (no stray "address" key) and surfaces the server's
// resolved address plus the user/email_address linkage on the response.
func TestEmailSendToUserID(t *testing.T) {
	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: &http.Client{
			Transport: &clerktest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"to": {"user_id": "user_123"},
					"from": {"address": "noreply@acme.com"},
					"subject": "Hello",
					"html": "<p>hi</p>"
				}`),
				Out: json.RawMessage(`{
					"id": "ema_123",
					"object": "email",
					"from_email_name": "noreply",
					"to_email_address": "member@acme.com",
					"email_address_id": "idn_123",
					"user_id": "user_123",
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
		To:      Recipient{UserID: "user_123"},
		From:    Mailbox{Address: "noreply@acme.com"},
		Subject: "Hello",
		HTML:    "<p>hi</p>",
	})
	require.NoError(t, err)
	require.Equal(t, "member@acme.com", email.ToEmailAddress)
	require.NotNil(t, email.UserID)
	require.Equal(t, "user_123", *email.UserID)
	require.NotNil(t, email.EmailAddressID)
	require.Equal(t, "idn_123", *email.EmailAddressID)
}
