package service

import (
	"log/slog"
	"strings"

	"github.com/LashkaPashka/EventConsumerService/internal/configs"
	"github.com/LashkaPashka/EventConsumerService/internal/lib/decode"
	"github.com/LashkaPashka/EventConsumerService/internal/lib/encode"
	"github.com/LashkaPashka/EventConsumerService/internal/service/model"
)

const (
	userPostCreated = "user.post.created"
	userPostLiked = "user.post.liked"
	orderPaid = "order.paid"

	feedLatestPost = "feed:latest_posts"
)

type MongoStorage interface {
	InsertDoc(eventModel any) any
}

type RedisStorage interface {
	SaveInMemory(data any, key string, id string)
}

type RabbitMQ interface{
	Producer(msg []byte)
}


type EventConsumerService struct {
	logger *slog.Logger
	cfg *configs.Configs
	MongoStorage
	RedisStorage
	RabbitMQ
}

func New(mongoSt MongoStorage, redisSt RedisStorage, rabbitMQ RabbitMQ, logger *slog.Logger, cfg *configs.Configs) *EventConsumerService{
	return &EventConsumerService{
		logger: logger,
		cfg: cfg,
		MongoStorage: mongoSt,
		RedisStorage: redisSt,
		RabbitMQ: rabbitMQ,
	}
}

func (s *EventConsumerService) GetEventFromKafka(topic string, msg []byte) {
	// TODO: unmarshal payload
	var eventDecode any
	var id string

	switch topic {
		case userPostCreated:
			eventDecode = decode.DecodeJSON[model.UserPostCreated](msg, s.logger)
			id = eventDecode.(model.UserPostCreated).ID
		case userPostLiked:
			eventDecode = decode.DecodeJSON[model.UserPostLiked](msg, s.logger)
			id = eventDecode.(model.UserPostCreated).ID
		case orderPaid:
			eventDecode = decode.DecodeJSON[model.OrderPaid](msg, s.logger)
			id = eventDecode.(model.OrderPaid).ID
		default:
			s.logger.Warn("Unknown topic")
			return
	}

	s.logger.Info("Message received from kafka", slog.Any("Msg", eventDecode))

	// TODO: Insert mongoDB
	ids := s.MongoStorage.InsertDoc(eventDecode)

	s.logger.Info("Payload add in MongoDb successfuly", slog.Any("ID", ids))

	// TODO: Insert redis
	s.RedisStorage.SaveInMemory(eventDecode, feedLatestPost, id)

	// TODO: Send notification on rabbitMQ
	if strings.Compare(topic, orderPaid) == 0 {
		s.RabbitMQ.Producer(encode.EncodeBytes(eventDecode, s.logger))
	}
}