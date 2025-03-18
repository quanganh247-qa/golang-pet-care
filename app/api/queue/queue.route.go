package queue

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	queue := routerGroup.RouterDefault.Group("/queue")
	authRoute := routerGroup.RouterAuth(queue)

	// Initialize API
	queueApi := &QueueApi{
		&QueueController{
			service: &QueueService{
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.GET("/", queueApi.controller.GetQueue)
		authRoute.PUT("/:id/status", queueApi.controller.UpdateQueueItemStatus)
	}
}
