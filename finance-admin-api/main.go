package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/api"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/event"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/filestorage"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/reports"
	"log/slog"
	"net/http"
	"net/url"
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

	filestorageclient, err := filestorage.NewClient(ctx)
	if err != nil {
		return err
	}

	eventClient := setupEventClient(ctx, logger)

	db := setupDbPool(ctx, logger)
	reportsClient := reports.NewClient(ctx, db)
	defer reportsClient.Close()

	server := api.NewServer(http.DefaultClient, reportsClient, eventClient, filestorageclient)

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

func setupEventClient(ctx context.Context, logger *slog.Logger) *event.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error("failed to load aws config", slog.Any("err", err))
	}

	// set endpoint to "" outside dev to use default AWS resolver
	if endpointURL := os.Getenv("AWS_BASE_URL"); endpointURL != "" {
		cfg.BaseEndpoint = aws.String(endpointURL)
	}

	return event.NewClient(cfg, os.Getenv("EVENT_BUS_NAME"))
}

func setupDbPool(ctx context.Context, logger *slog.Logger) *pgxpool.Pool {
	dbConn := os.Getenv("POSTGRES_CONN")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDb := os.Getenv("POSTGRES_DB")

	dbpool, err := pgxpool.New(ctx, fmt.Sprintf("postgresql://%s:%s@%s/%s", dbUser, url.QueryEscape(dbPassword), dbConn, pgDb))
	if err != nil {
		logger.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	return dbpool
}
