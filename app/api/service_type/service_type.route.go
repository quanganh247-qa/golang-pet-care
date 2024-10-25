package service_type

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	SVT := routerGroup.RouterDefault.Group("/service-type")
	authRoute := routerGroup.RouterAuth(SVT)
	// user.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	SVTApi := &ServiceTypeApi{
		&ServiceTypeController{
			service: &ServiceTypeService{
				storeDB: db.StoreDB, // This should refer to the actual instance

			},
		},
	}

	{
		authRoute.POST("/create", SVTApi.controller.CreateServiceType)
		// authRoute.POST("/create", userApi.controller.createUser)

	}

}
