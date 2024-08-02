package server

import (
	"io"
	"net/http"
)

type mockRoute struct { //nolint:golint,unused
	client   ApiClient
	data     any
	executed bool
	lastW    io.Writer
	error
}

func (r *mockRoute) Client() ApiClient { //nolint:golint,unused
	return r.client
}

func (r *mockRoute) execute(w http.ResponseWriter, req *http.Request, data any) error { //nolint:golint,unused
	r.executed = true
	r.lastW = w
	r.data = data
	return r.error
}

type mockApiClient struct {
	error error //nolint:golint,unused
}
