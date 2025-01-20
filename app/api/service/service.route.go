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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> b393bb9 (add service and add permission)
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
<<<<<<< HEAD
=======
		authRoute.POST("/create", SVApi.controller.CreateService)
		authRoute.POST("/delete", SVApi.controller.DeleteService)
		authRoute.GET("/list", SVApi.controller.GetAllServices)
		authRoute.PUT("/update/:serviceid", SVApi.controller.UpdateService)
		authRoute.GET("/getbyid/:serviceid", SVApi.controller.GetServiceByID)
<<<<<<< HEAD
<<<<<<< HEAD
		// authRoute.POST("/create", userApi.controller.createUser)
>>>>>>> cfbe865 (updated service response)
=======
>>>>>>> e30b070 (Get list appoinment by user)
=======
		authRoute.GET("/", SVApi.controller.getAllServices)
>>>>>>> 5e493e4 (get all services)
=======
>>>>>>> b393bb9 (add service and add permission)

	}

}
