package services

import (
	"context"
	"worker-service/config"
	tasks "worker-service/internal/dto"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(ctx context.Context, task tasks.EmailTask) error
}

type emailService struct {
	cfg config.AppConfig
}

func NewEmailService(cfg config.AppConfig) EmailService {
	return &emailService{
		cfg: cfg,
	}
}

func (e *emailService) SendEmail(ctx context.Context, task tasks.EmailTask) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", e.cfg.Smtp.Email)
	mailer.SetHeader("To", task.To)
	mailer.SetHeader("Subject", task.Subject)
	mailer.SetBody("text/html", task.Body)

	dialer := gomail.NewDialer(
		e.cfg.Smtp.Host,
		e.cfg.Smtp.Port,
		e.cfg.Smtp.Email,
		e.cfg.Smtp.Password,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	logrus.Info("email sent successfully!")
	return nil
}
