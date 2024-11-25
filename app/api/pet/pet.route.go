package pet

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
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
			},
		},
	}

	{
		authRoute.POST("/create", petApi.controller.CreatePet)
		authRoute.GET("/:petid", petApi.controller.GetPetByID)
		authRoute.GET("/list", petApi.controller.ListPets)
		authRoute.GET("/", petApi.controller.ListPetsByUsername)
		authRoute.PUT("/:petid", petApi.controller.UpdatePet)
		authRoute.DELETE("/delete/:petid", petApi.controller.DeletePet)
		authRoute.PUT("/avatar/:petid", petApi.controller.UpdatePetAvatar)

		// Pet logs
		authRoute.GET("/logs/:petid", petApi.controller.GetPetLogsByPetID)
		authRoute.POST("/logs", petApi.controller.InsertPetLog)
		authRoute.DELETE("/logs/:logid", petApi.controller.DeletePetLog)
		authRoute.PUT("/logs/:logid", petApi.controller.UpdatePetLog)

	}

}
