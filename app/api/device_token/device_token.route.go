package device_token

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	device_token := routerGroup.RouterDefault.Group("")
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

	// Public routes
	device_token.GET("/device-tokens", deviceTokenApi.controller.getDeviceTokens)

	{
		// authRoute.POST("/create", DeviceTokenApi.controller.createDeviceToken)
		authRoute.POST("/device-token/create", deviceTokenApi.controller.insertDeviceToken)
		authRoute.DELETE("/device-token/:token", deviceTokenApi.controller.deleteDeviceToken)

	}

}
