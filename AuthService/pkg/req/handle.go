package req

import (
	"log/slog"
	"net/http"
)


func HandleBody[T any](w http.ResponseWriter, r *http.Request, logger *slog.Logger) *T{
	// TODO: decode body	
	body := decode[T](r.Body, logger)

	// TODO: validate decode body
	if err := isValid(body, logger); err != nil {
		return nil
	}

	return &body
}