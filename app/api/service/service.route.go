package service

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup) {
	SV := routerGroup.RouterDefault.Group("/services")
	authRoute := routerGroup.RouterAuth(SV)
	perRoute := routerGroup.RouterPermission(SV)

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
		// authRoute.POST("/", SVApi.controller.CreateServiceController)
		// authRoute.DELETE("/:id", SVApi.controller.DeleteService)
		authRoute.GET("/", SVApi.controller.GetAllServices)
		// authRoute.PUT("/:id", SVApi.controller.UpdateService)
		authRoute.GET("/:id", SVApi.controller.GetServiceByID)

	}
	{
		// Example: Only users with the "MANAGE_SERVICES" permission can access this route
		perRoute([]perms.Permission{perms.ManageServices}).POST("/", SVApi.controller.CreateServiceController)
		perRoute([]perms.Permission{perms.ManageServices}).DELETE("/:id", SVApi.controller.DeleteService)
		perRoute([]perms.Permission{perms.ManageServices}).PUT("/:id", SVApi.controller.UpdateService)

	}

}
