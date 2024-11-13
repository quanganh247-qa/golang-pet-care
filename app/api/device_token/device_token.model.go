package device_token

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
>>>>>>> 0fb3f30 (user images)
)

type DeviceTokenController struct {
	service DeviceTokenServiceInterface
}

type DeviceTokenService struct {
<<<<<<< HEAD
	storeDB db.Store
=======
	storeDB    db.Store
	emailQueue *rabbitmq.EmailQueue
>>>>>>> 0fb3f30 (user images)
}

// route
type DeviceTokenApi struct {
	controller DeviceTokenControllerInterface
}

type DVTRequest struct {
	Token      string `json:"token"`
	DeviceType string `json:"device_type"`
	LastUsedAt string `json:"last_used_at"`
	ExpiredAt  string `json:"expired_at"`
}

type DVTResponse struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Token      string `json:"token"`
	DeviceType string `json:"device_type"`
	CreatedAt  string `json:"created_at"`
	LastUsedAt string `json:"last_used_at"`
	ExpiredAt  string `json:"expired_at"`
}
