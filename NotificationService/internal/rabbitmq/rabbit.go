package rabbitmq

import (
	"fmt"
	"log/slog"

	"github.com/LashkaPashka/notification-service/internal/configs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string, logger *slog.Logger) {
        if err != nil {
                logger.Error(fmt.Sprintf("%s:%s", msg, err))
        }
}

type Notifier interface {
	NotifyAboutPaid(payload string)
}

type RabbitClient struct {
	ch *amqp.Channel
	logger *slog.Logger
	Notifier
}


func New(notifier Notifier, cfg *configs.Configs, logger *slog.Logger) *RabbitClient {
	conn, err := amqp.Dial(cfg.RabbitAddr)
	failOnError(err, "Failed to connect to RabbitMQ", logger)
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel", logger)
	
	return &RabbitClient{
		ch: ch,
		logger: logger,
		Notifier: notifier,
	}
}

func (r *RabbitClient) Consumer() {
	err := r.ch.ExchangeDeclare(
		"event.notification.service",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange", r.logger)

	q, err := r.ch.QueueDeclare(
		"send.message.service",
		true,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue", r.logger)

	err = r.ch.QueueBind(
		q.Name,
		"key_events",
		"event.notification.service",
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue", r.logger)

	msgs, err := r.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer", r.logger)

	go func() {
		for d := range msgs {
			r.logger.Info("Recieved meassage from RabbitMQ", slog.String("Data", string(d.Body)))

			r.Notifier.NotifyAboutPaid(string(d.Body))
			
		}
	}()

	select{}
}