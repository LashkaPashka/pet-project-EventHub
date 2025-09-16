package sendemail

import (
	"crypto/tls"
	"fmt"
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

	from := n.cfg.Email.From
	password := n.cfg.Email.Password
	smtpHost := n.cfg.Email.SmtpHost
	smtpPort := n.cfg.Email.SmtpPort
	
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"\r\n%s\r\n",
		from, to, subject, body,
	))
	
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	// TLS-сессия
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		n.logger.Error("TLS connect error", slog.Any("Reason", err), slog.String("Error", op))
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		n.logger.Error("SMTP client error", slog.Any("Reason", err), slog.String("Error", op))
		return err
	}
	defer client.Close()

	// Авторизация
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err := client.Auth(auth); err != nil {
		return err
	}

	// Адреса
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	// Тело письма
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	if err := client.Quit(); err != nil {
		return err
	}

	n.logger.Info("Message was sent successfully on email!")
	return nil
}