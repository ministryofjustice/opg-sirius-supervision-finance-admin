package main

import (
	"context"
	"errors"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/api"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/auth"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/server"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"unicode"
)

type Envs struct {
	webDir          string
	siriusURL       string
	siriusPublicURL string
	backendURL      string
	hubURL          string
	prefix          string
	port            string
	jwtSecret       string
}

func parseEnvs() (*Envs, error) {
	envs := map[string]string{
		"SIRIUS_URL":        os.Getenv("SIRIUS_URL"),
		"SIRIUS_PUBLIC_URL": os.Getenv("SIRIUS_PUBLIC_URL"),
		"PREFIX":            os.Getenv("PREFIX"),
		"BACKEND_URL":       os.Getenv("BACKEND_URL"),
		"HUB_URL":           os.Getenv("HUB_URL"),
		"JWT_SECRET":        os.Getenv("JWT_SECRET"),
		"PORT":              os.Getenv("PORT"),
	}

	var missing []error
	for k, v := range envs {
		if v == "" {
			missing = append(missing, errors.New("missing environment variable: "+k))
		}
	}

	if len(missing) > 0 {
		return nil, errors.Join(missing...)
	}

	return &Envs{
		siriusURL:       envs["SIRIUS_URL"],
		siriusPublicURL: envs["SIRIUS_PUBLIC_URL"],
		prefix:          envs["PREFIX"],
		backendURL:      envs["BACKEND_URL"],
		hubURL:          envs["HUB_URL"],
		jwtSecret:       envs["JWT_SECRET"],
		webDir:          "web",
		port:            envs["PORT"],
	}, nil
}

func main() {
	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-finance-admin")

	err := run(ctx, logger)
	if err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	exportTraces := os.Getenv("TRACING_ENABLED") == "1"

	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	envs, err := parseEnvs()
	if err != nil {
		return err
	}

	client := api.NewClient(
		http.DefaultClient,
		&auth.JWT{
			Secret: envs.jwtSecret,
		},
		api.EnvVars{
			SiriusURL:  envs.siriusURL,
			BackendURL: envs.backendURL,
			HubURL:     envs.hubURL,
		})

	templates := createTemplates(envs)

	s := &http.Server{
		Addr: ":" + envs.port,
		Handler: server.New(logger, client, templates, server.Envs{
			WebDir:          envs.webDir,
			SiriusURL:       envs.siriusURL,
			SiriusPublicURL: envs.siriusPublicURL,
			Prefix:          envs.prefix,
		}),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Error("listen and server error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Running at :" + envs.port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", "sig", sig)

	tc, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.Shutdown(tc)
}

func createTemplates(envs *Envs) map[string]*template.Template {
	templates := map[string]*template.Template{}
	templateFunctions := map[string]interface{}{
		"contains": func(xs []string, needle string) bool {
			for _, x := range xs {
				if x == needle {
					return true
				}
			}

			return false
		},
		"title": func(s string) string {
			r := []rune(s)
			r[0] = unicode.ToUpper(r[0])

			return string(r)
		},
		"prefix": func(s string) string {
			return envs.prefix + s
		},
		"sirius": func(s string) string {
			return envs.siriusPublicURL + s
		},
	}

	templateDirPath := envs.webDir + "/template"
	templateDir, _ := os.Open(templateDirPath)
	templateDirs, _ := templateDir.Readdir(0)
	_ = templateDir.Close()

	mainTemplates, _ := filepath.Glob(templateDirPath + "/*.gotmpl")

	for _, file := range mainTemplates {
		tmpl := template.New(filepath.Base(file)).Funcs(templateFunctions)
		for _, dir := range templateDirs {
			if dir.IsDir() {
				tmpl, _ = tmpl.ParseGlob(templateDirPath + "/" + dir.Name() + "/*.gotmpl")
			}
		}
		tmpl, _ = tmpl.Parse(paginate.Template)
		templates[tmpl.Name()] = template.Must(tmpl.ParseFiles(file))
	}

	return templates
}
