package clerk

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomStringCommon(t *testing.T) {
	value := String{}
	assert.Equal(t, "", value.Value)
	assert.Equal(t, false, value.IsNull)
	assert.Equal(t, false, value.IsSet)

	value = NewString("foo")
	assert.Equal(t, "foo", value.Value)
	assert.Equal(t, false, value.IsNull)
	assert.Equal(t, true, value.IsSet)

	value = NewNullString()
	assert.Equal(t, "", value.Value)
	assert.Equal(t, true, value.IsNull)
	assert.Equal(t, true, value.IsSet)
}

func TestCustomStringJsonDecode(t *testing.T) {
	jsonPayload := `{
	"first_name": "foo",
	"middle_name": null
}`
	payload := struct {
		FirstName  String `json:"first_name"`
		MiddleName String `json:"middle_name"`
		LastName   String `json:"last_name"`
	}{}
	err := json.NewDecoder(strings.NewReader(jsonPayload)).Decode(&payload)
	if err != nil {
		t.Fatalf("failed to decode json: %v\n%v", err, jsonPayload)
	}

	assert.Equal(t, "foo", payload.FirstName.Value)
	assert.Equal(t, false, payload.FirstName.IsNull)
	assert.Equal(t, true, payload.FirstName.IsSet)

	assert.Equal(t, "", payload.MiddleName.Value)
	assert.Equal(t, true, payload.MiddleName.IsNull)
	assert.Equal(t, true, payload.MiddleName.IsSet)

	assert.Equal(t, "", payload.LastName.Value)
	assert.Equal(t, false, payload.LastName.IsNull)
	assert.Equal(t, false, payload.LastName.IsSet)
}

func TestCustomStringJsonEncode(t *testing.T) {
	firstName := NewString("foo")
	middleName := NewNullString()
	payload := struct {
		FirstName  *String `json:"first_name,omitempty"`
		MiddleName *String `json:"middle_name,omitempty"`
		LastName   *String `json:"last_name,omitempty"`
	}{
		FirstName:  &firstName,
		MiddleName: &middleName,
		LastName:   nil,
	}

	writer := bytes.NewBufferString("")
	err := json.NewEncoder(writer).Encode(payload)
	if err != nil {
		t.Fatalf("failed to encode to json: %v", err)
	}

	expected := `{
	"first_name": "foo",
	"middle_name": null
}`
	assert.JSONEq(t, expected, writer.String())
}
