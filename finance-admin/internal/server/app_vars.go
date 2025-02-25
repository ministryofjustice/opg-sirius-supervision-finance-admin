package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"net/http"
)

type AppVars struct {
	Path             string
	XSRFToken        string
	Tabs             []Tab
	EnvironmentVars  Envs
	ValidationErrors model.ValidationErrors
	Error            string
}

type Tab struct {
	Title    string
	Id       string
	Selected bool
}

func (t Tab) Path() string {
	return "/" + t.Id
}

func NewAppVars(r *http.Request, envVars Envs) AppVars {
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

	vars := AppVars{
		Path:            r.URL.Path,
		XSRFToken:       r.Context().(auth.Context).XSRFToken,
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
