package clerk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTTemplatesService_ListAll(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	dummyResponse := "[" + dummyJWTTemplateJSON + "]"

	mux.HandleFunc("/jwt_templates", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	got, err := c.JWTTemplates().ListAll()
	assert.Nil(t, err)

	expected := make([]JWTTemplate, 0)
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestJWTTemplatesService_Read(t *testing.T) {
	dummyResponse := dummyTemplateJSON

	c, mux, _, teardown := setup("token")
	defer teardown()

	url := fmt.Sprintf("/jwt_templates/%s", dummyJWTTemplateID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodGet)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	got, err := c.JWTTemplates().Read(dummyJWTTemplateID)
	assert.Nil(t, err)

	expected := JWTTemplate{}
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestJWTTemplatesService_Create(t *testing.T) {
	dummyResponse := dummyJWTTemplateJSON

	c, mux, _, teardown := setup("token")
	defer teardown()

	mux.HandleFunc("/jwt_templates", func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPost)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	newJWTTmpl := &CreateUpdateJWTTemplate{
		Name: "Testing",
		Claims: map[string]interface{}{
			"name": "{{user.first_name}}",
			"role": "tester",
		},
	}

	got, err := c.JWTTemplates().Create(newJWTTmpl)
	assert.Nil(t, err)

	expected := JWTTemplate{}
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestJWTTemplatesService_Update(t *testing.T) {
	dummyResponse := dummyJWTTemplateUpdateJSON

	c, mux, _, teardown := setup("token")
	defer teardown()

	url := fmt.Sprintf("/jwt_templates/%s", dummyJWTTemplateID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodPatch)
		testHeader(t, req, "Authorization", "Bearer token")
		_, _ = fmt.Fprint(w, dummyResponse)
	})

	updateJWTTmpl := &CreateUpdateJWTTemplate{
		Name: "New-Testing",
		Claims: map[string]interface{}{
			"name": "{{user.first_name}}",
			"age":  28,
		},
	}

	got, err := c.JWTTemplates().Update(dummyJWTTemplateID, updateJWTTmpl)
	assert.Nil(t, err)

	expected := JWTTemplate{}
	_ = json.Unmarshal([]byte(dummyResponse), &expected)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", got, expected)
	}
}

func TestJWTTemplatesService_Delete(t *testing.T) {
	c, mux, _, teardown := setup("token")
	defer teardown()

	url := fmt.Sprintf("/jwt_templates/%s", dummyJWTTemplateID)

	mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		testHttpMethod(t, req, http.MethodDelete)
		testHeader(t, req, "Authorization", "Bearer token")
		response := fmt.Sprintf(`{ "deleted": true, "id": "%s", "object": "jwt_template" }`, dummyJWTTemplateID)
		_, _ = fmt.Fprint(w, response)
	})

	expected := DeleteResponse{
		ID:      dummyJWTTemplateID,
		Object:  "jwt_template",
		Deleted: true,
	}

	got, err := c.JWTTemplates().Delete(dummyJWTTemplateID)
	assert.Nil(t, err)

	if !reflect.DeepEqual(*got, expected) {
		t.Errorf("Response = %v, want %v", *got, expected)
	}
}

const (
	dummyJWTTemplateID = "jtmp_21xC2Ziqscwjq43MtC3CN6Pngbo"

	dummyJWTTemplateJSON = `
{
    "object": "jwt_template",
	"id": "` + dummyJWTTemplateID + `",
    "name": "Testing",
    "claims": {
		"name": "{{user.first_name}}",
		"role": "tester"
	},
	"lifetime": 60,
	"allowed_clock_skew": 5
}`

	dummyJWTTemplateUpdateJSON = `
{
    "object": "jwt_template",
	"id": "` + dummyJWTTemplateID + `",
    "name": "New-Testing",
    "claims": {
		"name": "{{user.first_name}}",
		"age": 28
	},
	"lifetime": 60,
	"allowed_clock_skew": 5
}`
)
