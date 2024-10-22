package api

import (
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/finance-admin-api/awsclient"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log/slog"
	"net/http"
)

type Server struct {
	awsClient awsclient.AWSClient
}

func NewServer(awsClient awsclient.AWSClient) Server {
	return Server{awsClient}
}

func (s *Server) SetupRoutes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/health-check", healthCheck())

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, h handlerFunc) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, h)
		mux.Handle(pattern, handler)
	}

	handleFunc("GET /downloads/{filename}", s.download)
	handleFunc("POST /uploads", s.upload)

	return otelhttp.NewHandler(telemetry.Middleware(logger)(securityheaders.Use(s.RequestLogger(mux))), "supervision-finance-admin-api")
}

func (s *Server) RequestLogger(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health-check" {
			telemetry.LoggerFromContext(r.Context()).Info(
				"API Request",
				"method", r.Method,
				"uri", r.URL.RequestURI(),
			)
		}
		h.ServeHTTP(w, r)
	}
}
