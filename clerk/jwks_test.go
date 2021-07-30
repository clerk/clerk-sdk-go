package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestJWKSService_ListAll_happyPath(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/jwks", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, dummyJWKSJson)
	})

	want := &JWKS{}
	_ = json.Unmarshal([]byte(dummyJWKSJson), want)

	got, _ := c.JWKS().ListAll()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestJWKSService_ListAll_invalidServer(t *testing.T) {
	c, _ := NewClient("token")

	jwks, err := c.JWKS().ListAll()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if jwks != nil {
		t.Errorf("Was not expecting any jwks to be returned, instead got %+v", jwks)
	}
}

const dummyJWKSJson = `
{
	"keys":
		[
			{
				"use":"sig",
				"kty":"RSA",
				"kid":"kid",
				"alg":"RS256",
				"n":"8ffrRMLd1z50B1hJcEfoxPac2wm9U_SXCnoXxSg5frRyW1oI1t9e78y8sOOwUt-IU4FXNcNK93dsCDQMeDBc6EfLxPBHuCB4SbVvsbpdMH8XSy9qLH6AJmS1GqOldYG0VkP1YzSwGXTkflgcDLCtYOHxkjiK6m5TnhJ4tu77bkjPrINiWAo4jAYBCjk1gqiW3LZWZwzwvqF_7n8g50JbhoTiJi2z6rd0anSFgi1A9AbViKwlzdxkll1uW90W1kn_Zs6lC6Yz7-X9WmelhxxUoLVE49BcCQ82PtmlBvxDQk7rREPLRbvzJSI0RIw1HMChRkZC_KtsLNkgPKq5tY_YSw",
				"e":"AQAB"
			}
		]
	}
`
