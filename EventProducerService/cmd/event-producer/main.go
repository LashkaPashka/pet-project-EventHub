package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LashkaPashka/event-producer/internal/configs"
	likeevent "github.com/LashkaPashka/event-producer/internal/http-server/handlers/like-event"
	orderevent "github.com/LashkaPashka/event-producer/internal/http-server/handlers/order-event"
	"github.com/LashkaPashka/event-producer/internal/http-server/handlers/post-event"
	"github.com/LashkaPashka/event-producer/internal/http-server/middleware/auth"
	"github.com/LashkaPashka/event-producer/internal/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// TODO: init kafka-client

	client := kafka.New([]string{cfg.KafkaBroker.Broker}, logger)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/producer", func(r chi.Router) {
		r.Post("/post-event", auth.IsAuthed(postevent.New(client, logger), cfg))
		r.Get("/like-event", auth.IsAuthed(likeevent.New(client, logger), cfg))
		r.Post("/order-event", auth.IsAuthed(orderevent.New(client, logger), cfg))
	})
	
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
		Handler: router,
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
