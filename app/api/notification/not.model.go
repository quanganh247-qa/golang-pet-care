package notification

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type NotApi struct {
	controller NotControllerInterface
}

type NotController struct {
	service NotServiceInterface
}

type NotService struct {
	storeDB db.Store
}

type NotificationRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DateTime    string `json:"date_time"`
}

type NotificationResponse struct {
	NotificationID int64     `json:"notification_id"`
	Username       string    `json:"username"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	DateTime       time.Time `json:"date_time"`
	IsRead         bool      `json:"is_read"`
}
