package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LashkaPashka/EventConsumerService/internal/configs"
	fetchposts "github.com/LashkaPashka/EventConsumerService/internal/http-server/handlers/fetchPosts"
	"github.com/LashkaPashka/EventConsumerService/internal/kafka"
	"github.com/LashkaPashka/EventConsumerService/internal/rabbitmq"
	"github.com/LashkaPashka/EventConsumerService/internal/service"
	mongodb "github.com/LashkaPashka/EventConsumerService/internal/storage/mongoDb"
	"github.com/LashkaPashka/EventConsumerService/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal   = "local"
	envDev     = "dev"
	envProd    = "prod"
	postEvent  = "user.post.created"
	likeEvent  = "user.post.liked"
	orderEvent = "order.paid"
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

	// TODO: init redis-memcache
	redisSt := redis.New(cfg, logger)

	// TODO: init mongo-storage
	mongoSt := mongodb.New(cfg, logger)

	//TODO: init rabbitMQ
	// TODO: init Rabbit
	rbCli := rabbitmq.New(cfg, logger)

	// TODO: init Service
	service := service.New(mongoSt, redisSt, rbCli, logger, cfg)

	// TODO: init kafka-client
	kfCli := kafka.New(service, []string{cfg.KafkaBroker.Broker}, []string{postEvent, likeEvent, orderEvent}, cfg.GroupID, cfg, logger)

	// TODO: Start kafka in other gorutine
	go kfCli.Consumer()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/get-posts", fetchposts.New(redisSt, logger))

	// TODO: Creation server
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