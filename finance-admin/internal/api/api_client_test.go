package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func getContext(cookies []*http.Cookie) Context {
	return Context{
		Context:   context.Background(),
		Cookies:   cookies,
		XSRFToken: "abcde",
	}
}

func TestClientError(t *testing.T) {
	assert.Equal(t, "message", ClientError("message").Error())
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, err, err.Data())
}