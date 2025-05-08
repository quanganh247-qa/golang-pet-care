package products

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	product := routerGroup.RouterDefault.Group("/products")
	authRoute := routerGroup.RouterAuth(product)
	// product.Use(middleware.IPbasedRateLimitingMiddleware())

	// Apply cache middleware to GET endpoints
	product.Use(middleware.CacheMiddleware(time.Minute*5, "products", []string{"GET"}))

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
		authRoute.POST("/", petApi.controller.CreateProduct)
		authRoute.GET("/:product_id", petApi.controller.GetProductByID)

		// Routes for stock management
		authRoute.POST("/:product_id/import", petApi.controller.ImportStock)
		authRoute.POST("/:product_id/export", petApi.controller.ExportStock)
		authRoute.GET("/:product_id/movements", petApi.controller.GetProductStockMovements)
		authRoute.GET("/movements", petApi.controller.GetAllProductStockMovements)

	}

}
