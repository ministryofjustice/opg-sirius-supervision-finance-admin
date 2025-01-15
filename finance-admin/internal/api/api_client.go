package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
}

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

func NewClient(httpClient HTTPClient, siriusURL string, backendURL string, hubURL string) (*Client, error) {
	return &Client{
		http:       httpClient,
		SiriusURL:  siriusURL,
		BackendURL: backendURL,
		HubURL:     hubURL,
	}, nil
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	http       HTTPClient
	SiriusURL  string
	BackendURL string
	HubURL     string
}

type StatusError struct {
	Code   int    `json:"code"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.URL, e.Code)
}

func (e StatusError) Data() interface{} {
	return e
}

// Deprecated: newBackendRequest will be removed once backend is migrated to the Hub
func (c *Client) newBackendRequest(ctx Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx.Context, method, c.BackendURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("X-XSRF-TOKEN", ctx.XSRFToken)

	return req, err
}

func (c *Client) newHubRequest(ctx Context, method, path string, body io.Reader) (*http.Request, error) {
	fmt.Println("Making hub request: " + c.HubURL + path)
	req, err := http.NewRequestWithContext(ctx.Context, method, c.HubURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("X-XSRF-TOKEN", ctx.XSRFToken)

	return req, err
}

func (c *Client) newSessionRequest(ctx Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx.Context, "GET", c.SiriusURL+"/supervision-api/v1/users/current", nil)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("OPG-Bypass-Membrane", "1")

	return req, err
}

func newStatusError(resp *http.Response) StatusError {
	return StatusError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL.String(),
		Method: resp.Request.Method,
	}
}
