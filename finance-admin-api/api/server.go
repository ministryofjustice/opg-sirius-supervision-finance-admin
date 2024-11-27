package api

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/event"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/db"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"log/slog"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatch interface {
	FinanceAdminUpload(ctx context.Context, event event.FinanceAdminUpload) error
}

type FileStorage interface {
	GetFile(ctx context.Context, bucketName string, filename string, versionID string) (*s3.GetObjectOutput, error)
	PutFile(ctx context.Context, bucketName string, fileName string, file io.Reader) (*string, error)
	FileExists(ctx context.Context, bucketName string, filename string, versionID string) bool
}

type DbConn interface {
	Run(ctx context.Context, query db.ReportQuery) ([][]string, error)
}

type Server struct {
	http        HTTPClient
	conn        DbConn
	dispatch    Dispatch
	filestorage FileStorage
}

func NewServer(httpClient HTTPClient, conn DbConn, dispatch Dispatch, filestorage FileStorage) Server {
	return Server{httpClient, conn, dispatch, filestorage}
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

	handleFunc("GET /download", s.download)
	handleFunc("HEAD /download", s.checkDownload)

	handleFunc("POST /downloads", s.requestReport)
	handleFunc("POST /uploads", s.upload)

	handleFunc("POST /events", s.handleEvents)

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

type StatusError struct {
	Code   int    `json:"code"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

func newStatusError(resp *http.Response) StatusError {
	return StatusError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL.String(),
		Method: resp.Request.Method,
	}
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.URL, e.Code)
}
