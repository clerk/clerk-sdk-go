package template

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestTemplateClientList(t *testing.T) {
	t.Parallel()
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T: t,
			Out: json.RawMessage(`{
"data": [{"template_type":"email","name":"the-name"}],
"total_count": 1
}`),
			Method: http.MethodGet,
			Path:   "/v1/templates/email",
			Query: &url.Values{
				"paginated": []string{"true"},
			},
		},
	}
	client := NewClient(config)
	list, err := client.List(context.Background(), &ListParams{
		TemplateType: "email",
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), list.TotalCount)
	require.Equal(t, 1, len(list.Templates))
	require.Equal(t, clerk.TemplateType("email"), list.Templates[0].TemplateType)
	require.Equal(t, "the-name", list.Templates[0].Name)
}

func TestTemplateClientGet(t *testing.T) {
	t.Parallel()
	templateType := clerk.TemplateTypeSMS
	slug := "the-slug"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"template_type":"%s","slug":"%s"}`, templateType, slug)),
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/v1/templates/%s/%s", templateType, slug),
		},
	}
	client := NewClient(config)
	template, err := client.Get(context.Background(), &GetParams{
		TemplateType: templateType,
		Slug:         slug,
	})
	require.NoError(t, err)
	require.Equal(t, slug, template.Slug)
	require.Equal(t, templateType, template.TemplateType)
}

func TestTemplateClientUpdate(t *testing.T) {
	t.Parallel()
	templateType := clerk.TemplateTypeEmail
	subject := "subject"
	slug := "the-slug"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"subject":"%s"}`, subject)),
			Out:    json.RawMessage(fmt.Sprintf(`{"template_type":"%s","subject":"%s","slug":"%s"}`, templateType, subject, slug)),
			Method: http.MethodPut,
			Path:   fmt.Sprintf("/v1/templates/%s/%s", templateType, slug),
		},
	}
	client := NewClient(config)
	template, err := client.Update(context.Background(), &UpdateParams{
		TemplateType: templateType,
		Slug:         slug,
		Subject:      clerk.String(subject),
	})
	require.NoError(t, err)
	require.Equal(t, slug, template.Slug)
	require.Equal(t, subject, template.Subject)
	require.Equal(t, templateType, template.TemplateType)
}

func TestTemplateClientDelete(t *testing.T) {
	t.Parallel()
	templateType := clerk.TemplateTypeEmail
	slug := "the-slug"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"slug":"%s","deleted":true}`, slug)),
			Method: http.MethodDelete,
			Path:   fmt.Sprintf("/v1/templates/%s/%s", templateType, slug),
		},
	}
	client := NewClient(config)
	template, err := client.Delete(context.Background(), &DeleteParams{
		TemplateType: templateType,
		Slug:         slug,
	})
	require.NoError(t, err)
	require.Equal(t, slug, template.Slug)
	require.True(t, template.Deleted)
}

func TestTemplateClientRevert(t *testing.T) {
	t.Parallel()
	templateType := clerk.TemplateTypeEmail
	slug := "the-slug"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"slug":"%s"}`, slug)),
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/v1/templates/%s/%s/revert", templateType, slug),
		},
	}
	client := NewClient(config)
	template, err := client.Revert(context.Background(), &RevertParams{
		TemplateType: templateType,
		Slug:         slug,
	})
	require.NoError(t, err)
	require.Equal(t, slug, template.Slug)
}

func TestTemplateClientPreview(t *testing.T) {
	t.Parallel()
	templateType := clerk.TemplateTypeEmail
	slug := "the-slug"
	body := "body"
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"body":"%s"}`, body)),
			Out:    json.RawMessage(fmt.Sprintf(`{"body":"%s"}`, body)),
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/v1/templates/%s/%s/preview", templateType, slug),
		},
	}
	client := NewClient(config)
	preview, err := client.Preview(context.Background(), &PreviewParams{
		TemplateType: templateType,
		Slug:         slug,
		Body:         clerk.String(body),
	})
	require.NoError(t, err)
	require.Equal(t, body, preview.Body)
}

func TestTemplateClientToggleDelivery(t *testing.T) {
	t.Parallel()
	templateType := clerk.TemplateTypeEmail
	slug := "the-slug"
	deliveredByClerk := true
	config := &ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"delivered_by_clerk":%v}`, deliveredByClerk)),
			Out:    json.RawMessage(fmt.Sprintf(`{"delivered_by_clerk":%v}`, deliveredByClerk)),
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/v1/templates/%s/%s/toggle_delivery", templateType, slug),
		},
	}
	client := NewClient(config)
	template, err := client.ToggleDelivery(context.Background(), &ToggleDeliveryParams{
		TemplateType:     templateType,
		Slug:             slug,
		DeliveredByClerk: clerk.Bool(deliveredByClerk),
	})
	require.NoError(t, err)
	require.Equal(t, deliveredByClerk, template.DeliveredByClerk)
}
