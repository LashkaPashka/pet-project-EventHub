package encode

import (
	"encoding/json"
	"log/slog"
)


func EncodeBytes(payload any, logger *slog.Logger) []byte {
	const op = "EventConsumerService.lib.encode.EncodeBytes"
	
	bytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Invalid marshal payload", slog.String("Error: ", op))
		return nil
	}

	return bytes
}