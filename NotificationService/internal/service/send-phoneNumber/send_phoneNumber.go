package sendphonenumber

import (
	"log/slog"

	"github.com/LashkaPashka/notification-service/internal/configs"
)

type NotifyPhone struct {
	cfg *configs.Configs
	logger *slog.Logger
}

func (n *NotifyPhone) Send(to, subject, body string) error {
	return nil
}