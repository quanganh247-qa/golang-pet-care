package component

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	component := routerGroup.RouterDefault.Group("/component")
	authRoute := routerGroup.RouterAuth(component)

	componentAPI := &ComponentApi{
		&ComponentController{
			service: &ComponentService{
				store: db.StoreDB,
			},
		},
	}

	{
		authRoute.POST("/create", componentAPI.controller.createComponent)
		authRoute.GET("/:id", componentAPI.controller.getComponentByID)
		authRoute.GET("/", componentAPI.controller.getComponentsByUser)
		authRoute.DELETE("/:id", componentAPI.controller.removeComponent)

	}
}
