package clerk

// Clerk supports different types of templates.
type TemplateType string

// List of supported values for template types.
const (
	TemplateTypeEmail TemplateType = "email"
	TemplateTypeSMS   TemplateType = "sms"
)

type Template struct {
	APIResource
	Object             string       `json:"object"`
	Slug               string       `json:"slug"`
	ResourceType       string       `json:"resource_type"`
	TemplateType       TemplateType `json:"template_type"`
	Name               string       `json:"name"`
	Position           int          `json:"position"`
	CanRevert          bool         `json:"can_revert"`
	CanDelete          bool         `json:"can_delete"`
	FromEmailName      *string      `json:"from_email_name,omitempty"`
	ReplyToEmailName   *string      `json:"reply_to_email_name,omitempty"`
	DeliveredByClerk   bool         `json:"delivered_by_clerk"`
	Subject            string       `json:"subject"`
	Markup             string       `json:"markup"`
	Body               string       `json:"body"`
	AvailableVariables []string     `json:"available_variables"`
	RequiredVariables  []string     `json:"required_variables"`
	CreatedAt          int64        `json:"created_at"`
	UpdatedAt          int64        `json:"updated_at"`
}

type TemplateList struct {
	APIResource
	Templates  []*Template `json:"data"`
	TotalCount int64       `json:"total_count"`
}

type TemplatePreview struct {
	APIResource
	Subject             string  `json:"subject,omitempty"`
	Body                string  `json:"body"`
	FromEmailAddress    *string `json:"from_email_address,omitempty"`
	ReplyToEmailAddress *string `json:"reply_to_email_address,omitempty"`
}
