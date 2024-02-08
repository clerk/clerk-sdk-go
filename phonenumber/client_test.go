package phonenumber

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhoneNumberClientCreate(t *testing.T) {
	t.Parallel()

	id := "idn_123"
	userID := "usr_123"
	verified := true
	phone := "+30210555555"

	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			In: json.RawMessage(fmt.Sprintf(`
{"phone_number":"+30210555555", "user_id":"%s", "verified": true}
			`, userID)),
			Out: json.RawMessage(fmt.Sprintf(`
{"object": "phone_number", "id":"%s", "phone_number":"%s", "reserved": false, "verification": null}`, id, userID)),
			Method: http.MethodPost,
			Path:   "/v1/phone_numbers",
		},
	}
	client := NewClient(config)

	phoneNumber, err := client.Create(context.Background(), &CreateParams{
		UserID:      &userID,
		PhoneNumber: &phone,
		Verified:    &verified,
	})
	require.NoError(t, err)
	assert.Equal(t, id, phoneNumber.ID)
	assert.Equal(t, "phone_number", phoneNumber.Object)
	assert.Equal(t, false, phoneNumber.Reserved)
	assert.Nil(t, phoneNumber.Verification)
}
