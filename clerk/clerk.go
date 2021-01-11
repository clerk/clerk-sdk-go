package clerk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	clerkBaseUrl = "https://api.clerk.dev/v1/"
)

type Client interface {
	NewRequest(method string, url string) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)

	Users() *UsersService
}

type service struct {
	client Client
}

type client struct {
	client  *http.Client
	baseURL *url.URL
	token   string

	users *UsersService
}

// NewClient creates a new Clerk client.
// Because the token supplied will be used for all authenticated requests,
// the created client should not be used across different users
func NewClient(token string) (Client, error) {
	return NewClientWithBaseUrl(token, clerkBaseUrl)
}

func NewClientWithBaseUrl(token string, baseUrl string) (Client, error) {
	baseURL, _ := url.Parse(baseUrl)
	httpClient := http.Client{}

	client := &client{client: &httpClient, baseURL: baseURL, token: token}

	commonService := &service{client: client}
	client.users = (*UsersService)(commonService)

	return client, nil
}

// NewRequest creates an API request.
// A relative URL can be specified which is resolved relative to the BaseURL of the client.
// Relative URLs should be specified without a preceding slash.
func (c *client) NewRequest(method string, url string) (*http.Request, error) {
	fullUrl, err := c.baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fullUrl.String(), nil)
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

func (c *client) Users() *UsersService {
	return c.users
}
