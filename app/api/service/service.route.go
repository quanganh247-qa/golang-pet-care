package service

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	SV := routerGroup.RouterDefault.Group("/service")
	authRoute := routerGroup.RouterAuth(SV)
	// user.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	SVApi := &ServiceApi{
		&ServiceController{
			service: &ServiceService{
				storeDB: db.StoreDB, // This should refer to the actual instance

			},
		},
	}

	{
		authRoute.POST("/create", SVApi.controller.CreateService)
		// authRoute.POST("/create", userApi.controller.createUser)

	}

}
