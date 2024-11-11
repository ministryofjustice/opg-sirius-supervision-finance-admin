package server

import (
	"context"
	"errors"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandler struct {
	w      http.ResponseWriter
	r      *http.Request
	called bool
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.w = w
	m.r = r
	m.called = true
}

type mockAuthClient struct {
	validSession bool
	error        error
	called       bool
}

func (m *mockAuthClient) CheckUserSession(ctx api.Context) (bool, error) {
	m.called = true
	return m.validSession, m.error
}

func Test_authenticate_success(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("test"))
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "test-url/1", nil)

	client := &mockAuthClient{validSession: true}

	auth := Authenticator{
		Client: client,
		EnvVars: EnvironmentVars{
			SiriusURL: "v1/api/",
		},
	}
	next := &mockHandler{}
	auth.Authenticate(next).ServeHTTP(w, r)

	assert.Equal(t, true, client.called)
	assert.Equal(t, w, next.w)
	assert.Equal(t, r, next.r)
	assert.Equal(t, true, next.called)
	assert.Equal(t, 200, w.Result().StatusCode)
}

func Test_authenticate_unauthorised(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("test"))
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "test-url/1", nil)

	client := &mockAuthClient{validSession: false}

	auth := Authenticator{
		Client: client,
		EnvVars: EnvironmentVars{
			SiriusURL: "https://sirius.gov.uk/v1/api",
			Prefix:    "finance-admin",
		},
	}
	next := &mockHandler{}
	auth.Authenticate(next).ServeHTTP(w, r)

	assert.Equal(t, true, client.called)
	assert.Equal(t, false, next.called)
	assert.Equal(t, 302, w.Result().StatusCode)
	assert.Equal(t, "https://sirius.gov.uk/v1/api/auth?redirect=finance-admin%2Ftest-url%2F1", w.Result().Header.Get("Location"))
}

func Test_authenticate_error(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("test"))
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "test-url/1", nil)

	client := &mockAuthClient{validSession: false, error: errors.New("something went wrong")}

	auth := Authenticator{
		Client: client,
		EnvVars: EnvironmentVars{
			SiriusURL: "https://sirius.gov.uk/v1/api",
			Prefix:    "finance-admin",
		},
	}
	next := &mockHandler{}
	auth.Authenticate(next).ServeHTTP(w, r)

	assert.Equal(t, true, client.called)
	assert.Equal(t, false, next.called)
	assert.Equal(t, 302, w.Result().StatusCode)
	assert.Equal(t, "https://sirius.gov.uk/v1/api/auth?redirect=finance-admin%2Ftest-url%2F1", w.Result().Header.Get("Location"))
}
