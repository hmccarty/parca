package smtp

import (
	"errors"
	"fmt"
	"net/smtp"

	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

var (
	FailedAuthErr = errors.New("couldn't authenticate with smtp server")
)

func NewSMTPClient(config *c.Config) (m.EmailClient, error) {
	auth := smtp.PlainAuth("", config.SMTPEmail, config.SMTPPassword, config.SMTPServer)
	if auth == nil {
		return nil, FailedAuthErr
	}
	return &SMTPClient{
		sender: config.SMTPEmail,
		auth:   auth,
		server: config.SMTPServer,
		port:   config.SMTPPort,
	}, nil
}

type SMTPClient struct {
	sender string
	auth   smtp.Auth
	server string
	port   string
}

func (client *SMTPClient) SendEmail(recipient, subject, content string) error {
	address := fmt.Sprintf("%s:%s", client.server, client.port)
	msg := []byte(fmt.Sprintf("To: %s\r\n Subject: %s\r\n\r\n %s\r\n",
		recipient, subject, content))
	return smtp.SendMail(address, client.auth,
		client.sender, []string{recipient}, msg)
}
