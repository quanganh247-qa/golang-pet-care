package products

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	product := routerGroup.RouterDefault.Group("/products")
	authRoute := routerGroup.RouterAuth(product)
	// product.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	petApi := &ProductApi{
		&ProductController{
			service: &ProductService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
		authRoute.GET("/", petApi.controller.GetProducts)
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.POST("/", petApi.controller.CreateProduct)
		authRoute.GET("/:product_id", petApi.controller.GetProductByID)
=======
>>>>>>> bd5945b (get list products)
=======
		authRoute.GET("/:product_id", petApi.controller.GetProductByID)
>>>>>>> 63e2c90 (get product by id)

	}

}
