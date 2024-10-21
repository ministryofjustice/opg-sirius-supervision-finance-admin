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
		ctx := r.Context()
		logger := telemetry.LoggerFromContext(ctx)

		sessionCookie, err := r.Cookie("sirius")
		if err != nil {
			logger.Debug("Missing session cookie. Redirecting.")
			http.Redirect(w, r, fmt.Sprintf("%s/auth?redirect=%s", a.EnvVars.SiriusURL, url.QueryEscape(a.EnvVars.Prefix+r.URL.Path)), http.StatusFound)
			return
		}

		_, _, sessionValid := a.Client.CheckUserSession(ctx, sessionCookie)
		if !sessionValid {
			logger.Debug("User session not valid. Redirecting.")
			http.Redirect(w, r, fmt.Sprintf("%s/auth?redirect=%s", a.EnvVars.SiriusURL, url.QueryEscape(a.EnvVars.Prefix+r.URL.Path)), http.StatusFound)
			return
		}

		//// Exchange session for JWT token
		//jwtToken := exchangeSessionForJWT()
		//
		//// Attach the JWT token to the request context
		//ctx := context.WithValue(r.Context(), jwtContextKey, jwtToken)

		next.ServeHTTP(w, r)
	})
}

// Helper function to get JWT token from context
//func getJWTFromContext(ctx context.Context) (string, bool) {
//	jwtToken, ok := ctx.Value(jwtContextKey).(string)
//	return jwtToken, ok
//}

// Simulate session to JWT exchange
//func exchangeSessionForJWT() string {
//	// Placeholder for actual JWT generation
//	return "generated-jwt-token"
//}
