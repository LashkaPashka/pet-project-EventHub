package register

import (
	"log/slog"

	"github.com/LashkaPashka/EventHub/AuthService/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	FindByEmail(email string) *model.User
	Create(user *model.User) int
}

type AuthSerice struct {
	storage Storage
}

func NewAuthSerice(storage Storage) *AuthSerice{
	return &AuthSerice{
		storage: storage,
	}
}

func (s *AuthSerice) Register(email, password, name string, logger *slog.Logger) string {
	const op = "AuthService.service.Register"
	
	exsiting_user := s.storage.FindByEmail(email)
	if exsiting_user != nil {
		logger.Error("User've already created", slog.String("Error: ", op))
		return ""
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Invalid hashedPassword", slog.String("Error: ", op))
		return ""
	}

	user := &model.User{
		Email: email,
		Password: string(hashedPassword),
		Username: name,
	}
	
	s.storage.Create(user)

	return user.Email
}

func (s *AuthSerice) Login(email, password string, logger *slog.Logger) string {
	user := s.storage.FindByEmail(email)
	
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error("Wrong Credetials!")
		return ""
	}

	return user.Email
}