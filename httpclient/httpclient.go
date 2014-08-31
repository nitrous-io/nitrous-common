package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	*http.Client
	AuthToken string
}

var DefaultClient = &Client{http.DefaultClient, ""}

var ErrAuthTokenEmpty = errors.New("Auth token is not provided")

type ErrResponse struct {
	Hash        string `json:"error_hash"`
	Description string `json:"error_description"`
}

func (e *ErrResponse) Error() string {
	switch {
	case e.Description != "":
		return e.Description
	case e.Hash != "":
		return fmt.Sprintf("something went wrong (contact support with this error code: %s)", e.Hash)
	default:
		return "unknown error occurred"
	}
}

// New returns a client that uses net/http's DefaultClient with the token given by authToken.
// If AuthToken is an empty string, *Authed() calls will fail with a "must be logged in" error message
func New(authToken string) *Client {
	return &Client{http.DefaultClient, authToken}
}

// CheckResponse retrieves the error_hash and error_description on error responses
func CheckResponse(resp *http.Response, expectedCode int) error {
	if resp.StatusCode == expectedCode {
		return nil
	}
	var j ErrResponse
	err := json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		return err
	}

	return &j
}

func (c *Client) Get(urlStr string, params *url.Values) (*http.Response, error) {
	return c.DoRequest("GET", urlStr, params)
}

func (c *Client) Delete(urlStr string, params *url.Values) (*http.Response, error) {
	return c.DoRequest("DELETE", urlStr, params)
}

// http.Client already includes PostForm, but this is here for consistency
func (c *Client) PostForm(urlStr string, data *url.Values) (*http.Response, error) {
	return c.DoFormRequest("POST", urlStr, data)
}

func (c *Client) PutForm(urlStr string, data *url.Values) (*http.Response, error) {
	return c.DoFormRequest("PUT", urlStr, data)
}

func (c *Client) GetAuthed(urlStr string, params *url.Values) (*http.Response, error) {
	return c.DoRequestAuthed("GET", urlStr, params)
}

func (c *Client) DeleteAuthed(urlStr string, params *url.Values) (*http.Response, error) {
	return c.DoRequestAuthed("DELETE", urlStr, params)
}

func (c *Client) PostFormAuthed(urlStr string, data *url.Values) (*http.Response, error) {
	return c.DoFormRequestAuthed("POST", urlStr, data)
}

func (c *Client) PutFormAuthed(urlStr string, data *url.Values) (*http.Response, error) {
	return c.DoFormRequestAuthed("PUT", urlStr, data)
}

func (c *Client) DoRequest(method string, urlStr string, params *url.Values) (*http.Response, error) {
	req, err := c.PrepareRequest(method, urlStr, params)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) DoFormRequest(method string, urlStr string, data *url.Values) (*http.Response, error) {
	req, err := c.PrepareFormRequest(method, urlStr, data)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) DoRequestAuthed(method string, urlStr string, params *url.Values) (*http.Response, error) {
	req, err := c.PrepareRequest(method, urlStr, params)
	if err != nil {
		return nil, err
	}
	return c.DoAuthed(req)
}

func (c *Client) DoFormRequestAuthed(method string, urlStr string, data *url.Values) (*http.Response, error) {
	req, err := c.PrepareFormRequest(method, urlStr, data)
	if err != nil {
		return nil, err
	}
	return c.DoAuthed(req)
}

func (c *Client) PrepareRequest(method string, urlStr string, params *url.Values) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) PrepareFormRequest(method string, urlStr string, data *url.Values) (*http.Request, error) {
	var body *strings.Reader
	if data == nil {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(data.Encode())
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c *Client) DoAuthed(req *http.Request) (*http.Response, error) {
	if c.AuthToken == "" {
		return nil, ErrAuthTokenEmpty
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	return c.Do(req)
}
