package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsService_ListAll_HappyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	expectedResponse := fmt.Sprintf(`{
		"total_count": 1,
		"data": [%s]
	}`, domainJSON)

	mux.HandleFunc("/domains", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want *DomainListResponse
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, _ := client.Domains().ListAll()
	if got.TotalCount != want.TotalCount {
		t.Errorf("Was expecting %d domains to be returned, instead got %d", want.TotalCount, got.TotalCount)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestDomainsService_Create_HappyPath(t *testing.T) {
	token := "token"
	expectedResponse := domainJSON

	name := "foobar.com"

	payload := CreateDomainParams{
		Name:        name,
		IsSatellite: true,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want Domain
	err := json.Unmarshal([]byte(expectedResponse), &want)
	assert.Nil(t, err)

	got, err := client.Domains().Create(payload)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestDomainsService_Update_HappyPath(t *testing.T) {
	token := "token"
	domainID := "dmn_banana"
	expectedResponse := domainJSON

	name := "foobar.com"

	payload := UpdateDomainParams{
		Name: &name,
	}

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/domains/%s", domainID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PATCH")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want Domain
	err := json.Unmarshal([]byte(expectedResponse), &want)
	assert.Nil(t, err)

	got, err := client.Domains().Update(domainID, payload)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestDomainsService_Delete_HappyPath(t *testing.T) {
	token := "token"
	domainID := "dmn_banana"
	expectedResponse := `{
		"object": "domain",
		"id": "dmn_banana",
		"delete": "true"
	}`

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/domains/%s", domainID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want DeleteResponse
	err := json.Unmarshal([]byte(expectedResponse), &want)
	assert.Nil(t, err)

	got, err := client.Domains().Delete(domainID)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %#v, want %#v", got, want)
	}
}

const domainJSON = `{
	"id": "dmn_banana",
	"name": "foobar.com"
}`
