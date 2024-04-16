//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func TestTemplates(t *testing.T) {
	client := createClient()
	templateType := "email"

	// Get all email templates
	templates, err := client.Templates().ListAll(templateType)
	if err != nil {
		t.Fatalf("Templates.ListAll returned error: %v", err)
	}
	if templates == nil {
		t.Fatalf("Templates.ListAll returned nil")
	}

	for _, template := range templates {
		slug := template.Slug

		// Make sure we can read each template
		tmpl, err := client.Templates().Read(templateType, slug)
		if err != nil {
			t.Fatalf("Templates.Read returned error: %v", err)
		}
		if tmpl == nil {
			t.Fatalf("Templates.Read returned nil")
		}

		// Preview each template with sample data
		fromEmailName := "marketing"
		replyToEmailName := "support"
		templatePreview, err := client.Templates().Preview(templateType, slug, &clerk.PreviewTemplateRequest{
			Subject:          "{{AppName}} is da bomb",
			Body:             "<p><a href=\"{{AppURL}}\">{{AppName}}</a> is the greatest app of all time!</p>",
			FromEmailName:    &fromEmailName,
			ReplyToEmailName: &replyToEmailName,
		})
		if err != nil {
			t.Fatalf("Templates.Preview returned error: %v", err)
		}
		if templatePreview == nil {
			t.Errorf("Templates.Preview returned nil")
		}
	}
}

func TestTemplates_Upsert(t *testing.T) {
	client := createClient()

	// Update one of the templates, just to make sure that the Upsert method works.
	templateType := "email"
	slug := "organization_invitation"
	requiredVariable := "{{action_url}}"
	deliveredByClerk := false
	fromEmailName := "marketing"
	replyToEmailName := "support"
	upsertedTemplate, err := client.Templates().Upsert(templateType, slug, &clerk.UpsertTemplateRequest{
		Name:             "Remarketing email",
		Subject:          "Unmissable opportunity",
		Markup:           "",
		Body:             fmt.Sprintf("Click %s for free unicorns", requiredVariable),
		FromEmailName:    &fromEmailName,
		ReplyToEmailName: &replyToEmailName,
		DeliveredByClerk: &deliveredByClerk,
	})
	if err != nil {
		t.Fatalf("Templates.Update returned error: %v", err)
	}
	if upsertedTemplate == nil {
		t.Errorf("Templates.Upsert returned nil")
	}
}
