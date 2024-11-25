package server

import (
	"github.com/a-h/templ"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"io"
	"net/http"
)

type mockRoute struct { //nolint:golint,unused
	client ApiClient
	error
}

func (r *mockRoute) Client() ApiClient { //nolint:golint,unused
	return r.client
}

func (r *mockRoute) execute(w io.Writer, req *http.Request, component templ.Component) error { //nolint:golint,unused
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
