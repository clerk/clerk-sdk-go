package phonenumber

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

func TestPhoneNumberClientCreate(t *testing.T) {
	t.Parallel()
	phone := "+10123456789"
	userID := "user_123"
	id := "idn_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"phone_number":"%s","user_id":"%s","verified":false}`, phone, userID)),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","phone_number":"%s"}`, id, phone)),
			Method: http.MethodPost,
			Path:   "/v1/phone_numbers",
		},
	}
	client := NewClient(config)
	phoneNumber, err := client.Create(context.Background(), &CreateParams{
		UserID:      clerk.String(userID),
		PhoneNumber: clerk.String(phone),
		Verified:    clerk.Bool(false),
	})
	require.NoError(t, err)
	require.Equal(t, id, phoneNumber.ID)
	require.Equal(t, phone, phoneNumber.PhoneNumber)
}

func TestPhoneNumberClientUpdate(t *testing.T) {
	t.Parallel()
	id := "idn_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(`{"verified":true, "reserved_for_second_factor":true}`),
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","reserved_for_second_factor":true,"verification":{"status":"verified"}}`, id)),
			Method: http.MethodPatch,
			Path:   "/v1/phone_numbers/" + id,
		},
	}
	client := NewClient(config)
	phoneNumber, err := client.Update(context.Background(), "idn_123", &UpdateParams{
		Verified:                clerk.Bool(true),
		ReservedForSecondFactor: clerk.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, id, phoneNumber.ID)
	require.Equal(t, "verified", phoneNumber.Verification.Status)
	require.Equal(t, true, phoneNumber.ReservedForSecondFactor)
}

func TestPhoneNumberClientGet(t *testing.T) {
	t.Parallel()
	id := "idn_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","verification":{"status":"verified"}}`, id)),
			Method: http.MethodGet,
			Path:   "/v1/phone_numbers/" + id,
		},
	}
	client := NewClient(config)
	phoneNumber, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, phoneNumber.ID)
	require.Equal(t, "verified", phoneNumber.Verification.Status)
}

func TestPhoneNumberClientDelete(t *testing.T) {
	t.Parallel()
	id := "idn_456"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)),
			Method: http.MethodDelete,
			Path:   "/v1/phone_numbers/" + id,
		},
	}
	client := NewClient(config)
	phoneNumber, err := client.Delete(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, phoneNumber.ID)
	require.True(t, phoneNumber.Deleted)
}
