package kafka

import (
	"context"
	"log/slog"

	"github.com/LashkaPashka/EventConsumerService/internal/configs"
	"github.com/segmentio/kafka-go"
)

type Serivcer interface {
	GetEventFromKafka(topic string, msg []byte)
}

type Client struct {
	Reader *kafka.Reader
	cfg *configs.Configs
	logger *slog.Logger
	serv Serivcer
}

func New(service Serivcer, brokers []string, topics []string, groupID string, cfg *configs.Configs, logger *slog.Logger) *Client {
	if len(brokers) == 0 || brokers[0] == "" || topics == nil || groupID == "" {
		logger.Error("Kafka connection parameters not specified")
		return nil
	}

	c := Client{logger: logger, cfg: cfg, serv: service}

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupTopics: topics,
		GroupID: groupID,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})

	return &c
}

func (c *Client) Consumer() []byte {
	const op = "EventConsumerService.kafka.Consumer"

	for {
		msg, err := c.Reader.FetchMessage(context.Background())
		if err != nil {
			c.logger.Error("Error fetchmessage", slog.String("Error: ", op))
			return nil
		}
		
		c.serv.GetEventFromKafka(msg.Topic, msg.Value)

		if err = c.Reader.CommitMessages(context.Background(), msg); err != nil {
			c.logger.Error("Error commit messages", slog.String("Error: ", op))
			return nil
		}
	}
}