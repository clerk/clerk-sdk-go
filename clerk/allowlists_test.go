package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllowlistService_CreateIdentifier_happyPath(t *testing.T) {
	token := "token"
	var allowlistIdentifier AllowlistIdentifierResponse
	_ = json.Unmarshal([]byte(dummyAllowlistIdentifierJson), &allowlistIdentifier)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/allowlist_identifiers", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyAllowlistIdentifierJson)
	})

	got, _ := client.Allowlists().CreateIdentifier(CreateAllowlistIdentifierParams{
		Identifier: allowlistIdentifier.Identifier,
	})

	assert.Equal(t, &allowlistIdentifier, got)
}

func TestAllowlistService_CreateIdentifier_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Allowlists().CreateIdentifier(CreateAllowlistIdentifierParams{
		Identifier: "dummy@example.com",
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestAllowlistService_DeleteIdentifier_happyPath(t *testing.T) {
	token := "token"
	var allowlistIdentifier BlocklistIdentifierResponse
	_ = json.Unmarshal([]byte(dummyBlocklistIdentifierJson), &allowlistIdentifier)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/allowlist_identifiers/"+allowlistIdentifier.ID, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodDelete)
		testHeader(t, req, "Authorization", "Bearer "+token)
		response := fmt.Sprintf(`{ "deleted": true, "id": "%s", "object": "allowlist_identifier" }`, allowlistIdentifier.ID)
		_, _ = fmt.Fprint(w, response)
	})

	got, _ := client.Allowlists().DeleteIdentifier(allowlistIdentifier.ID)
	assert.Equal(t, allowlistIdentifier.ID, got.ID)
	assert.True(t, got.Deleted)
}

func TestAllowlistService_DeleteIdentifier_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Allowlists().DeleteIdentifier("alid_1mvFol71HiKCcypBd6xxg0IpMBN")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestAllowlistService_ListAllIdentifiers_happyPath(t *testing.T) {
	token := "token"
	var allowlistIdentifier AllowlistIdentifierResponse
	_ = json.Unmarshal([]byte(dummyAllowlistIdentifierJson), &allowlistIdentifier)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/allowlist_identifiers", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, fmt.Sprintf(`[%s]`, dummyAllowlistIdentifierJson))
	})

	got, _ := client.Allowlists().ListAllIdentifiers()

	assert.Len(t, got.Data, 1)
	assert.Equal(t, int64(1), got.TotalCount)
	assert.Equal(t, &allowlistIdentifier, got.Data[0])
}

func TestAllowlistService_ListAllIdentifiers_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Allowlists().ListAllIdentifiers()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummyAllowlistIdentifierJson = `{
    "id": "alid_1mvFol71HiKCcypBd6xxg0IpMBN",
    "object": "allowlist_identifier",
	"identifier": "dummy@example.com",
	"invitation_id": "inv_1mvFol71PeRCcypBd628g0IuRmF",
	"created_at": 1610783813,
	"updated_at": 1610783813
}`
