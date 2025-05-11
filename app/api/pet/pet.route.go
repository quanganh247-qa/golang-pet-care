package pet

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

func Routes(routerGroup middleware.RouterGroup) {
	Pet := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(Pet)
	// Pet.Use(middleware.IPbasedRateLimitingMiddleware())

	// Apply cache middleware to GET endpoints - 5 minute cache duration for main endpoints
	Pet.Use(middleware.CacheMiddleware(time.Minute*5, "pets", []string{"GET"}))

	Pet.Use(middleware.CacheMiddleware(time.Minute*15, "pet_logs", []string{"GET"}))

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
		authRoute.POST("/pet/create", petApi.controller.CreatePet)
		authRoute.GET("/pet/:pet_id", petApi.controller.GetPetByID)
		authRoute.GET("/pets", petApi.controller.ListPets)
		authRoute.GET("/pet", petApi.controller.ListPetsByUsername)
		authRoute.PUT("/pet/:pet_id", petApi.controller.UpdatePet)
		authRoute.DELETE("/pet/:pet_id", petApi.controller.DeletePet)
		authRoute.PUT("/pet/:pet_id/avatar", petApi.controller.UpdatePetAvatar)
		authRoute.GET("/pet/:pet_id/owner", petApi.controller.GetPetOwnerByPetID)

		// Pet logs
		authRoute.GET("/pet/:pet_id/logs", petApi.controller.GetPetLogsByPetID) // Changed from /pet/logs/:pet_id
		authRoute.POST("/pet/logs", petApi.controller.InsertPetLog)
		authRoute.DELETE("/pet/logs/:log_id", petApi.controller.DeletePetLog)
		authRoute.PUT("/pet/logs/:log_id", petApi.controller.UpdatePetLog)
		authRoute.GET("/pet/logs", petApi.controller.GetAllPetLogsByUsername)
		authRoute.GET("/pet/log/:log_id/details", petApi.controller.GetDetailsPetLog) // Changed from /pet/logs/:log_id/details

		// Pet weight tracking
		authRoute.GET("/pet/:pet_id/weights", petApi.controller.GetWeightHistory)
		authRoute.DELETE("/pet/:pet_id/weights/:record_id", petApi.controller.DeleteWeightRecord)
	}
}
