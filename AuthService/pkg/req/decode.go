package req

import (
	"encoding/json"
	"io"
	"log/slog"
)

func decode[T any](body io.ReadCloser, logger *slog.Logger) T {
	const op = "AuthService.pkg.req.decode.go"

	var payload T

	err := json.NewDecoder(body).Decode(&payload)

	if err != nil {
		logger.Error("Invalid decode body request", slog.String("Error:", op))
		return payload
	}

	return payload
}