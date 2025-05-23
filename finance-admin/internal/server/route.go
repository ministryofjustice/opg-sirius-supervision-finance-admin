package server

import (
	"fmt"
	"net/http"
)

type PageData struct {
	Data           any
	SuccessMessage string
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
			Data:           data,
			SuccessMessage: r.getSuccess(req),
		}
		return r.tmpl.Execute(w, data)
	}
}

func (r route) getSuccess(req *http.Request) string {
	switch req.URL.Query().Get("success") {
	case "upload":
		return "File successfully uploaded"
	case "download":
		return "Your file has been successfully downloaded"
	case "request_report":
		return fmt.Sprintf("Your %s report is being prepared. You will be notified by email when it’s ready for download.", req.URL.Query().Get("report_type"))
	}
	return ""
}

func IsHxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
