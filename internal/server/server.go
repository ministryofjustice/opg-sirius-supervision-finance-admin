package server

import (
	"github.com/a-h/templ"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/internal/components"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"log/slog"
	"net/http"
)

type ApiClient interface{}

type router interface {
	Client() ApiClient
	execute(io.Writer, *http.Request, templ.Component) error
}

func New(logger *slog.Logger, client ApiClient, envVars components.EnvironmentVars) http.Handler {
	r := route{client: client, envVars: envVars}
	mux := http.NewServeMux()

	handleMux := func(pattern string, h Handler) {
		errors := wrapHandler(envVars)
		mux.Handle(pattern, telemetry.Middleware(logger)(errors(h)))
	}

	handleMux("GET /downloads", &GetDownloadsHandler{&r})
	handleMux("GET /uploads", &GetUploadsHandler{&r})
	handleMux("GET /annual-invoicing-letters", &GetAnnualInvoicingLettersHandler{&r})

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, securityheaders.Use(mux)), "supervision-finance-admin")
}
