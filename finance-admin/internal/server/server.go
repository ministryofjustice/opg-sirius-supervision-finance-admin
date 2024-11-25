package server

import (
	"github.com/a-h/templ"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type ApiClient interface {
	RequestReport(api.Context, model.ReportRequest) error
	Upload(api.Context, shared.Upload) error
	CheckDownload(api.Context, string) error
	Download(api.Context, string) (*http.Response, error)
}

type AuthClient interface {
	CheckUserSession(ctx api.Context) (bool, error)
}

type router interface {
	Client() ApiClient
	execute(io.Writer, *http.Request, templ.Component) error
}

func New(logger *slog.Logger, client *api.Client, envVars components.EnvironmentVars) http.Handler {
	r := route{client: client, envVars: envVars}
	mux := http.NewServeMux()
	auth := Authenticator{
		Client:  client,
		EnvVars: envVars,
	}

	handleMux := func(pattern string, h Handler) {
		errors := wrapHandler(envVars)
		mux.Handle(pattern, telemetry.Middleware(logger)(errors(h)))
	}

	handleMux("GET /downloads", &DownloadsTabHandler{&r})
	handleMux("GET /uploads", &UploadsTabHandler{&r})
	handleMux("GET /annual-invoicing-letters", &AnnualInvoicingLettersTabHandler{&r})

	//forms
	handleMux("POST /request-report", &RequestReportHandler{&r})
	handleMux("POST /uploads", &UploadFormHandler{&r})

	// file download
	handleMux("GET /download", &GetDownloadHandler{&r})

	downloadMux := func(pattern string, h http.Handler) {
		mux.Handle(pattern, telemetry.Middleware(logger)(auth.Authenticate(h)))
	}

	downloadMux("GET /download/callback", downloadCallback(client))

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, securityheaders.Use(mux)), "supervision-finance-admin")
}

func getContext(r *http.Request) api.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return api.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
