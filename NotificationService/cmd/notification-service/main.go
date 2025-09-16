package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LashkaPashka/notification-service/internal/configs"
	"github.com/LashkaPashka/notification-service/internal/rabbitmq"
	"github.com/LashkaPashka/notification-service/internal/service"
	sendemail "github.com/LashkaPashka/notification-service/internal/service/send-email"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := configs.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info(
		"starting event-producer serivce",
		slog.String("env", cfg.Env),
		slog.String("Addr", cfg.Address),
		slog.String("version", "123"),
	)

	// TODO: init notify-email
	nEmai := sendemail.New(cfg, logger)

	// TODO: init service
	service := service.New(nEmai, cfg, logger)

	// TODO: init rabbitMQ
	rbCli := rabbitmq.New(service, cfg, logger)
	
	// TODO: start rabbit conusmer
	go rbCli.Consumer()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server")
		}
	}()

	logger.Info("server started")
	
	<-done
	logger.Info("stopping server")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server")
		return
	}

	logger.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log  = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
