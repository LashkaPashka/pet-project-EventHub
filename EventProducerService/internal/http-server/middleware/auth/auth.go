package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/LashkaPashka/event-producer/internal/configs"
	"github.com/LashkaPashka/event-producer/internal/lib/jwt"
)

type key string

const (
	Emailkey key = "Emailkey"
)

func WriteStatusCode(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}


func IsAuthed(next http.Handler, cfg *configs.Configs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer ") {
			WriteStatusCode(w)
			return 
		}

		token := strings.TrimPrefix(authedHeader, "Bearer ")
		isValid, data := jwt.NewJwt(cfg.JWT_Secret).Parse(token)
		if !isValid {
			WriteStatusCode(w)
			return 
		}

		ctx := context.WithValue(r.Context(), Emailkey, data.Email)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	}
}