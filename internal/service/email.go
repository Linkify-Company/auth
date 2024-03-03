package service

import (
	"auth/internal/config"
	"auth/internal/repository"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"net/smtp"
	"os"
)

type EmailService struct {
	log        logger.Logger
	emailRepos repository.Email
	cfg        config.EmailServiceConfig
}

func NewEmailService(
	log logger.Logger,
	emailRepos repository.Email,
	cfg config.EmailServiceConfig,
) Email {
	return &EmailService{
		log:        log,
		emailRepos: emailRepos,
		cfg:        cfg,
	}
}

func (m *EmailService) Send(ctx context.Context, title string, toEmail string, message string) errify.IError {
	var errCh = make(chan errify.IError, 1)
	var ok = make(chan struct{})

	go func() {
		tlsConfig := &tls.Config{
			ServerName: m.cfg.SmtpServer,
		}

		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", m.cfg.SmtpServer, m.cfg.SmtpPort), tlsConfig)
		if err != nil {
			errCh <- errify.NewInternalServerError(err.Error(), "Send/Dial")
			return
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.cfg.SmtpServer)
		if err != nil {
			errCh <- errify.NewInternalServerError(err.Error(), "Send/PlainAuth")
			return
		}
		defer client.Quit()

		err = client.Auth(smtp.PlainAuth("", os.Getenv(config.AuthEmail), os.Getenv(config.AuthEmailCredentials), m.cfg.SmtpServer))
		if err != nil {
			errCh <- errify.NewInternalServerError(err.Error(), "Send/Auth").SetDetails("Auth error")
			return
		}

		err = client.Mail(os.Getenv(config.AuthEmail))
		if err != nil {
			errCh <- errify.NewInternalServerError(err.Error(), "Send/Auth").SetDetails("Error sender email")
			return
		}

		if err = client.Rcpt(toEmail); err != nil {
			errCh <- errify.NewBadRequestError(err.Error(), ErrInvalidCredentials.Error(), "Send/Rcpt")
			return
		}

		w, err := client.Data()
		if err != nil {
			errCh <- errify.NewInternalServerError(err.Error(), "Send/Data").SetDetails("Error getting writers")
			return
		}
		defer w.Close()

		_, err = w.Write([]byte(fmt.Sprintf("From: %s\r\n", os.Getenv(config.AuthEmail)) +
			fmt.Sprintf("To: %s\r\n", toEmail) +
			fmt.Sprintf("Subject: %s\r\n", title) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"\r\n" +
			message))
		if err != nil {
			errCh <- errify.NewInternalServerError(err.Error(), "Send/Write").SetDetails("Error write message")
			return
		}
		ok <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errify.NewInternalServerError(ctx.Err().Error(), "Send")
	case err := <-errCh:
		return err
	case <-ok:
		m.log.Debugf("Send message successfully")
		return nil
	}
}
