package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlocklistService_CreateIdentifier_happyPath(t *testing.T) {
	token := "token"
	var blocklistIdentifier BlocklistIdentifierResponse
	_ = json.Unmarshal([]byte(dummyBlocklistIdentifierJson), &blocklistIdentifier)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/blocklist_identifiers", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, dummyBlocklistIdentifierJson)
	})

	got, _ := client.Blocklists().CreateIdentifier(CreateBlocklistIdentifierParams{
		Identifier: blocklistIdentifier.Identifier,
	})

	assert.Equal(t, &blocklistIdentifier, got)
}

func TestBlocklistService_CreateIdentifier_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Blocklists().CreateIdentifier(CreateBlocklistIdentifierParams{
		Identifier: "dummy@example.com",
	})
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestBlocklistService_DeleteIdentifier_happyPath(t *testing.T) {
	token := "token"
	var blocklistIdentifier BlocklistIdentifierResponse
	_ = json.Unmarshal([]byte(dummyBlocklistIdentifierJson), &blocklistIdentifier)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/blocklist_identifiers/"+blocklistIdentifier.ID, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodDelete)
		testHeader(t, req, "Authorization", "Bearer "+token)
		response := fmt.Sprintf(`{ "deleted": true, "id": "%s", "object": "blocklist_identifier" }`, blocklistIdentifier.ID)
		_, _ = fmt.Fprint(w, response)
	})

	got, _ := client.Blocklists().DeleteIdentifier(blocklistIdentifier.ID)
	assert.Equal(t, blocklistIdentifier.ID, got.ID)
	assert.True(t, got.Deleted)
}

func TestBlocklistService_DeleteIdentifier_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Blocklists().DeleteIdentifier("blid_1mvFol71HiKCcypBd6xxg0IpMBN")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestBlocklistService_ListAllIdentifiers_happyPath(t *testing.T) {
	token := "token"
	var blocklistIdentifier BlocklistIdentifierResponse
	_ = json.Unmarshal([]byte(dummyBlocklistIdentifierJson), &blocklistIdentifier)

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/blocklist_identifiers", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, fmt.Sprintf(`{"data": [%s], "total_count": 1}`, dummyBlocklistIdentifierJson))
	})

	got, _ := client.Blocklists().ListAllIdentifiers()

	assert.Len(t, got.Data, 1)
	assert.Equal(t, int64(1), got.TotalCount)
	assert.Equal(t, &blocklistIdentifier, got.Data[0])
}

func TestBlocklistService_ListAllIdentifiers_invalidServer(t *testing.T) {
	client, _ := NewClient("token")

	_, err := client.Blocklists().ListAllIdentifiers()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

const dummyBlocklistIdentifierJson = `{
    "id": "blid_1mvFol71HiKCcypBd6xxg0IpMBN",
    "object": "blocklist_identifier",
	"identifier": "dummy@example.com",
	"identifier_type": "email_address",
	"created_at": 1610783813,
	"updated_at": 1610783813
}`
