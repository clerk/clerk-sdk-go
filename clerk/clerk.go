package clerk

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	ProdUrl = "https://api.clerk.dev/v1/"

	ClientsUrl       = "clients"
	ClientsVerifyUrl = ClientsUrl + "/verify"
	EmailsUrl        = "emails"
	SessionsUrl      = "sessions"
	SMSUrl           = "sms_messages"
	UsersUrl         = "users"
	WebhooksUrl      = "webhooks"
)

type Client interface {
	NewRequest(method string, url string, body ...interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)

	Clients() *ClientsService
	Emails() *EmailService
	Sessions() *SessionsService
	SMS() *SMSService
	Users() *UsersService
	Webhooks() *WebhooksService
	Verification() *VerificationService
	Tokens() *TokensService
}

type service struct {
	client Client
}

type client struct {
	client  *http.Client
	baseURL *url.URL
	token   string
	jwks    []*jwk

	clients      *ClientsService
	emails       *EmailService
	sessions     *SessionsService
	sms          *SMSService
	users        *UsersService
	webhooks     *WebhooksService
	verification *VerificationService
	tokens       *TokensService
}

type jwk struct {
	// Sig (for signature) or Enc (for encryption)
	publicKeyUse string `json:"use"`

	// Algorithm family (RSA, ECDSA etc.)
	keyType string `json:"kty"`

	// RSA256
	Algorithm string `json:"alg"`

	// Clerk instance ID
	KeyID string `json:"kid"`

	n string `json:"n"`
	e string `json:"e"`
}

// NewClient creates a new Clerk client.
// Because the token supplied will be used for all authenticated requests,
// the created client should not be used across different users
func NewClient(token string) (Client, error) {
	return NewClientWithBaseUrl(token, ProdUrl)
}

func NewClientWithBaseUrl(token string, baseUrl string) (Client, error) {
	httpClient := http.Client{}

	return NewClientWithCustomHTTP(token, baseUrl, &httpClient)
}

func NewClientWithCustomHTTP(token string, urlStr string, httpClient *http.Client) (Client, error) {
	baseURL := toURLWithEndingSlash(urlStr)
	client := &client{client: httpClient, baseURL: baseURL, token: token}

	commonService := &service{client: client}
	client.clients = (*ClientsService)(commonService)
	client.emails = (*EmailService)(commonService)
	client.sessions = (*SessionsService)(commonService)
	client.sms = (*SMSService)(commonService)
	client.users = (*UsersService)(commonService)
	client.webhooks = (*WebhooksService)(commonService)
	client.verification = (*VerificationService)(commonService)
	client.tokens = (*TokensService)(commonService)

	return client, nil
}

func toURLWithEndingSlash(u string) *url.URL {
	baseURL, _ := url.Parse(u)
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	return baseURL
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

	errorResponse := &ErrorResponse{Response: resp}

	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		// it's ok if we cannot unmarshal to Clerk's error response
		_ = json.Unmarshal(data, errorResponse)
	}

	return errorResponse
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

func (c *client) Webhooks() *WebhooksService {
	return c.webhooks
}

func (c *client) Verification() *VerificationService {
	return c.verification
}

func (c *client) Tokens() *TokensService {
	return c.tokens
}
