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
	}

}
