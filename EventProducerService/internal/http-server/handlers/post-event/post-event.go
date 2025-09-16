package postevent

import (
	"log/slog"
	"net/http"

	convertjson "github.com/LashkaPashka/event-producer/internal/convert-json"
	"github.com/LashkaPashka/event-producer/internal/http-server/middleware/auth"
	"github.com/LashkaPashka/event-producer/internal/lib/encode"
	"github.com/LashkaPashka/event-producer/internal/lib/req"
	"github.com/LashkaPashka/event-producer/internal/payload"
)

const (
	topic = "user.post.created"
)


type Clienter interface {
	Producer(topic string, data []byte) 
}

func New(client Clienter, logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate payload
		body := req.HandleBody[payload.PostEventRequest](w, r, logger)
		if body == nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		logger.Debug("Validation went successfully!")

		// TODO: convert-json-model
		email := r.Context().Value(auth.Emailkey).(string)
		
		logger.Debug("Context was recieved successfully!", slog.Any("Email: ", email))

		eventModel := convertjson.ConvertForKafka(*body, email, logger)

		// TODO: convert to []bytes
		eventBytes := encode.EncodeBytes(eventModel)

		// TODO: write message to Kafka
		client.Producer(topic, eventBytes)
	}
}