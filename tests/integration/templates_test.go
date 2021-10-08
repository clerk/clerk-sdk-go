// +build integration

package integration

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"testing"
)

func TestTemplates(t *testing.T) {
	client := createClient()

	templateType := "email"

	templates, err := client.Templates().ListAll(templateType)
	if err != nil {
		t.Fatalf("Templates.ListAll returned error: %v", err)
	}
	if templates == nil {
		t.Fatalf("Templates.ListAll returned nil")
	}

	for _, template := range templates {
		slug := template.Slug

		tmpl, err := client.Templates().Read(templateType, slug)
		if err != nil {
			t.Fatalf("Templates.Read returned error: %v", err)
		}
		if tmpl == nil {
			t.Fatalf("Templates.Read returned nil")
		}

		upsertTemplateRequest := clerk.UpsertTemplateRequest{
			Name:               "Remarketing SMS",
			Markup:             "",
			Body:               "Click {link} for free unicorns",
			MandatoryVariables: []string{"lorem", "ipsum"},
		}

		upsertedTemplate, err := client.Templates().Upsert(templateType, slug, &upsertTemplateRequest)
		if err != nil {
			t.Fatalf("Templates.Update returned error: %v", err)
		}
		if upsertedTemplate == nil {
			t.Errorf("Templates.Upsert returned nil")
		}
	}
}
