package jwt

import (
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type JWTdata struct {
	Email string
}

type JWT struct {
	Secret string
	logger *slog.Logger
}

func NewJWT(secret string, logger *slog.Logger) *JWT {
	return &JWT{
		Secret: secret,
		logger: logger,
	}
}

func (j *JWT) CreateJWT(email string) (string, error) {
	const op = "AuthService.pkg.jwt.CreateJWT"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	sToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {	
		j.logger.Error("Error signatute token", slog.String("Error: ", op))
		return  "", err
	}

	return sToken, nil
}