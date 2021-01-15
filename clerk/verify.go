package clerk

type verifyRequest struct {
	Token string `json:"token"`
}

func verify(client Client, url string, token string, response interface{}) error {
	tokenPayload := verifyRequest{Token: token}
	req, _ := client.NewRequest("POST", url, &tokenPayload)

	_, err := client.Do(req, response)
	if err != nil {
		return err
	}
	return nil
}
