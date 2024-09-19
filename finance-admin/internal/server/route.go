package server

import (
	"net/http"
)

type PageData struct {
	Data any
}

type route struct {
	client  ApiClient
	tmpl    Template
	partial string
}

func (r route) Client() ApiClient {
	return r.client
}

// execute is an abstraction of the Template execute functions in order to conditionally render either a full template or
// a block, in response to a header added by HTMX. If the header is not present, the function will also fetch all
// additional data needed by the page for a full page load.
func (r route) execute(w http.ResponseWriter, req *http.Request, data any) error {
	if IsHxRequest(req) {
		return r.tmpl.ExecuteTemplate(w, r.partial, data)
	} else {
		data := PageData{
			Data: data,
		}
		return r.tmpl.Execute(w, data)
	}
}

func IsHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
