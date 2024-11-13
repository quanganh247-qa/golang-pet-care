package device_token

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 9d28896 (image pet)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 9d28896 (image pet)
)

type DeviceTokenController struct {
	service DeviceTokenServiceInterface
}

type DeviceTokenService struct {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	storeDB db.Store
=======
	storeDB    db.Store
	emailQueue *rabbitmq.EmailQueue
>>>>>>> 0fb3f30 (user images)
=======
	storeDB db.Store
>>>>>>> 9d28896 (image pet)
=======
	storeDB    db.Store
	emailQueue *rabbitmq.EmailQueue
>>>>>>> 0fb3f30 (user images)
=======
	storeDB db.Store
>>>>>>> 9d28896 (image pet)
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
