package device_token

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	device_token := routerGroup.RouterDefault.Group("/device-token")
	authRoute := routerGroup.RouterAuth(device_token)
	// DeviceToken.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	deviceTokenApi := &DeviceTokenApi{
		&DeviceTokenController{
			service: &DeviceTokenService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				// emailQueue: rabbitmq.Client.Email,
			},
		},
	}

	{
		// authRoute.POST("/create", DeviceTokenApi.controller.createDeviceToken)
		authRoute.POST("/create", deviceTokenApi.controller.insertDeviceToken)
<<<<<<< HEAD
		authRoute.DELETE("/:token", deviceTokenApi.controller.deleteDeviceToken)
=======
>>>>>>> 0fb3f30 (user images)

	}

}
