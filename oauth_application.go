package clerk

type OAuthApplication struct {
	APIResource
	Object                string  `json:"object"`
	ID                    string  `json:"id"`
	InstanceID            string  `json:"instance_id"`
	Name                  string  `json:"name"`
	ClientID              string  `json:"client_id"`
	Public                bool    `json:"public"`
	DynamicallyRegistered bool    `json:"dynamically_registered"`
	ConsentScreenEnabled  bool    `json:"consent_screen_enabled"`
	Scopes                string  `json:"scopes"`
	CallbackURL           string  `json:"callback_url"`
	DiscoveryURL          string  `json:"discovery_url"`
	AuthorizeURL          string  `json:"authorize_url"`
	TokenFetchURL         string  `json:"token_fetch_url"`
	UserInfoURL           string  `json:"user_info_url"`
	CreatedAt             int64   `json:"created_at"`
	UpdatedAt             int64   `json:"updated_at"`
	ClientSecret          *string `json:"client_secret,omitempty"`
}

type OAuthApplicationList struct {
	APIResource
	OAuthApplications []*OAuthApplication `json:"data"`
	TotalCount        int64               `json:"total_count"`
}
