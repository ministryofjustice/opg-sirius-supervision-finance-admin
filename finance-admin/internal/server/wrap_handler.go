package server

import (
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/components"
	"net/http"
	"time"
)

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler interface {
	render(app components.AppVars, w http.ResponseWriter, r *http.Request) error
}

func wrapHandler(envVars components.EnvironmentVars) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()

			vars := components.NewAppVars(r, envVars)
			err := next.render(vars, w, r)

			logger := telemetry.LoggerFromContext(ctx)

			logger.Info(
				"Page Request",
				"duration", time.Since(start),
				"hx-request", r.Header.Get("HX-Request") == "true",
			)

			if err != nil {
				if errors.Is(err, api.ErrUnauthorized) {
					http.Redirect(w, r, envVars.SiriusURL+"/auth", http.StatusFound)
					return
				}

				logger.Error("Page Error", "error", err)

				code := http.StatusInternalServerError
				var serverStatusError StatusError
				if errors.As(err, &serverStatusError) {
					code = serverStatusError.Code()
				}
				var siriusStatusError api.StatusError
				if errors.As(err, &siriusStatusError) {
					code = siriusStatusError.Code
				}

				w.Header().Add("HX-Retarget", "#main-container")
				w.WriteHeader(code)
				errVars := components.ErrorVars{
					Code:   code,
					Error:  err.Error(),
					Prefix: envVars.Prefix,
				}
				component := components.Error(errVars)
				if IsHxRequest(r) {
					err = component.Render(ctx, w)
				} else {
					var data components.PageVars
					data.EnvironmentVars = envVars

					err = components.Page(data, component).Render(ctx, w)
				}

				if err != nil {
					logger.Error("failed to render error template", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
