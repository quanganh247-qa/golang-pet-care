package cart

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

func Routes(routerGroup middleware.RouterGroup) {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	cart := routerGroup.RouterDefault.Group("/")
=======
	cart := routerGroup.RouterDefault.Group("/cart")
>>>>>>> c449ffc (feat: cart api)
=======
	cart := routerGroup.RouterDefault.Group("/")
>>>>>>> dc47646 (Optimize SQL query)
=======
	cart := routerGroup.RouterDefault.Group("/cart")
>>>>>>> c449ffc (feat: cart api)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> dc47646 (Optimize SQL query)
		authRoute.POST("/cart", cartApi.controller.AddToCart)
		authRoute.GET("/cart", cartApi.controller.GetCartItems)
		authRoute.DELETE("/cart/product/:product_id", cartApi.controller.RemoveItemFromCart)

	}
	{
		authRoute.POST("/order", cartApi.controller.CreateOrder)
		authRoute.GET("/order", cartApi.controller.GetOrders)
		authRoute.GET("/order/:order_id", cartApi.controller.GetOrderByID)
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.GET("/orders", cartApi.controller.GetAllOrders)
=======
		authRoute.POST("/", cartApi.controller.AddToCart)
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> c449ffc (feat: cart api)
=======
		authRoute.GET("/", cartApi.controller.GetCartItems)
		authRoute.POST("/order", cartApi.controller.CreateOrder)
<<<<<<< HEAD

>>>>>>> 21608b5 (cart and order api)
=======
		authRoute.GET("/order", cartApi.controller.GetOrders)
		authRoute.GET("/order/:order_id", cartApi.controller.GetOrderByID)
>>>>>>> b0fe977 (place order and make payment)
=======
		authRoute.DELETE("/product/:product_id", cartApi.controller.RemoveItemFromCart)
>>>>>>> 4a16bfc (remove item in cart)
=======
		authRoute.GET("/orders", cartApi.controller.GetAllOrders)
>>>>>>> dc47646 (Optimize SQL query)
=======
		authRoute.POST("/", cartApi.controller.AddToCart)
>>>>>>> c449ffc (feat: cart api)
=======
		authRoute.GET("/", cartApi.controller.GetCartItems)
		authRoute.POST("/order", cartApi.controller.CreateOrder)
<<<<<<< HEAD

>>>>>>> 21608b5 (cart and order api)
=======
		authRoute.GET("/order", cartApi.controller.GetOrders)
		authRoute.GET("/order/:order_id", cartApi.controller.GetOrderByID)
>>>>>>> b0fe977 (place order and make payment)
	}

}
