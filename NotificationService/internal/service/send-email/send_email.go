package sendemail

import (
	"log/slog"
	"net/smtp"

	"github.com/LashkaPashka/notification-service/internal/configs"
)

type NotifyEmail struct {
	cfg *configs.Configs
	logger *slog.Logger
}

func New(cfg *configs.Configs, logger *slog.Logger) *NotifyEmail{
	return  &NotifyEmail{
		cfg: cfg,
		logger: logger,
	}
}

func (n *NotifyEmail) Send(to, subject, body string) error {
	const op = "NotificationService.send-email.Send"
	
	_ = n.cfg

	from := "balasanianraf@yandex.ru"
	password := "tgjulhqwvgtjwups"
	smtpHost := "smtp.yandex.ru"
	smtpPort := "587"
	
	msg := []byte(body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		n.logger.Error("Invalid sending message on email", slog.Any("Reason: ", err), slog.String("Error: ", op))
		return err
	}

	n.logger.Info("Message was sent successfully on email!")

	return nil
}