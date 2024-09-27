package page

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	project := routerGroup.RouterDefault.Group("/page")
	authRoute := routerGroup.RouterAuth(project)
	// user.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	pageAPI := &PageApi{
		&PageController{
			page: &PageService{
				store: db.StoreDB,
			},
		},
	}

	{
		// authRoute.POST("/create", userApi.controller.createUser)
		authRoute.POST("/create", pageAPI.controller.createPage)

	}

}
