package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"net/http"
	"net/url"
)

type Authenticator struct {
	Client  AuthClient
	EnvVars EnvironmentVars
}

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getContext(r)
		logger := telemetry.LoggerFromContext(ctx.Context)

		_, err, sessionValid := a.Client.CheckUserSession(ctx)
		if err != nil {
			logger.Error("Error validating session.", "error", err)
		}
		if !sessionValid {
			logger.Info("User session not valid. Redirecting.")
			http.Redirect(w, r, fmt.Sprintf("%s/auth?redirect=%s", a.EnvVars.SiriusURL, url.QueryEscape(a.EnvVars.Prefix+r.URL.Path)), http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
