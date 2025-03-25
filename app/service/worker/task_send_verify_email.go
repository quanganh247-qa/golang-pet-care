package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option) error {

	//Convert object to Json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal tasj payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

// func (processor *RedisTaskProccessor) ProccessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
// 	var payload PayloadSendVerifyEmail
// 	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
// 		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
// 	}

// 	user, err := processor.store.GetUser(ctx, payload.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return fmt.Errorf("user doesn't exists: %w", asynq.SkipRetry)
// 		}
// 		return fmt.Errorf("failed to get user: %w", err)
// 	}

// 	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
// 		Username:   user.Username,
// 		Email:      user.Email,
// 		SecretCode: util.RandomInt(10000000, 999999999),
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create verify email %w", err)
// 	}
// 	subject := "Welcome to Pet care app"
// 	// TODO: replace this URL with an environment variable that points to a front-end page
// 	verifyUrl := fmt.Sprintf("http://localhost:8088/api/v1/user/verify_email?email_id=%d&secret_code=%d",
// 		verifyEmail.ID, verifyEmail.SecretCode)
// 	content := fmt.Sprintf(`Hello %s,<br/>
// 	Thank you for registering with us!<br/>
// 	Please <a href="%s">click here</a> to verify your email address.<br/>
// 	`, user.FullName, verifyUrl)
// 	to := []string{user.Email}

// 	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to send verify email: %w", err)
// 	}

// 	return nil
// }

func (processor *RedisTaskProccessor) ProccessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exists: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Tạo secret code với độ dài cố định là 9 chữ số
	secretCode := util.RandomInt(100000000, 999999999)

	// Tạo verify email record với expiration time
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: secretCode,
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// Lấy base URL từ config
	// baseURL := processor.config.BaseURL
	// if baseURL == "" {
	// 	baseURL = ""
	// }

	verifyUrl := fmt.Sprintf("http://localhost:8088/api/v1/user/verify_email?email_id=%d&secret_code=%d",
		verifyEmail.ID, verifyEmail.SecretCode)

	// HTML template cho email
	emailTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            padding: 10px 20px;
            background-color: #4CAF50;
            color: white;
            text-decoration: none;
            border-radius: 5px;
        }
        .footer { margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Welcome to Pet Care App!</h2>
        <p>Hello %s,</p>
        <p>Thank you for registering with us! To complete your registration and verify your email address, please click the button below:</p>
        <p><a href="%s" class="button">Verify Email Address</a></p>
        <p>Or copy and paste this link into your browser:</p>
        <p>%s</p>
        <p>This link will expire in 24 hours.</p>
        <div class="footer">
            <p>If you did not create an account with Pet Care App, please ignore this email.</p>
            <p>Need help? Contact our support team at support@petcareapp.com</p>
        </div>
    </div>
</body>
</html>`

	content := fmt.Sprintf(emailTemplate, user.FullName, verifyUrl, verifyUrl)
	subject := "Verify Your Email - Pet Care App"
	to := []string{user.Email}

	// Gửi email với retry logic
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			// Wait before retrying (với exponential backoff)
			time.Sleep(time.Second * time.Duration(1<<uint(i)))
			continue
		}
		return fmt.Errorf("failed to send verify email after %d attempts: %w", maxRetries, err)
	}

	return nil
}
