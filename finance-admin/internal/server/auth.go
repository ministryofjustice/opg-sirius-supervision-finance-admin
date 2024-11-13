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

		sessionValid, err := a.Client.CheckUserSession(ctx)
		if err != nil {
			logger.Error("Error validating session.", "error", err)
			http.Redirect(w, r, a.redirectPath(r.URL.RequestURI()), http.StatusFound)
			return
		}
		if !sessionValid {
			logger.Info("User session not valid. Redirecting.")
			http.Redirect(w, r, a.redirectPath(r.URL.RequestURI()), http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) redirectPath(to string) string {
	return fmt.Sprintf("%s/auth?redirect=%s", a.EnvVars.SiriusPublicURL, url.QueryEscape(fmt.Sprintf("%s%s", a.EnvVars.Prefix, to)))
}
