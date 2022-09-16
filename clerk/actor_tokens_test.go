package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActorTokenService_CreateIdentifier_happyPath(t *testing.T) {
	token := "token"
	var actorTokenResponse ActorTokenResponse
	_ = json.Unmarshal([]byte(dummyActorTokenJson), &actorTokenResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/actor_tokens", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyActorTokenJson)
	})

	got, _ := client.ActorTokens().Create(CreateActorTokenParams{
		Actor:  actorTokenResponse.Actor,
		UserID: actorTokenResponse.UserID,
	})

	assert.Equal(t, &actorTokenResponse, got)
}

func TestActorTokenService_Create_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.ActorTokens().Create(CreateActorTokenParams{
		Actor:  []byte(`{"sub":"some_actor_id"}`),
		UserID: "some_user_id",
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestActorTokenService_Revoke_happyPath(t *testing.T) {
	token := "token"
	var actorTokenResponse ActorTokenResponse
	_ = json.Unmarshal([]byte(dummyActorTokenJson), &actorTokenResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/actor_tokens/"+actorTokenResponse.ID+"/revoke", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		_, _ = fmt.Fprint(w, dummyActorTokenJson)
	})

	got, _ := client.ActorTokens().Revoke(actorTokenResponse.ID)
	assert.Equal(t, &actorTokenResponse, got)
}

func TestActorTokenService_Revoke_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.ActorTokens().Revoke("impt_2EKxJqKTYcBMlzMh6BGe2C7kh6b")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummyActorTokenJson = `{
    "id": "impt_2EKxJqKTYcBMlzMh6BGe2C7kh6b",
    "object": "actor_token",
    "actor": {
			"sub": "some_actor_id",
			"iss": "the-issuer"
		},
    "user_id": "user_2EKwinD5cZID96QP1ruHsnfGx50",
    "status": "pending",
    "token": "my-token",
    "created_at": 1662358007949,
    "updated_at": 1662358007949
}`
