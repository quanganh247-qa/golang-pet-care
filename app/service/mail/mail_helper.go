package mail

import (
	"context"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// GetEmailSender returns the appropriate email sender based on configuration
func GetEmailSender(config util.Config, store db.Store) EmailSender {
	// Attempt to get default SMTP config from database
	smtpConfig, err := store.GetDefaultSMTPConfig(context.Background())
	if err == nil {
		// Use custom SMTP config if available
		return NewCustomSMTPSender(
			smtpConfig.Name,
			smtpConfig.Email,
			smtpConfig.Password,
			smtpConfig.SmtpHost,
			smtpConfig.SmtpPort,
		)
	}

	return NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)
}
