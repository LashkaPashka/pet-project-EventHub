package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/LashkaPashka/EventHub/AuthService/internal/configs"
	"github.com/LashkaPashka/EventHub/AuthService/internal/server"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	// TODO: init configs
	cfg := configs.MustLoad()

	// TODO: init logger
	logger := setupLogger(cfg.Env)

	// TODO: init server
	server := server.NewServer(cfg, logger)

	// TODO: Run server

	go func() {
		logger.Info("Server' running", slog.String("ip-address", "127.0.0.1"), slog.Int("port", 8080))
		server.Run()
	}()

	// TODO: Stop server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	<-stop

	server.Stop()

	logger.Info("application stopped")
}


func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
		case envLocal:
			logger = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		case envDev:
			logger = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			)
		case envProd:
			logger = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
			)
	}

	return logger
}