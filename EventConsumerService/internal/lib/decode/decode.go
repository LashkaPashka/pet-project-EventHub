package decode

import (
	"encoding/json"
	"log/slog"
)

func DecodeJSON[T any](payload []byte, logger *slog.Logger) T {
	const op = "EventConsumerService"
	
	var res T
	if err := json.Unmarshal(payload, &res); err != nil {
		logger.Error("Error unmarshal payload", slog.String("Error: ", op))
	}
	
	return res
} 