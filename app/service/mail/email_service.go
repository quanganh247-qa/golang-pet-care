package mail

import (
	"context"
	"fmt"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// EmailService handles sending emails using the configured email sender
type EmailService struct {
	store db.Store
}

// NewEmailService creates a new email service
func NewEmailService(store db.Store) *EmailService {
	return &EmailService{
		store: store,
	}
}

// SendMessageEmail sends a message email
func (s *EmailService) SendMessageEmail(ctx context.Context, config util.Config, to, subject, content, from string) error {

	sender := GetEmailSender(config, s.store)

	// Prepare email recipients
	toList := []string{to}

	// Send the email
	err := sender.SendEmail(
		subject,
		content,
		toList,
		nil, // No CC
		nil, // No BCC
		nil, // No attachments
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
