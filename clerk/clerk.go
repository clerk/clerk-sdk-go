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
	"strings"
	"time"
)

const version = "1.5.0"

const (
	ProdUrl = "https://api.clerk.dev/v1/"

	ClientsUrl       = "clients"
	ClientsVerifyUrl = ClientsUrl + "/verify"
	EmailsUrl        = "emails"
	OrganizationsUrl = "organizations"
	SessionsUrl      = "sessions"
	SMSUrl           = "sms_messages"
	TemplatesUrl     = "templates"
	UsersUrl         = "users"
	UsersCountUrl    = UsersUrl + "/count"
	WebhooksUrl      = "webhooks"
	JWTTemplatesUrl  = "jwt_templates"
)

var defaultHTTPClient = &http.Client{Timeout: time.Second * 5}

type Client interface {
	NewRequest(method, url string, body ...interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)

	DecodeToken(token string) (*TokenClaims, error)
	VerifyToken(token string, opts ...VerifyTokenOption) (*SessionClaims, error)

	Clients() *ClientsService
	Emails() *EmailService
	JWKS() *JWKSService
	JWTTemplates() *JWTTemplatesService
	Organizations() *OrganizationsService
	Sessions() *SessionsService
	SMS() *SMSService
	Templates() *TemplatesService
	Users() *UsersService
	Webhooks() *WebhooksService
	Verification() *VerificationService
	Interstitial() ([]byte, error)

	APIKey() string
}

type service struct {
	client Client
}

type client struct {
	client    *http.Client
	baseURL   *url.URL
	jwksCache *jwksCache
	token     string

	clients       *ClientsService
	emails        *EmailService
	jwks          *JWKSService
	jwtTemplates  *JWTTemplatesService
	organizations *OrganizationsService
	sessions      *SessionsService
	sms           *SMSService
	templates     *TemplatesService
	users         *UsersService
	webhooks      *WebhooksService
	verification  *VerificationService
}

// NewClient creates a new Clerk client.
// Because the token supplied will be used for all authenticated requests,
// the created client should not be used across different users
func NewClient(token string, options ...ClerkOption) (Client, error) {
	if token == "" {
		return nil, errors.New("you must provide an API token")
	}

	defaultBaseURL, err := toURLWithEndingSlash(ProdUrl)
	if err != nil {
		return nil, err
	}

	client := &client{
		client:  defaultHTTPClient,
		baseURL: defaultBaseURL,
		token:   token,
	}

	for _, option := range options {
		if err = option(client); err != nil {
			return nil, err
		}
	}

	commonService := &service{client: client}
	client.clients = (*ClientsService)(commonService)
	client.emails = (*EmailService)(commonService)
	client.jwks = (*JWKSService)(commonService)
	client.jwtTemplates = (*JWTTemplatesService)(commonService)
	client.organizations = (*OrganizationsService)(commonService)
	client.sessions = (*SessionsService)(commonService)
	client.sms = (*SMSService)(commonService)
	client.templates = (*TemplatesService)(commonService)
	client.users = (*UsersService)(commonService)
	client.webhooks = (*WebhooksService)(commonService)
	client.verification = (*VerificationService)(commonService)

	client.jwksCache = &jwksCache{}

	return client, nil
}

// Deprecated: NewClientWithBaseUrl is deprecated. Use the NewClient instead e.g. NewClient(token, WithBaseURL(baseUrl))
func NewClientWithBaseUrl(token, baseUrl string) (Client, error) {
	return NewClient(token, WithBaseURL(baseUrl))
}

// Deprecated: NewClientWithCustomHTTP is deprecated. Use the NewClient instead e.g. NewClient(token, WithBaseURL(urlStr), WithHTTPClient(httpClient))
func NewClientWithCustomHTTP(token, urlStr string, httpClient *http.Client) (Client, error) {
	return NewClient(token, WithBaseURL(urlStr), WithHTTPClient(httpClient))
}

func toURLWithEndingSlash(u string) (*url.URL, error) {
	baseURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	return baseURL, err
}

// NewRequest creates an API request.
// A relative URL `url` can be specified which is resolved relative to the baseURL of the client.
// Relative URLs should be specified without a preceding slash.
// The `body` parameter can be used to pass a body to the request. If no body is required, the parameter can be omitted.
func (c *client) NewRequest(method, url string, body ...interface{}) (*http.Request, error) {
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

	// Add custom header with the current SDK version
	req.Header.Set("X-Clerk-SDK", fmt.Sprintf("go/%s", version))

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

func (c *client) JWKS() *JWKSService {
	return c.jwks
}

func (c *client) JWTTemplates() *JWTTemplatesService {
	return c.jwtTemplates
}

func (c *client) Organizations() *OrganizationsService {
	return c.organizations
}

func (c *client) Sessions() *SessionsService {
	return c.sessions
}

func (c *client) SMS() *SMSService {
	return c.sms
}

func (c *client) Templates() *TemplatesService {
	return c.templates
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

func (c *client) APIKey() string {
	return c.token
}

func (c *client) Interstitial() ([]byte, error) {
	req, err := c.NewRequest("GET", "internal/interstitial")
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	interstitial, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return interstitial, err
	}

	return interstitial, nil
}
