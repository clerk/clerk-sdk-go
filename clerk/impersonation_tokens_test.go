package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImpersonationTokenService_CreateIdentifier_happyPath(t *testing.T) {
	token := "token"
	var impersonationTokenResponse ImpersonationTokenRespose
	_ = json.Unmarshal([]byte(dummyImpersonationTokenJson), &impersonationTokenResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/impersonation_tokens", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyImpersonationTokenJson)
	})

	got, _ := client.ImpersonationTokens().Create(CreateImpersonationTokenParams{
		ActorID:   impersonationTokenResponse.ActorID,
		SubjectID: impersonationTokenResponse.SubjectID,
	})

	assert.Equal(t, &impersonationTokenResponse, got)
}

func TestImpersonationTokenService_Create_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.ImpersonationTokens().Create(CreateImpersonationTokenParams{
		ActorID:   "some_actor_id",
		SubjectID: "some_subject_id",
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestImpersonationTokenService_Revoke_happyPath(t *testing.T) {
	token := "token"
	var impersonationTokenResponse ImpersonationTokenRespose
	_ = json.Unmarshal([]byte(dummyImpersonationTokenJson), &impersonationTokenResponse)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/impersonation_tokens/"+impersonationTokenResponse.ID+"/revoke", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		_, _ = fmt.Fprint(w, dummyImpersonationTokenJson)
	})

	got, _ := client.ImpersonationTokens().Revoke(impersonationTokenResponse.ID)
	assert.Equal(t, &impersonationTokenResponse, got)
}

func TestImpersonationTokenService_Revoke_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.ImpersonationTokens().Revoke("impt_2EKxJqKTYcBMlzMh6BGe2C7kh6b")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummyImpersonationTokenJson = `{
    "id": "impt_2EKxJqKTYcBMlzMh6BGe2C7kh6b",
    "object": "impersonation_token",
    "actor_id": "some_actor_id",
    "subject_id": "user_2EKwinD5cZID96QP1ruHsnfGx50",
    "status": "pending",
    "token": "my-token",
    "created_at": 1662358007949,
    "updated_at": 1662358007949
}`
