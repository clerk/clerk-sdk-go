package clerk

type OAuthApplication struct {
	APIResource
	Object                string   `json:"object"`
	ID                    string   `json:"id"`
	InstanceID            string   `json:"instance_id"`
	Name                  string   `json:"name"`
	ClientID              string   `json:"client_id"`
	ClientSecret          *string  `json:"client_secret,omitempty"`
	ClientImageURL        *string  `json:"client_image_url"`
	ClientURL             *string  `json:"client_url"`
	PKCERequired          bool     `json:"pkce_required"`
	Public                bool     `json:"public"`
	DynamicallyRegistered bool     `json:"dynamically_registered"`
	ConsentScreenEnabled  bool     `json:"consent_screen_enabled"`
	Scopes                string   `json:"scopes"`
	RedirectURIs          []string `json:"redirect_uris"`
	DiscoveryURL          string   `json:"discovery_url"`
	AuthorizeURL          string   `json:"authorize_url"`
	TokenFetchURL         string   `json:"token_fetch_url"`
	UserInfoURL           string   `json:"user_info_url"`
	TokenIntrospectionURL string   `json:"token_introspection_url"`
	CreatedAt             int64    `json:"created_at"`
	UpdatedAt             int64    `json:"updated_at"`

	// Deprecated: Use RedirectURIs instead
	CallbackURL string `json:"callback_url"`
}

type OAuthApplicationList struct {
	APIResource
	OAuthApplications []*OAuthApplication `json:"data"`
	TotalCount        int64               `json:"total_count"`
}
