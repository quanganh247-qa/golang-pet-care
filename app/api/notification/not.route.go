package notification

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	notification := routerGroup.RouterDefault.Group("/notification")
	authRoute := routerGroup.RouterAuth(notification)

	// Khởi tạo service với WebSocket Hub
	notificationApi := &NotificationApi{
		&NotificationController{
			service: &NotificationService{
				storeDB: db.StoreDB,
			},
		},
	}

	{
		// Routes hiện tại
		authRoute.POST("/", notificationApi.controller.InsertNotification)
		authRoute.GET("/", notificationApi.controller.ListNotificationsByUsername)
		authRoute.PUT("/:id", notificationApi.controller.MarkAsRead)
		authRoute.POST("/subscribe", notificationApi.controller.SubscribeToNotifications)
		authRoute.POST("/unsubscribe", notificationApi.controller.UnsubscribeFromNotifications)

	}
}
