package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatesService_List_All_happyPath(t *testing.T) {
	client, mux, _, teardown := setup("token")
	defer teardown()

	templateType := "email"

	expectedResponse := "[" + dummyTemplateJSON + "]"

	url := fmt.Sprintf("/templates/%s", templateType)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer token")
		fmt.Fprint(w, expectedResponse)
	})

	var want []Template
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.Templates().ListAll(templateType)
	assert.Nil(t, err)

	if len(got) != len(want) {
		t.Errorf("Was expecting %d user to be returned, instead got %d", len(want), len(got))
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response = %v, want %v", got, want)
	}
}

func TestTemplatesService_Read_happyPath(t *testing.T) {
	token := "token"
	templateType := "email"
	slug := "metalslug"
	expectedResponse := dummyTemplateJSON

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/templates/%s/%s", templateType, slug)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "GET")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want TemplateExtended
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.Templates().Read(templateType, slug)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestTemplatesService_Upsert_happyPath(t *testing.T) {
	token := "token"
	templateType := "email"
	slug := "metalslug"
	expectedResponse := dummyTemplateJSON

	var payload UpsertTemplateRequest
	_ = json.Unmarshal([]byte(dummyUpsertRequestJSON), &payload)

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/templates/%s/%s", templateType, slug)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "PUT")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Templates().Upsert(templateType, slug, &payload)
	assert.Nil(t, err)

	var want TemplateExtended
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, payload)
	}
}

func TestTemplatesService_RevertToSystemTemplate_happyPath(t *testing.T) {
	token := "token"
	templateType := "email"
	slug := "metalslug"
	expectedResponse := dummyTemplateJSON

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/templates/%s/%s/revert", templateType, slug)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	var want TemplateExtended
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	got, err := client.Templates().Revert(templateType, slug)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestTemplatesService_Delete_happyPath(t *testing.T) {
	token := "token"
	templateType := "email"
	slug := "metalslug"

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/templates/%s/%s", templateType, slug)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "DELETE")
		testHeader(t, req, "Authorization", "Bearer "+token)
		response := fmt.Sprintf(`{ "deleted": true, "slug": "%v", "object": "template" }`, slug)
		fmt.Fprint(w, response)
	})

	want := DeleteResponse{Slug: slug, Object: "template", Deleted: true}

	got, err := client.Templates().Delete(templateType, slug)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, want)
	}
}

func TestTemplatesService_Preview_happyPath(t *testing.T) {
	token := "token"
	templateType := "sms"
	slug := "snail"
	expectedResponse := dummyPreviewResponseJSON

	var payload PreviewTemplateRequest
	_ = json.Unmarshal([]byte(dummyPreviewRequestJSON), &payload)

	client, mux, _, teardown := setup(token)
	defer teardown()

	url := fmt.Sprintf("/templates/%s/%s/preview", templateType, slug)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, "POST")
		testHeader(t, req, "Authorization", "Bearer "+token)
		fmt.Fprint(w, expectedResponse)
	})

	got, err := client.Templates().Preview(templateType, slug, &payload)
	assert.Nil(t, err)

	var want TemplatePreview
	_ = json.Unmarshal([]byte(expectedResponse), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Response = %v, want %v", *got, payload)
	}
}

const dummyTemplateJSON = `{
    "object": "template",
    "slug": "derp",
    "resource_type": "user",
    "template_type": "email",
    "name": "Vin Diesel",
    "position": 0,
    "created_at": 1633541368454,
    "updated_at": 1633541368454,
    "subject": "Choo choo train",
    "markup": "<p>Hee Hee</p>",
    "body": "<p>Ho Ho</p>",
	"available_variables": [
        "michael",
        "jackson"
    ],
    "required_variables": [
        "michael"
    ]
}`

const dummyUpsertRequestJSON = `{
	"name": "Dominic Toretto",
	"subject": "NOS bottles for sale",
	"markup": "<p>Family</p>",
	"body": "<p>One quarter of a mile at a time<p>"
}`

const dummyPreviewRequestJSON = `{
	"body": "{{OTPCode}} is your code for {{AppName}}, valid for {{TTLMinutes}} minutes"
}`

const dummyPreviewResponseJSON = `{
	"body": "123456 is your code for ACME, valid for 10 minutes"
}`
