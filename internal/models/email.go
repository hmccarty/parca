package models

type EmailClient interface {
	SendEmail(recipient, subject, content string) error
}
