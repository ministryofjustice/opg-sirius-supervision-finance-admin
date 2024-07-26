package server

import (
	"net/http"
	"net/url"
)

type AppVars struct {
	Path            string
	XSRFToken       string
	Tabs            []Tab
	EnvironmentVars EnvironmentVars
	Error           string
}

type Tab struct {
	Title    string
	Id       string
	Selected bool
}

func (t Tab) Path() string {
	return "/" + t.Id
}

func NewAppVars(r *http.Request, envVars EnvironmentVars) AppVars {
	tabs := []Tab{
		{
			Id:    "downloads",
			Title: "Downloads",
		},
		{
			Id:    "uploads",
			Title: "Uploads",
		},
		{
			Id:    "annual-invoicing-letters",
			Title: "Annual Invoicing Letters",
		},
	}

	var token string
	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	vars := AppVars{
		Path:            r.URL.Path,
		XSRFToken:       token,
		EnvironmentVars: envVars,
		Tabs:            tabs,
	}

	return vars
}

func (a *AppVars) selectTab(s string) {
	for i, tab := range a.Tabs {
		if tab.Id == s {
			a.Tabs[i] = Tab{
				Title:    tab.Title,
				Id:       tab.Id,
				Selected: true,
			}
		}
	}
}
