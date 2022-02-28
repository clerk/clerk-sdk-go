//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func TestJWTTemplates(t *testing.T) {
	c := createClient()

	jwtTemplates, err := c.JWTTemplates().ListAll()
	if err != nil {
		t.Fatalf("JWTTemplates.ListAll returned error: %v", err)
	}
	if jwtTemplates == nil {
		t.Fatalf("JWTTemplates.ListAll returned nil")
	}

	for _, jwtTemplate := range jwtTemplates {
		tmpl, err := c.JWTTemplates().Read(jwtTemplate.ID)
		if err != nil {
			t.Fatalf("JWTTemplates.Read returned error: %v", err)
		}
		if tmpl == nil {
			t.Fatalf("JWTTemplates.Read returned nil")
		}
	}

	newJWTTemplate := &clerk.CreateUpdateJWTTemplate{
		Name: fmt.Sprintf("Integration-%d", time.Now().Unix()),
		Claims: map[string]interface{}{
			"name": "{{user.first_name}}",
			"role": "tester",
		},
	}

	jwtTemplate, err := c.JWTTemplates().Create(newJWTTemplate)
	if err != nil {
		t.Fatalf("JWTTemplates.Create returned error: %v", err)
	}

	assert.Equal(t, newJWTTemplate.Name, jwtTemplate.Name)
	assert.Equal(t, 60, jwtTemplate.Lifetime)
	assert.Equal(t, 5, jwtTemplate.AllowedClockSkew)

	updateJWTTemplate := &clerk.CreateUpdateJWTTemplate{
		Name: fmt.Sprintf("Updated-Integration-%d", time.Now().Unix()),
		Claims: map[string]interface{}{
			"name": "{{user.first_name}}",
			"age":  28,
		},
	}

	updated, err := c.JWTTemplates().Update(jwtTemplate.ID, updateJWTTemplate)
	if err != nil {
		t.Fatalf("JWTTemplates.Create returned error: %v", err)
	}

	assert.Equal(t, jwtTemplate.ID, updated.ID)
	assert.Equal(t, updateJWTTemplate.Name, updated.Name)

	expectedClaimsBytes, _ := json.Marshal(updateJWTTemplate.Claims)
	assert.JSONEq(t, string(expectedClaimsBytes), string(updated.Claims))

	_, err = c.JWTTemplates().Delete(jwtTemplate.ID)
	if err != nil {
		t.Fatalf("JWTTemplates.Delete returned error: %v", err)
	}
}
