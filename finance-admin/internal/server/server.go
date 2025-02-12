package server

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/api"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

type ApiClient interface {
	RequestReport(context.Context, shared.ReportRequest) error
	Upload(context.Context, shared.Upload) error
	CheckDownload(context.Context, string) error
	Download(context.Context, string) (*http.Response, error)
}

type router interface {
	Client() ApiClient
	execute(http.ResponseWriter, *http.Request, any) error
}

type Template interface {
	Execute(wr io.Writer, data any) error
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type HtmxHandler interface {
	render(app AppVars, w http.ResponseWriter, r *http.Request) error
}

func New(logger *slog.Logger, client *api.Client, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	mux := http.NewServeMux()
	authenticator := auth.Authenticator{
		Client: client,
		EnvVars: auth.EnvVars{
			SiriusPublicURL: envVars.SiriusPublicURL,
			Prefix:          envVars.Prefix,
		},
	}

	handleMux := func(pattern string, h HtmxHandler) {
		errors := wrapHandler(templates["error.gotmpl"], "main", envVars)
		mux.Handle(pattern, telemetry.Middleware(logger)(authenticator.Authenticate(errors(h))))
	}

	// tabs
	handleMux("GET /downloads", &DownloadsTabHandler{&route{client: client, tmpl: templates["downloads.gotmpl"], partial: "downloads"}})
	handleMux("GET /uploads", &UploadsTabHandler{&route{client: client, tmpl: templates["uploads.gotmpl"], partial: "uploads"}})
	handleMux("GET /annual-invoicing-letters", &AnnualInvoicingLettersTabHandler{&route{client: client, tmpl: templates["annual_invoicing_letters.gotmpl"], partial: "annual-invoicing-letters"}})

	//forms
	handleMux("POST /request-report", &RequestReportHandler{&route{client: client, tmpl: templates["downloads.gotmpl"], partial: "error-summary"}})
	handleMux("POST /uploads", &UploadFormHandler{&route{client: client, tmpl: templates["uploads.gotmpl"], partial: "error-summary"}})

	// file download
	handleMux("GET /download", &GetDownloadHandler{&route{client: client, tmpl: templates["download-button.gotmpl"], partial: "download"}})

	downloadMux := func(pattern string, h http.Handler) {
		mux.Handle(pattern, telemetry.Middleware(logger)(authenticator.Authenticate(h)))
	}

	downloadMux("GET /download/callback", downloadCallback(client))

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, securityheaders.Use(mux)), "supervision-finance-admin")
}
