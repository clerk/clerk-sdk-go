package clerk

type WebhooksService service

type DiahookResponse struct {
	DiahookURL string `json:"diahook_url"`
}

func (s *WebhooksService) CreateDiahook() (*DiahookResponse, error) {
	diahookUrl := WebhooksUrl + "/diahook"
	req, _ := s.client.NewRequest("POST", diahookUrl)

	var diahookResponse DiahookResponse
	if _, err := s.client.Do(req, &diahookResponse); err != nil {
		return nil, err
	}
	return &diahookResponse, nil
}

func (s *WebhooksService) DeleteDiahook() error {
	diahookUrl := WebhooksUrl + "/diahook"
	req, _ := s.client.NewRequest("DELETE", diahookUrl)

	_, err := s.client.Do(req, nil)

	return err
}

func (s *WebhooksService) RefreshDiahookURL() (*DiahookResponse, error) {
	diahookUrl := WebhooksUrl + "/diahook_url"
	req, _ := s.client.NewRequest("POST", diahookUrl)

	var diahookResponse DiahookResponse
	if _, err := s.client.Do(req, &diahookResponse); err != nil {
		return nil, err
	}
	return &diahookResponse, nil
}
