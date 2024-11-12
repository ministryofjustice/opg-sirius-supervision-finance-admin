package server

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"io"
	"net/http"
)

type mockTemplate struct {
	executed         bool
	executedTemplate bool
	lastVars         interface{}
	lastW            io.Writer
	error            error
}

func (m *mockTemplate) Execute(w io.Writer, vars any) error {
	m.executed = true
	m.lastVars = vars
	m.lastW = w
	return m.error
}

func (m *mockTemplate) ExecuteTemplate(w io.Writer, name string, vars any) error {
	m.executedTemplate = true
	m.lastVars = vars
	m.lastW = w
	return m.error
}

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
	error            error
	downloadResponse *http.Response
}

func (m mockApiClient) Upload(context api.Context, data shared.Upload) error {
	return m.error
}

func (m mockApiClient) RequestReport(context api.Context, data model.ReportRequest) error {
	return m.error
}

func (m mockApiClient) Download(ctx api.Context, uid string) (*http.Response, error) {
	return m.downloadResponse, m.error
}

func (m mockApiClient) CheckDownload(ctx api.Context, uid string) error {
	return m.error
}
