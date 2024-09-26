package main

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/opg-sirius-finance-admin/finance-admin-api/api"
	"github.com/opg-sirius-finance-admin/finance-admin-api/awsclient"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-finance-admin-api")

	err := run(ctx, logger)
	if err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	exportTraces := env.Get("TRACING_ENABLED", "0") == "1"

	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	awsClient, err := awsclient.NewClient(ctx)
	if err != nil {
		return err
	}

	server := api.NewServer(awsClient)

	s := &http.Server{
		Addr:    ":8080",
		Handler: server.SetupRoutes(logger),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Error("listen and server error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("Running at :8080")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", "sig", sig)

	tc, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.Shutdown(tc)
}
