package pet

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

func Routes(routerGroup middleware.RouterGroup) {
	Pet := routerGroup.RouterDefault.Group("/pet")
	authRoute := routerGroup.RouterAuth(Pet)
	// Pet.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	petApi := &PetApi{
		&PetController{
			service: &PetService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				redis:   redis.Client,
			},
		},
	}

	{
		authRoute.POST("/create", petApi.controller.CreatePet)
		authRoute.GET("/:pet_id", petApi.controller.GetPetByID)
		authRoute.GET("/list", petApi.controller.ListPets)
		authRoute.GET("/", petApi.controller.ListPetsByUsername)
		authRoute.PUT("/:pet_id", petApi.controller.UpdatePet)
		authRoute.DELETE("/delete/:pet_id", petApi.controller.DeletePet)
		authRoute.PUT("/avatar/:pet_id", petApi.controller.UpdatePetAvatar)

		// Pet logs
		authRoute.GET("/logs/:pet_id", petApi.controller.GetPetLogsByPetID)
		authRoute.POST("/logs", petApi.controller.InsertPetLog)
		authRoute.DELETE("/logs/:log_id", petApi.controller.DeletePetLog)
		authRoute.PUT("/logs/:log_id", petApi.controller.UpdatePetLog)

		// Pet summary
		Pet.GET("/summary/:pet_id", petApi.controller.GetPetProfileSummary)
	}

}
