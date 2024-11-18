package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
<<<<<<< HEAD
	"github.com/rs/zerolog/log"
)

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPasswordResponse represents the response for forgot password
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// CreateForgotPasswordPayload represents the payload for the async task
type PayloadForgotPassword struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
	OTP      int64  `json:"otp"`
=======
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
>>>>>>> 6610455 (feat: redis queue)
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

<<<<<<< HEAD
	// verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
	// 	Username:   user.Username,
	// 	Email:      user.Email,
	// 	SecretCode: payload.OTP,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create verify email %w", err)
	// }

	subject := "Welcome to Pet Care App - Verify Your Email"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body>
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #333;">Welcome to Pet Care App!</h2>
        <p>Hello %s,</p>
        <p>Thank you for registering with Pet Care App. To verify your email address, please use the following OTP code:</p>
        <div style="background-color: #f5f5f5; padding: 15px; text-align: center; margin: 20px 0;">
            <h3 style="color: #4a90e2; font-size: 24px; margin: 0;">%d</h3>
        </div>
        <p>This code will expire in 5 minutes.</p>
        <p>If you didn't request this verification, please ignore this email.</p>
        <p>Best regards,<br>Pet Care App Team</p>
    </div>
</body>
</html>`, user.Username, payload.OTP)

	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, htmlBody, to, nil, nil, nil)
=======
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomInt(10000000, 999999999),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}
	subject := "Welcome to Pet care app"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8088/api/v1/user/verify_email?email_id=%d&secret_code=%d",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
>>>>>>> 6610455 (feat: redis queue)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}
