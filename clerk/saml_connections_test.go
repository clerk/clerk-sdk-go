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

func TestSAMLConnectionsService_ListAll(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	dummyResponse := fmt.Sprintf(`{
		"data": [%s],
		"total_count": 1
	}`, dummySAMLConnectionJSON)

	mux.HandleFunc("/saml_connections", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer token")

		expectedQuery := url.Values{
			"limit":    {"5"},
			"offset":   {"6"},
			"query":    {"my-query"},
			"order_by": {"created_at"},
		}
		assert.Equal(t, expectedQuery, req.URL.Query())

		_, _ = fmt.Fprint(w, dummyResponse)
	})

	listParams := ListSAMLConnectionsParams{
		Limit:   intToPtr(5),
		Offset:  intToPtr(6),
		Query:   stringToPtr("my-query"),
		OrderBy: stringToPtr("created_at"),
	}

	got, err := c.SAMLConnections().ListAll(listParams)
	assert.NoError(t, err)

	expected := &ListSAMLConnectionsResponse{}
	_ = json.Unmarshal([]byte(dummyResponse), expected)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestSAMLConnectionsService_Read(t *testing.T) {
	dummyResponse := dummySAMLConnectionJSON

	c, mux, _, teardown := setup("token")
	defer teardown()

	url := fmt.Sprintf("/saml_connections/%s", dummySAMLConnectionID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	got, err := c.SAMLConnections().Read(dummySAMLConnectionID)
	assert.NoError(t, err)

	expected := SAMLConnection{}
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestSAMLConnectionsService_Create(t *testing.T) {
	dummyResponse := dummySAMLConnectionJSON

	c, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/saml_connections", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	createParams := &CreateSAMLConnectionParams{
		Name:           "Testing SAML",
		Domain:         "example.com",
		IdpEntityID:    stringToPtr("test-idp-entity-id"),
		IdpSsoURL:      stringToPtr("https://example.com/saml/sso"),
		IdpCertificate: stringToPtr(dummySAMLConnectionCertificate),
	}

	got, err := c.SAMLConnections().Create(createParams)
	assert.NoError(t, err)

	expected := SAMLConnection{}
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestSAMLConnectionsService_Update(t *testing.T) {
	expectedName := "New name for Testing SAML"
	expectedActive := true
	expectedSyncUserAttributes := false
	dummyResponse := dummySAMLConnectionUpdatedJSON

	c, mux, _, teardown := setup("token")
	defer teardown()

	url := fmt.Sprintf("/saml_connections/%s", dummySAMLConnectionID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPatch)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	updateParams := &UpdateSAMLConnectionParams{
		Name:               &expectedName,
		Active:             &expectedActive,
		SyncUserAttributes: &expectedSyncUserAttributes,
	}

	got, err := c.SAMLConnections().Update(dummySAMLConnectionID, updateParams)
	assert.NoError(t, err)

	expected := SAMLConnection{}
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestSAMLConnectionsService_Delete(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	url := fmt.Sprintf("/saml_connections/%s", dummySAMLConnectionID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodDelete)
		testHeader(t, req, "Authorization", "Bearer token")
		response := fmt.Sprintf(`{ "deleted": true, "id": "%s", "object": "saml_connection" }`, dummySAMLConnectionID)
		_, _ = fmt.Fprint(w, response)
	})

	expected := DeleteResponse{
		ID:      dummySAMLConnectionID,
		Object:  "saml_connection",
		Deleted: true,
	}

	got, err := c.SAMLConnections().Delete(dummySAMLConnectionID)
	assert.NoError(t, err)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", *got, expected)
	}
}

const (
	dummySAMLConnectionID = "samlc_2P17P4pXsx8MmunM1pkeYeimDDd"

	dummySAMLConnectionJSON = `
{
    "object": "saml_connection",
	"id": "` + dummySAMLConnectionID + `",
    "name": "Testing SAML",
    "domain": "example.com",
	"idp_entity_id": "test-idp-entity-id",
	"idp_sso_url": "https://example.com/saml/sso",
	"idp_certificate": "` + dummySAMLConnectionCertificate + `",
	"acs_url": "` + "https://clerk.example.com/v1/saml/acs" + dummySAMLConnectionID + `",
	"sp_entity_id": "` + "https://clerk.example.com/acs" + dummySAMLConnectionID + `",
	"active": false,
	"provider": "saml_custom",
	"user_count": 3,
	"sync_user_attributes": true
}`

	dummySAMLConnectionUpdatedJSON = `
{
    "object": "saml_connection",
	"id": "` + dummySAMLConnectionID + `",
    "name": "New name for Testing SAML",
    "domain": "example.com",
	"idp_entity_id": "test-idp-entity-id",
	"idp_sso_url": "https://example.com/saml/sso",
	"idp_certificate": "` + dummySAMLConnectionCertificate + `",
	"acs_url": "` + "https://clerk.example.com/v1/saml/acs" + dummySAMLConnectionID + `",
	"sp_entity_id": "` + "https://clerk.example.com/acs" + dummySAMLConnectionID + `",
	"active": true,
	"provider": "saml_custom",
	"user_count": 3,
	"sync_user_attributes": false
}`

	dummySAMLConnectionCertificate = `MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=`
)
