package convertjson

import (
	"log/slog"
	"time"

	"github.com/LashkaPashka/event-producer/internal/lib/random"
	"github.com/LashkaPashka/event-producer/internal/model"
	modelforevents "github.com/LashkaPashka/event-producer/internal/model/model-for-events"
	"github.com/LashkaPashka/event-producer/internal/payload"
)

type Constraints interface {
	payload.PostEventRequest | payload.LikeEventRequest | payload.OrderEventRequest
}

func ConvertForKafka[T Constraints](body T, email string, logger *slog.Logger) *model.EventForKafka{
	const op = "EventProducerservice.convert-json.ConvertForKafka"

	id := random.NewRandomString("evt_", 6)
	trace_id := random.NewRandomString("", 6)
	
	var eventModel = model.EventForKafka{
		ID: id,
		Timestamp: time.Now(),
		Source: "event-producer",
		MetaM: model.Meta{
			Trace_id: trace_id,
			Version: "1.0",
		},
	}

	logger.Debug("Model was created successfully", slog.String("Debug: ", op))

	switch v := any(body).(type) {
		case payload.PostEventRequest:
			WrapPostEvent(&eventModel, v, email)
		case payload.LikeEventRequest:
			WrapLikeEvent(&eventModel, v, email)
		case payload.OrderEventRequest:
			WrapOrderEvent(&eventModel, v, email)
		default:
			logger.Warn("Unknown event type", slog.String("op", op))
    		return nil
	}

	logger.Debug("EventModel was upgrated successfullt", slog.String("Debug: ", op))

	return &eventModel
}

func WrapPostEvent(event *model.EventForKafka, body payload.PostEventRequest, email string) {
	event.Type = "user.post.created"
	event.DataM = &modelforevents.PostData{
		Email: email,
		Title: body.Title,
		Content: body.Content,
		Tags: body.Tags,
		Timestamp: time.Now(),
	}
}

func WrapLikeEvent(event *model.EventForKafka, body payload.LikeEventRequest, email string) {
	event.Type = "user.post.liked"
	event.DataM = &modelforevents.LikeData{
		Email: email,
		Timestamp: time.Now(),
	}
}

func WrapOrderEvent(event *model.EventForKafka, body payload.OrderEventRequest, email string) {
	event.Type = "order.paid"
	event.DataM = &modelforevents.OrderData{
		Email: email,
		Amount: body.Amount,
		Currency: body.Currency,
		PaymentMethod: body.PaymentMethod,
		Timestamp: time.Now(),
	}
}