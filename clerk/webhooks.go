package clerk

type WebhooksService service

type SvixResponse struct {
	SvixURL string `json:"svix_url"`
}

func (s *WebhooksService) CreateSvix() (*SvixResponse, error) {
	svixUrl := WebhooksUrl + "/svix"
	req, _ := s.client.NewRequest("POST", svixUrl)

	var svixResponse SvixResponse
	if _, err := s.client.Do(req, &svixResponse); err != nil {
		return nil, err
	}
	return &svixResponse, nil
}

func (s *WebhooksService) DeleteSvix() error {
	svixUrl := WebhooksUrl + "/svix"
	req, _ := s.client.NewRequest("DELETE", svixUrl)

	_, err := s.client.Do(req, nil)

	return err
}

func (s *WebhooksService) RefreshSvixURL() (*SvixResponse, error) {
	svixUrl := WebhooksUrl + "/svix_url"
	req, _ := s.client.NewRequest("POST", svixUrl)

	var svixResponse SvixResponse
	if _, err := s.client.Do(req, &svixResponse); err != nil {
		return nil, err
	}
	return &svixResponse, nil
}
