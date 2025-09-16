package kafka

import (
	"context"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	Writer *kafka.Writer
	logger *slog.Logger
}

func New(brokers []string, logger *slog.Logger) *Client {
	if len(brokers) == 0 || brokers[0] == "" {
		logger.Error("Kafka connection parameters not specified")
		return nil 
	}

	c := Client{logger: logger}

	c.Writer = &kafka.Writer{
		Addr: kafka.TCP(brokers[0]),
		Balancer: &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &c
}

func (c *Client) Producer(topic string, data []byte) {
	const op = "EventProducerService.kafka.Producer"

	msg := kafka.Message {
		Topic: topic,
		Value: data,
	}
	
	if err := c.Writer.WriteMessages(context.Background(), msg); err != nil {
		c.logger.Error("Error write messeges", slog.String("Error: ", op))
	}

	c.logger.Info("Message was record successfuly!", slog.Any("Msg", msg.Value))
	time.Sleep(1*time.Second)
}