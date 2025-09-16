package rabbitmq

import (
	"context"
	"log/slog"
	"time"

	"github.com/LashkaPashka/EventConsumerService/internal/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string, logger *slog.Logger) {
  if err != nil {
	logger.Error(msg)
  }
}

type RabbitClient struct {
	ch *amqp.Channel
	logger *slog.Logger
}

func New(cfg *configs.Configs, logger *slog.Logger) *RabbitClient{
	conn, err := amqp.Dial(cfg.RabbitAddr)
	failOnError(err, "Failed to connect to RabbitMq", logger)
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel", logger)

	return &RabbitClient{
		ch: ch,
		logger: logger,
	}
}

func (r *RabbitClient) Producer(msg []byte) {
	// err := r.ch.ExchangeDeclare(
	// 	"event.notification.service",
	// 	"direct",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// failOnError(err, "Failed to declare a exchange", r.logger)

	// err = r.ch.QueueBind(
	// 	"send.message.service",
	// 	"key_events",
	// 	"event.notification.service",
	// 	false,
	// 	nil,
	// )
	// failOnError(err, "Failed to declare a binding", r.logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.ch.PublishWithContext(ctx,
		"event.notification.service",
		"key_events",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg,
		})
	failOnError(err, "Failed to publish a message", r.logger)

	r.logger.Info("Message was added successfully in RabbitMQ")
}