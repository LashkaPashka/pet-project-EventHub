package service

import (
	"log/slog"

	"github.com/LashkaPashka/notification-service/internal/configs"
	getemail "github.com/LashkaPashka/notification-service/internal/lib/get-email"
	sendemail "github.com/LashkaPashka/notification-service/internal/service/send-email"
	sendphonenumber "github.com/LashkaPashka/notification-service/internal/service/send-phoneNumber"
)

type Notifier interface {
	*sendemail.NotifyEmail | *sendphonenumber.NotifyPhone

	Send(to, subject, body string) error
}

type Service[T Notifier] struct {
	cfg *configs.Configs
	logger *slog.Logger
	notifier T
}

func New[T Notifier](notifier T, cfg *configs.Configs, logger *slog.Logger) *Service[T] {
	return  &Service[T]{
		cfg: cfg,
		logger: logger,
		notifier: notifier,
	}
}

func (s *Service[T]) NotifyAboutPaid(payload string) {
	subject := "Notification about paid"
	to := getemail.GetEmail([]byte(payload))

	s.logger.Debug("Prepare message")

	s.notifier.Send(to, subject, payload)
}