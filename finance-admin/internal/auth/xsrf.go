package auth

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/telemetry"
)

func XsrfCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			ctx := r.Context().(Context)

			// Limit request body size to 10MB to prevent memory exhaustion
			r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

			if r.FormValue("xsrfToken") != ctx.XSRFToken {
				logger := telemetry.LoggerFromContext(ctx)
				logger.Error("XSRF token mismatch")
				http.Error(w, "XSRF token mismatch", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
