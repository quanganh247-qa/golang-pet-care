package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	defaultSMTPAuthAddress   = "smtp.gmail.com"
	defaultSMTPServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	smtpAuthAddress   string
	smtpServerAddress string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		smtpAuthAddress:   defaultSMTPAuthAddress,
		smtpServerAddress: defaultSMTPServerAddress,
	}
}

// NewCustomSMTPSender creates a new email sender with custom SMTP settings
func NewCustomSMTPSender(name string, fromEmailAddress string, fromEmailPassword string, smtpHost string, smtpPort string) EmailSender {
	smtpAuthAddress := smtpHost
	smtpServerAddress := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		smtpAuthAddress:   smtpAuthAddress,
		smtpServerAddress: smtpServerAddress,
	}
}

func (sender *GmailSender) SendEmail(subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, sender.smtpAuthAddress)
	return e.Send(sender.smtpServerAddress, smtpAuth)
}

// send new password to user
func (sender *GmailSender) SendNewPassword(ctx context.Context, email string, newPassword string) error {
	subject := "New Password"
	content := fmt.Sprintf("Your new password is: %s", newPassword)
	to := []string{email}
	return sender.SendEmail(subject, content, to, nil, nil, nil)
}
