package clerk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	ProdUrl = "https://api.clerk.dev/v1/"

	ClientsUrl       = "clients"
	ClientsVerifyUrl = ClientsUrl + "/verify"
	EmailsUrl        = "emails"
	SessionsUrl      = "sessions"
	SMSUrl           = "sms_messages"
	UsersUrl         = "users"
)

type Client interface {
	NewRequest(method string, url string, body ...interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)

	Clients() *ClientsService
	Emails() *EmailService
	Sessions() *SessionsService
	SMS() *SMSService
	Users() *UsersService
	Verification() *VerificationService
}

type service struct {
	client Client
}

type client struct {
	client  *http.Client
	baseURL *url.URL
	token   string

	clients      *ClientsService
	emails       *EmailService
	sessions     *SessionsService
	sms          *SMSService
	users        *UsersService
	verification *VerificationService
}

// NewClient creates a new Clerk client.
// Because the token supplied will be used for all authenticated requests,
// the created client should not be used across different users
func NewClient(token string) (Client, error) {
	return NewClientWithBaseUrl(token, ProdUrl)
}

func NewClientWithBaseUrl(token string, baseUrl string) (Client, error) {
	baseURL, _ := url.Parse(baseUrl)
	httpClient := http.Client{}

	client := &client{client: &httpClient, baseURL: baseURL, token: token}

	commonService := &service{client: client}
	client.clients = (*ClientsService)(commonService)
	client.emails = (*EmailService)(commonService)
	client.sessions = (*SessionsService)(commonService)
	client.sms = (*SMSService)(commonService)
	client.users = (*UsersService)(commonService)
	client.verification = (*VerificationService)(commonService)

	return client, nil
}

// NewRequestWithBody creates an API request.
// A relative URL `url` can be specified which is resolved relative to the baseURL of the client.
// Relative URLs should be specified without a preceding slash.
// The `body` parameter can be used to pass a body to the request. If no body is required, the parameter can be omitted.
func (c *client) NewRequest(method string, url string, body ...interface{}) (*http.Request, error) {
	fullUrl, err := c.baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if len(body) > 0 && body[0] != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body[0])
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, fullUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do will send the given request using the client `c` on which it is called.
// If the response contains a body, it will be unmarshalled in `v`.
func (c *client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkForErrors(resp)
	if err != nil {
		return resp, err
	}

	if resp.Body != nil && v != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		err = json.Unmarshal(body, &v)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func checkForErrors(resp *http.Response) error {
	if c := resp.StatusCode; c >= 200 && c < 400 {
		return nil
	}

	data, _ := ioutil.ReadAll(resp.Body)
	if data != nil && len(data) > 0 {
		return errors.New(string(data))
	}
	return errors.New(fmt.Sprintf("Server returned unexpected error with status code %d", resp.StatusCode))
}

func (c *client) Clients() *ClientsService {
	return c.clients
}

func (c *client) Emails() *EmailService {
	return c.emails
}

func (c *client) Sessions() *SessionsService {
	return c.sessions
}

func (c *client) SMS() *SMSService {
	return c.sms
}

func (c *client) Users() *UsersService {
	return c.users
}

func (c *client) Verification() *VerificationService {
	return c.verification
}
