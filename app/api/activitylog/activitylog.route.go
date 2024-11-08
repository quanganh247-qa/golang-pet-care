package activitylog

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	activityLog := routerGroup.RouterDefault.Group("/activity-log")
	authRoute := routerGroup.RouterAuth(activityLog)

	activityLogApi := &ActivityLogApi{
		&ActivityLogController{
			service: &ActivityLogService{
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.POST("/create", activityLogApi.controller.CreateActivityLog)
		authRoute.GET("/get/:logid", activityLogApi.controller.GetActivityLogByID)
		authRoute.GET("/list", activityLogApi.controller.ListActivityLogs)
		authRoute.PUT("/update/:logid", activityLogApi.controller.UpdateActivityLog)
		authRoute.DELETE("/delete/:logid", activityLogApi.controller.DeleteActivityLog)
	}
}
