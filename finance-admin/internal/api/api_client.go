package api

import (
	"context"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"io"
	"net/http"
)

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type JWTClient interface {
	CreateJWT(ctx context.Context) string
}

type EnvVars struct {
	SiriusURL string
	HubURL    string
}

type Client struct {
	http HTTPClient
	jwt  JWTClient
	EnvVars
}

func NewClient(httpClient HTTPClient, jwtClient JWTClient, env EnvVars) *Client {
	return &Client{
		http:    httpClient,
		jwt:     jwtClient,
		EnvVars: env,
	}
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

func (c *Client) newHubRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.HubURL+path, body)
	if err != nil {
		return nil, err
	}

	addCookiesFromContext(ctx, req)
	req.Header.Add("Authorization", "Bearer "+c.jwt.CreateJWT(ctx))

	return req, err
}

func (c *Client) newSessionRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.SiriusURL+"/supervision-api/v1/users/current", nil)
	if err != nil {
		return nil, err
	}

	addCookiesFromContext(ctx, req)
	req.Header.Add("OPG-Bypass-Membrane", "1")

	return req, err
}

func addCookiesFromContext(ctx context.Context, req *http.Request) {
	for _, c := range ctx.(auth.Context).Cookies {
		req.AddCookie(c)
	}
}

func newStatusError(resp *http.Response) StatusError {
	return StatusError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL.String(),
		Method: resp.Request.Method,
	}
}
