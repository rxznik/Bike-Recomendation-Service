package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/app"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/config"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/observability"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/setup"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	sentryFlushTimeout = 3 * time.Second
	serviceName        = "recomendation-service"
)

var (
	ctx         = context.Background()
	cfg         *config.Config
	logger      *zap.Logger
	promobserv  *observability.Prometheus
	application *app.App
)

func init() {
	godotenv.Load()
	cfg = config.MustLoad()
	logger = setup.MustSetupLogger(cfg.Env)
	promobserv = observability.NewPrometheus(logger)
	application = setup.MustSetupRabbitMQApplication(logger, cfg, promobserv)
	observability.InitGlobalSentry(cfg.SentryDSN, cfg.Env, serviceName)
}

func main() {
	shutdownTracer := observability.MustInitTracer(ctx, serviceName)

	defer func() {
		logger.Sync()
		observability.FlushAndRecoverSentry(sentryFlushTimeout)
		shutdownTracer()
	}()
	logger.Info("starting application...")

	// application
	go func(ctx context.Context) {
		if err := application.Run(ctx); err != nil {
			logger.Fatal("failed to run application", zap.Error(err))
		}
	}(ctx)

	// prometheus
	go func() {
		if err := promobserv.Run(); err != nil {
			logger.Fatal("failed to run prometheus", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	logger.Info("shutting down prometheus...")

	if err := promobserv.Stop(); err != nil {
		logger.Fatal("failed to stop prometheus", zap.Error(err))
	}

	logger.Info("shutting down application...")

	if err := application.Stop(); err != nil {
		logger.Fatal("failed to stop application", zap.Error(err))
	}

	logger.Info("application shutdown")
}
