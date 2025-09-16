package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)


type JWTdata struct {
	Email string
}

type JWT struct {
	Secret string
}

func NewJwt(secret string) *JWT{
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Parse(token string) (bool, *JWTdata){
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"]

	return t.Valid, &JWTdata{
		Email: email.(string),
	}
}