package server

import (
	"github.com/a-h/templ"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"io"
	"net/http"
)

type route struct {
	client  ApiClient
	envVars components.EnvironmentVars
}

func (r route) Client() ApiClient {
	return r.client
}

func (r route) execute(w io.Writer, req *http.Request, component templ.Component) error {
	ctx := req.Context()
	if IsHxRequest(req) {
		return component.Render(ctx, w)
	} else {
		var data components.PageVars
		data.EnvironmentVars = r.envVars

		return components.Page(data, component).Render(ctx, w)
	}
}

func (r route) getSuccess(req *http.Request) string {
	switch req.URL.Query().Get("success") {
	case "upload":
		return "File successfully uploaded"
	case "download":
		return "Your file has been successfully downloaded"
	}
	return ""
}

func IsHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
