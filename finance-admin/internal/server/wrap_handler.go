package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/api"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"net/http"
	"time"
)

type ErrorVars struct {
	Code  int
	Error string
	Envs
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

func wrapHandler(errTmpl Template, errPartial string, envVars Envs) func(next HtmxHandler) http.Handler {
	return func(next HtmxHandler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			vars := NewAppVars(r, envVars)
			logger := telemetry.LoggerFromContext(r.Context())

			user := r.Context().(auth.Context).User
			if !user.IsFinanceReporting() {
				w.WriteHeader(http.StatusForbidden)
				errVars := ErrorVars{
					Code: http.StatusForbidden,
					Envs: envVars,
				}
				err := errTmpl.Execute(w, errVars)
				if err != nil {
					logger.Error("failed to render error template", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			err := next.render(vars, w, r)

			logger.Info(
				"page request",
				"duration", time.Since(start),
				"hx-request", r.Header.Get("HX-Request") == "true",
				"user-id", user.ID,
			)

			if err != nil {
				if errors.Is(err, api.ErrUnauthorized) {
					http.Redirect(w, r, envVars.SiriusURL+"/auth", http.StatusFound)
					return
				}

				code := http.StatusInternalServerError
				var serverStatusError StatusError
				if errors.As(err, &serverStatusError) {
					logger.Error("server error", "error", err)
					code = serverStatusError.Code()
				}
				var siriusStatusError api.StatusError
				if errors.As(err, &siriusStatusError) {
					logger.Error("sirius error", "error", err)
					code = siriusStatusError.Code
				}
				if errors.Is(err, context.Canceled) {
					code = 499 // Client Closed Request
				}

				w.Header().Add("HX-Retarget", "#main-container")
				w.WriteHeader(code)
				errVars := ErrorVars{
					Code:  code,
					Error: err.Error(),
					Envs:  envVars,
				}
				if IsHxRequest(r) {
					err = errTmpl.ExecuteTemplate(w, errPartial, errVars)
				} else {
					err = errTmpl.Execute(w, errVars)
				}

				if err != nil {
					logger.Error("failed to render error template", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
