package api

import (
	"context"
	"fmt"
	"net/http"
)

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
	ClientId  int
}

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

func NewApiClient(httpClient HTTPClient) *ApiClient {
	return &ApiClient{
		http: httpClient,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiClient struct {
	http HTTPClient
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
