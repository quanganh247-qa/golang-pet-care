package project

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	project := routerGroup.RouterDefault.Group("/project")
	authRoute := routerGroup.RouterAuth(project)
	// user.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	projectAPI := &ProjectApi{
		&ProjectController{
			service: &ProjectService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
		// authRoute.POST("/create", userApi.controller.createUser)
		authRoute.POST("/create", projectAPI.controller.createProject)

	}

}
