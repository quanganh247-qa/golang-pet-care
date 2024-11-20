package notification

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	notification := routerGroup.RouterDefault.Group("/notifications")
	authRoute := routerGroup.RouterAuth(notification)
	// Not.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	petApi := &NotApi{
		&NotController{
			service: &NotService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}
	{
		authRoute.POST("/", petApi.controller.InsertNotification)
		authRoute.GET("/", petApi.controller.ListNotificationsByUsername)
		authRoute.DELETE("/", petApi.controller.DeleteNotification)

	}

}
