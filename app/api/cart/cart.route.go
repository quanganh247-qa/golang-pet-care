package cart

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

func Routes(routerGroup middleware.RouterGroup) {
	cart := routerGroup.RouterDefault.Group("/cart")
	authRoute := routerGroup.RouterAuth(cart)
	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	cartApi := &CartApi{
		&CartController{
			service: &CartService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				redis:   redis.Client,
			},
		},
	}

	{
		authRoute.POST("/", cartApi.controller.AddToCart)
	}

}
