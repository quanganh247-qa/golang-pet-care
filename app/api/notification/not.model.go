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
	PetID            int64     `json:"pet_id"`
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	DueDate          time.Time `json:"due_date"`
	RepeatInterval   string    `json:"repeat_interval"`
	NotificationSent bool      `json:"notification_sent"`
}

type NotificationResponse struct {
	NotificationID   int64     `json:"notification_id"`
	PetID            int64     `json:"pet_id"`
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	DueDate          time.Time `json:"due_date"`
	RepeatInterval   string    `json:"repeat_interval"`
	IsCompleted      bool      `json:"is_completed"`
	NotificationSent bool      `json:"notification_sent"`
}
