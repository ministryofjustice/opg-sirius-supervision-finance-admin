package auth

import (
	"context"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
	"net/url"
)

type Context struct {
	context.Context
	Cookies   []*http.Cookie
	XSRFToken string
	user      *shared.User
}

func newContext(ctx context.Context, r *http.Request, user *shared.User) Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}
	return Context{
		Context:   ctx,
		Cookies:   r.Cookies(),
		XSRFToken: token,
		user:      user,
	}
}

type Client interface {
	GetUserSession(ctx context.Context) (*shared.User, error)
}

type EnvVars struct {
	SiriusPublicURL string
	Prefix          string
}

type Authenticator struct {
	Client  Client
	EnvVars EnvVars
}

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := telemetry.LoggerFromContext(ctx)

		user, err := a.Client.GetUserSession(ctx)
		if err != nil {
			logger.Error("Error validating session.", "error", err)
			http.Redirect(w, r, a.redirectPath(r.URL.RequestURI()), http.StatusFound)
			return
		}
		if user == nil {
			logger.Info("User session not valid. Redirecting.")
			http.Redirect(w, r, a.redirectPath(r.URL.RequestURI()), http.StatusFound)
			return
		}

		next.ServeHTTP(w, r.WithContext(newContext(ctx, r, user)))
	})
}

func (a *Authenticator) redirectPath(to string) string {
	return fmt.Sprintf("%s/auth?redirect=%s", a.EnvVars.SiriusPublicURL, url.QueryEscape(fmt.Sprintf("%s%s", a.EnvVars.Prefix, to)))
}
