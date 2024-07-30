package components

import (
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"net/http"
	"net/url"
	"os"
)

type EnvironmentVars struct {
	Port            string
	WebDir          string
	SiriusURL       string
	SiriusPublicURL string
	BackendURL      string
	Prefix          string
}

func NewEnvironmentVars() EnvironmentVars {
	return EnvironmentVars{
		Port:            getEnv("PORT", "1234"),
		WebDir:          getEnv("WEB_DIR", "web"),
		SiriusURL:       getEnv("SIRIUS_URL", "http://host.docker.internal:8080"),
		SiriusPublicURL: getEnv("SIRIUS_PUBLIC_URL", ""),
		Prefix:          getEnv("PREFIX", ""),
		BackendURL:      getEnv("BACKEND_URL", ""),
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}

type AppVars struct {
	Path             string
	XSRFToken        string
	EnvironmentVars  EnvironmentVars
	ValidationErrors model.ValidationErrors
	Error            string
}

func NewAppVars(r *http.Request, envVars EnvironmentVars) AppVars {
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
	}

	return vars
}
