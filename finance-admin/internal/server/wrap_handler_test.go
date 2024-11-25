package server

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestStatusError_Code(t *testing.T) {
	assert.Equal(t, 0, StatusError(0).Code())
	assert.Equal(t, 200, StatusError(200).Code())
}

func TestStatusError_Error(t *testing.T) {
	assert.Equal(t, "0 ", StatusError(0).Error())
	assert.Equal(t, "200 OK", StatusError(200).Error())
	assert.Equal(t, "999 ", StatusError(999).Error())
}

type mockHandler struct {
	app    components.AppVars
	w      http.ResponseWriter
	r      *http.Request
	Err    error
	Called int
}

func (m *mockHandler) render(app components.AppVars, w http.ResponseWriter, r *http.Request) error {
	m.app = app
	m.w = w
	m.r = r
	m.Called = m.Called + 1
	return m.Err
}

func Test_wrapHandler_successful_request(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("opg-sirius-finance-admin"))
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "test-url/1", nil)

	nextHandlerFunc := wrapHandler(components.EnvironmentVars{})
	next := &mockHandler{}
	httpHandler := nextHandlerFunc(next)
	httpHandler.ServeHTTP(w, r)

	assert.Nil(t, next.Err)
	assert.Equal(t, w, next.w)
	assert.Equal(t, r, next.r)
	assert.Equal(t, 1, next.Called)
	assert.Equal(t, "test-url/1", next.app.Path)
	assert.Equal(t, 200, w.Result().StatusCode)
}

func Test_wrapHandler_status_error_handling(t *testing.T) {
	tests := []struct {
		error    error
		wantCode int
	}{
		{error: StatusError(400), wantCode: 400},
		{error: StatusError(401), wantCode: 401},
		{error: StatusError(403), wantCode: 403},
		{error: StatusError(404), wantCode: 404},
		{error: StatusError(500), wantCode: 500},
		{error: api.StatusError{Code: 400}, wantCode: 400},
		{error: api.StatusError{Code: 401}, wantCode: 401},
		{error: api.StatusError{Code: 403}, wantCode: 403},
		{error: api.StatusError{Code: 404}, wantCode: 404},
		{error: api.StatusError{Code: 500}, wantCode: 500},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx := telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("opg-sirius-finance-admin"))
			r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "test-url/1", nil)

			nextHandlerFunc := wrapHandler(components.EnvironmentVars{})
			next := &mockHandler{Err: test.error}
			httpHandler := nextHandlerFunc(next)
			httpHandler.ServeHTTP(w, r)

			assert.Equal(t, 1, next.Called)
			assert.Equal(t, w, next.w)
			assert.Equal(t, r, next.r)
			assert.Equal(t, test.wantCode, w.Result().StatusCode)
		})
	}
}
