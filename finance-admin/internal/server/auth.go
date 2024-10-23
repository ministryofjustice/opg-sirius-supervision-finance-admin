package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"net/http"
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
			http.Redirect(w, r, fmt.Sprintf("%s/auth", a.EnvVars.SiriusURL), http.StatusFound)
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
