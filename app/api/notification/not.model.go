package notification

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

type NotificationService struct {
	storeDB db.Store
	wsHub   *websocket.Hub // Thêm WebSocket Hub
}

type NotificationController struct {
	service NotificationServiceInterface
}

type NotificationApi struct {
	controller NotificationControllerInterface
}

type NotificationRequest struct {
	Title       string `json:"title" binding:"required" omitempty`
	Content     string `json:"content" binding:"required" omitempty`
	NotifyType  string `json:"notify_type" binding:"required" omitempty`
	RelatedID   int64  `json:"related_id" binding:"required" omitempty`
	RelatedType string `json:"related_type" binding:"required" omitempty`
	DateTime    string `json:"datetime" binding:"required" omitempty`
	Username    string `json:"username" binding:"required" omitempty`
}

type Notification struct {
	Username    string
	Title       string
	Content     string
	NotifyType  string
	RelatedID   int64
	RelatedType string
}

type NotificationResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Title       string `json:"title" omitempty`
	Content     string `json:"content" omitempty`
	DateTime    string `json:"datetime"`
	IsRead      bool   `json:"is_read"`
	NotifyType  string `json:"notify_type" omitempty`
	RelatedID   int32  `json:"related_id" omitempty`
	RelatedType string `json:"related_type" omitempty`
}

// Thêm models cho notification preferences

type NotificationPreference struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Topic     string `json:"topic" omitempty`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at" omitempty`
	UpdatedAt string `json:"updated_at" omitempty`
}

type SubscriptionRequest struct {
	Topics []string `json:"topics" binding:"required" omitempty`
}

type UpdatePreferenceRequest struct {
	Topic   string `json:"topic" binding:"required" omitempty`
	Enabled bool   `json:"enabled" binding:"required" omitempty`
}
