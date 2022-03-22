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

		var requiredVariable string
		switch slug {
		case "invitation", "organization_invitation":
			requiredVariable = "{{ActionURL}}"
		case "magic_link":
			requiredVariable = "{{MagicLink}}"
		case "suspicious_activity":
			requiredVariable = "{{Reason}}"
		case "verification_code":
			requiredVariable = "{{OTPCode}}"
		}

		upsertTemplateRequest := clerk.UpsertTemplateRequest{
			Name:    "Remarketing email",
			Subject: "Unmissable opportunity",
			Markup:  "",
			Body:    fmt.Sprintf("Click %s for free unicorns", requiredVariable),
		}

		upsertedTemplate, err := client.Templates().Upsert(templateType, slug, &upsertTemplateRequest)
		if err != nil {
			t.Fatalf("Templates.Update returned error: %v", err)
		}
		if upsertedTemplate == nil {
			t.Errorf("Templates.Upsert returned nil")
		}

		previewTemplateRequest := clerk.PreviewTemplateRequest{
			Subject: "{{AppName}} is da bomb",
			Body:    "<p><a href=\"{{AppURL}}\">{{AppName}}</a> is the greatest app of all time!</p>",
		}

		templatePreview, err := client.Templates().Preview(templateType, slug, &previewTemplateRequest)
		if err != nil {
			t.Fatalf("Templates.Preview returned error: %v", err)
		}
		if templatePreview == nil {
			t.Errorf("Templates.Preview returned nil")
		}
	}
}
