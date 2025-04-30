package smtp

import (
	"time"
)

// CreateSMTPConfigRequest đại diện cho yêu cầu tạo cấu hình SMTP mới
type CreateSMTPConfigRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	SMTPHost  string `json:"smtp_host"`
	SMTPPort  string `json:"smtp_port"`
	IsDefault bool   `json:"is_default"`
}

// SMTPConfigResponse đại diện cho phản hồi thông tin cấu hình SMTP
type SMTPConfigResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	SMTPHost  string    `json:"smtp_host"`
	SMTPPort  string    `json:"smtp_port"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateSMTPConfigRequest đại diện cho yêu cầu cập nhật cấu hình SMTP
type UpdateSMTPConfigRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password"`
	SMTPHost  string `json:"smtp_host"`
	SMTPPort  string `json:"smtp_port"`
	IsDefault bool   `json:"is_default"`
}

// SetDefaultSMTPConfigRequest đại diện cho yêu cầu đặt cấu hình mặc định
type SetDefaultSMTPConfigRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// TestSMTPConfigRequest đại diện cho yêu cầu kiểm tra cấu hình SMTP
type TestSMTPConfigRequest struct {
	// Name      string `json:"name" binding:"required"`
	// Email     string `json:"email" binding:"required,email"`
	// Password  string `json:"password" binding:"required"`
	// SMTPHost  string `json:"smtp_host"`
	// SMTPPort  string `json:"smtp_port"`
	TestEmail string `json:"test_email" binding:"required,email"`
	SMTPId    int64  `json:"smtp_id" binding:"required"`
}

// TestSMTPConfigResponse đại diện cho phản hồi kiểm tra cấu hình SMTP
type TestSMTPConfigResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
