package clerk

type SvixWebhook struct {
	APIResource
	SvixURL string `json:"svix_url"`
}
