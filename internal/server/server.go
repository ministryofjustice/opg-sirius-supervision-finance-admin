package server

import (
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

type ApiClient interface{}

type router interface {
	Client() ApiClient
	execute(http.ResponseWriter, *http.Request, any) error
}

type Template interface {
	Execute(wr io.Writer, data any) error
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

func New(logger *slog.Logger, client ApiClient, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {

	mux := http.NewServeMux()

	handleMux := func(pattern string, h Handler) {
		errors := wrapHandler(templates["error.gotmpl"], "main", envVars)
		mux.Handle(pattern, telemetry.Middleware(logger)(errors(h)))
	}

	handleMux("GET /downloads", &GetDownloadsHandler{&route{client: client, tmpl: templates["downloads.gotmpl"], partial: "downloads"}})
	handleMux("GET /uploads", &GetUploadsHandler{&route{client: client, tmpl: templates["uploads.gotmpl"], partial: "uploads"}})
	handleMux("GET /annual-invoicing-letters", &GetAnnualInvoicingLettersHandler{&route{client: client, tmpl: templates["annual_invoicing_letters.gotmpl"], partial: "annual-invoicing-letters"}})

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, securityheaders.Use(mux)), "supervision-finance-admin")
}
