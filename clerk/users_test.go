package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_ListAll_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := "[" + dummyUserJson + "]"

	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want []User
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Users.ListAll()
	if len(got) != len(want) {
		t.Errorf("Was expecting %d user to be returned, instead got %d", len(want), len(got))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestUsersService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	users, err := client.Users.ListAll()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if users != nil {
		t.Errorf("Was not expecting any users to be returned, instead got %v", users)
	}
}

const dummyUserJson = `{
        "birthday": "",
        "created_at": "2021-01-05T14:29:48.385449Z",
        "email_addresses": [
            {
                "email_address": "iron_man@avengers.com",
                "id": "idn_1mebQ9KkZWrhbQGF8Yj",
                "linked_to": [
                    {
                        "id": "idn_1mebQ8sPZOtb7UQgptk",
                        "type": "oauth_google"
                    }
                ],
                "object": "email_address",
                "verification": {
                    "status": "verified",
                    "strategy": "from_oauth_google"
                }
            }
        ],
        "external_accounts": [
            {
                "approved_scopes": "email https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile openid profile",
                "email_address": "iron_man@avengers.com",
                "family_name": "Stark",
                "given_name": "Tony",
                "google_id": "11031040442607",
                "id": "idn_1mebQ8sPZOtb7UQgptk",
                "object": "google_account",
                "picture": "https://lh3.googleusercontent.com/a-/AOh14uJQzsltH-3r-VQ=s96-c"
            }
        ],
        "first_name": "Anthony",
        "gender": "",
        "id": "user_1mebQggrD3xO5JfuHk7clQ94ysA",
        "last_name": "Stark",
        "metadata": {},
        "object": "user",
        "password_enabled": false,
        "phone_numbers": [],
        "primary_email_address_id": "idn_1mebQ9KkZWrhbQGF8Yj",
        "primary_phone_number_id": null,
        "private_metadata": {},
        "profile_image_url": "https://lh3.googleusercontent.com/a-/AOh14uJQzsltH-3r-VQ=s96-c",
        "two_factor_enabled": false,
        "updated_at": "2021-01-05T14:29:48.385449Z",
        "username": null
    }`
