package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersService_Create_happyPath(t *testing.T) {
	token := "token"
	var payload CreateUserParams
	_ = json.Unmarshal([]byte(dummyCreateUserRequestJson), &payload)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyUserJson)
	})

	got, err := client.Users().Create(payload)

	var want User
	_ = json.Unmarshal([]byte(dummyUserJson), &want)

	assert.Nil(t, err)
	assert.Equal(t, want, *got)
}

func TestUsersService_Create_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Users().Create(CreateUserParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

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

	got, _ := client.Users().ListAll(ListAllUsersParams{})
	if len(got) != len(want) {
		t.Errorf("Was expecting %d user to be returned, instead got %d", len(want), len(got))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestUsersService_ListAll_happyPathWithParameters(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := "[" + dummyUserJson + "]"

	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")

		actualQuery := req.URL.Query()
		expectedQuery := url.Values(map[string][]string{
			"limit":         {"5"},
			"offset":        {"6"},
			"email_address": {"email1", "email2"},
			"phone_number":  {"phone1", "phone2"},
			"web3_wallet":   {"wallet1", "wallet2"},
			"username":      {"username1", "username2"},
			"user_id":       {"userid1", "userid2"},
			"query":         {"my-query"},
			"order_by":      {"created_at"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	var want []User
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	limit := 5
	offset := 6
	queryString := "my-query"
	orderBy := "created_at"
	got, _ := client.Users().ListAll(ListAllUsersParams{
		Limit:          &limit,
		Offset:         &offset,
		EmailAddresses: []string{"email1", "email2"},
		PhoneNumbers:   []string{"phone1", "phone2"},
		Web3Wallets:    []string{"wallet1", "wallet2"},
		Usernames:      []string{"username1", "username2"},
		UserIDs:        []string{"userid1", "userid2"},
		Query:          &queryString,
		OrderBy:        &orderBy,
	})
	if len(got) != len(want) {
		t.Errorf("Was expecting %d user to be returned, instead got %d", len(want), len(got))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestUsersService_ListAll_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	users, err := client.Users().ListAll(ListAllUsersParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if users != nil {
		t.Errorf("Was not expecting any users to be returned, instead got %v", users)
	}
}

func TestUsersService_Count_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyUserCountJson

	mux.HandleFunc("/users/count", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want UserCount
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Users().Count(ListAllUsersParams{})
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestUsersService_Count_happyPathWithParameters(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := dummyUserCountJson

	mux.HandleFunc("/users/count", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")

		actualQuery := req.URL.Query()
		expectedQuery := url.Values(map[string][]string{
			"email_address": {"email1", "email2"},
			"phone_number":  {"phone1", "phone2"},
			"web3_wallet":   {"wallet1", "wallet2"},
			"username":      {"username1", "username2"},
			"user_id":       {"userid1", "userid2"},
			"query":         {"my-query"},
		})
		assert.Equal(t, expectedQuery, actualQuery)
		fmt.Fprint(w, expectedResponse)
	})

	var want UserCount
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	queryString := "my-query"
	got, _ := client.Users().Count(ListAllUsersParams{
		EmailAddresses: []string{"email1", "email2"},
		PhoneNumbers:   []string{"phone1", "phone2"},
		Web3Wallets:    []string{"wallet1", "wallet2"},
		Usernames:      []string{"username1", "username2"},
		UserIDs:        []string{"userid1", "userid2"},
		Query:          &queryString,
	})
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestUsersService_Count_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	users, err := client.Users().Count(ListAllUsersParams{})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if users != nil {
		t.Errorf("Was not expecting any users to be returned, instead got %v", users)
	}
}

func TestUsersService_Read_happyPath(t *testing.T) {
	token := "token"
	userId := "someUserId"
	expectedResponse := dummyUserJson

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/users/"+userId, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want User
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Users().Read(userId)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestUsersService_Read_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	user, err := client.Users().Read("someUserId")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if user != nil {
		t.Errorf("Was not expecting any user to be returned, instead got %v", user)
	}
}

func TestUsersService_Delete_happyPath(t *testing.T) {
	token := "token"
	userId := "someUserId"

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/users/"+userId, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer "+token)
		response := fmt.Sprintf(`{ "deleted": true, "id": "%v", "object": "user" }`, userId)
		fmt.Fprint(w, response)
	})

	want := DeleteResponse{ID: userId, Object: "user", Deleted: true}

	got, _ := client.Users().Delete(userId)
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestUsersService_Delete_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	delResponse, err := client.Users().Delete("someUserId")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if delResponse != nil {
		t.Errorf("Was not expecting any reponse to be returned, instead got %v", delResponse)
	}
}

func TestUsersService_Update_happyPath(t *testing.T) {
	token := "token"
	userId := "someUserId"
	var payload UpdateUser
	_ = json.Unmarshal([]byte(dummyUpdateRequestJson), &payload)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/users/"+userId, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyUserJson)
	})

	got, _ := client.Users().Update(userId, &payload)

	var want User
	_ = json.Unmarshal([]byte(dummyUserJson), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, payload)
	}
}

func TestUsersService_Update_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Users().Update("someUserId", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestUsersService_UpdateMetadata_happyPath(t *testing.T) {
	token := "token"
	userId := "someUserId"
	var payload UpdateUserMetadata
	_ = json.Unmarshal([]byte(dummyUpdateMetadataRequestJson), &payload)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/users/"+userId+"/metadata", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPatch)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyUserJson)
	})

	got, _ := client.Users().UpdateMetadata(userId, &payload)

	var want User
	_ = json.Unmarshal([]byte(dummyUserJson), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, payload)
	}
}

func TestUsersService_UpdateMetadata_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Users().UpdateMetadata("someUserId", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestUsersService_DisableMFA_happyPath(t *testing.T) {
	token := "token"
	userID := "test-user-id"
	var payload UpdateUserMetadata
	_ = json.Unmarshal([]byte(dummyUpdateMetadataRequestJson), &payload)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/users/"+userID+"/mfa", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodDelete)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, `{"user_id":"`+userID+`"}`)
	})

	got, err := client.Users().DisableMFA(userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, got.UserID)
}

const dummyUserJson = `{
        "birthday": "",
        "created_at": 1610783813,
        "email_addresses": [
            {
                "email_address": "iron_man@avengers.com",
                "id": "idn_1mebQ9KkZWrhbQGF8Yj",
                "linked_to": [
                    {
                        "id": "idn_1n8tzrjmoKzHtQkaFe1pvK1OqLr",
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
                "picture": "https://lh3.googleusercontent.com/a-/AOh14Gg-UlYe7Pzd8vngVKdFlNCuGTn7cqxx=s96-c"
            }
        ],
        "first_name": "Anthony",
        "gender": "",
        "id": "user_1mebQggrD3xO5JfuHk7clQ94ysA",
        "last_name": "Stark",
        "object": "user",
        "password_enabled": false,
        "phone_numbers": [],
        "primary_email_address_id": "idn_1n8tzqi8K5ydvb1K7RJEKjT7Wb8",
        "primary_phone_number_id": null,
        "profile_image_url": "https://lh3.googleusercontent.com/a-/AOh14Gg-UlYe7PzddYKJRu2r8vGTn7cqxx=s96-c",
        "two_factor_enabled": false,
        "updated_at": 1610783813,
        "username": null,
		"public_metadata": {
			"address": {
				"street": "Pennsylvania Avenue",
				"number": "1600"
			}
		},
		"private_metadata": {
			"app_id": 5
		},
		"last_sign_in_at": 1610783813
    }`

const dummyCreateUserRequestJson = `{
		"first_name": "Tony",
		"last_name": "Stark",
		"email_address": ["email@example.com"],
		"phone_number": ["+30123456789"],
		"password": "new_password",
		"public_metadata": {
			"address": {
				"street": "Pennsylvania Avenue",
				"number": "1600"
			}
		},
		"private_metadata": {
			app_id: 5
		},
		"unsafe_metadata": {
			viewed_profile: true
		},
		"totp_secret": "AICJ3HCXKO4KOY6NDH6RII4E3ZYL5ZBH",
	}`

const dummyUpdateRequestJson = `{
		"first_name": "Tony",
		"last_name": "Stark",
		"primary_email_address_id": "some_image_id",
		"primary_phone_number_id": "some_phone_id",
		"profile_image": "some_profile_image",
		"password": "new_password",
		"public_metadata": {
			"address": {
				"street": "Pennsylvania Avenue",
				"number": "1600"
			}
		},
		"private_metadata": {
			app_id: 5
		},
		"unsafe_metadata": {
			viewed_profile: true
		},
	}`

const dummyUpdateMetadataRequestJson = `{
		"public_metadata": {
			"value": "public_value",
		},
		"private_metadata": {
			"contact_id": "some_contact_id",
		},
		"unsafe_metadata": {
			viewed_profile: true
		},
	}`

const dummyUserCountJson = `{
		"object": "total_count",
		"total_count": 2
	}`
