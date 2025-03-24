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
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.GET("/:pet_id", petApi.controller.GetPetByID)
=======
		authRoute.GET("/:petid", petApi.controller.GetPetByID)
>>>>>>> 7a9ad08 (updated pet api)
		authRoute.GET("/list", petApi.controller.ListPets)
		authRoute.GET("/", petApi.controller.ListPetsByUsername)
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.PUT("/:pet_id", petApi.controller.UpdatePet)
		authRoute.DELETE("/delete/:pet_id", petApi.controller.DeletePet)
		authRoute.PUT("/avatar/:pet_id", petApi.controller.UpdatePetAvatar)

		// Pet logs
		authRoute.GET("/logs/:pet_id", petApi.controller.GetPetLogsByPetID)
		authRoute.POST("/logs", petApi.controller.InsertPetLog)
		authRoute.DELETE("/logs/:log_id", petApi.controller.DeletePetLog)
		authRoute.PUT("/logs/:log_id", petApi.controller.UpdatePetLog)

<<<<<<< HEAD
<<<<<<< HEAD
=======
		authRoute.PUT("/update/:petid", petApi.controller.UpdatePet)
=======
		authRoute.PUT("/:petid", petApi.controller.UpdatePet)
>>>>>>> 5ea33aa (PUT pet info)
		authRoute.DELETE("/delete/:petid", petApi.controller.DeletePet)
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> c73e2dc (pagination function)
=======
=======
		authRoute.PUT("/avatar/:petid", petApi.controller.UpdatePetAvatar)
>>>>>>> e30b070 (Get list appoinment by user)

		// Pet logs
		authRoute.GET("/logs/:petid", petApi.controller.GetPetLogsByPetID)
<<<<<<< HEAD
>>>>>>> 7e616af (add pet log schema)
=======
=======
		authRoute.GET("/:pet_id", petApi.controller.GetPetByID)
		authRoute.GET("/list", petApi.controller.ListPets)
		authRoute.GET("/", petApi.controller.ListPetsByUsername)
		authRoute.PUT("/:pet_id", petApi.controller.UpdatePet)
		authRoute.DELETE("/delete/:pet_id", petApi.controller.DeletePet)
		authRoute.PUT("/avatar/:pet_id", petApi.controller.UpdatePetAvatar)

		// Pet logs
		authRoute.GET("/logs/:pet_id", petApi.controller.GetPetLogsByPetID)
>>>>>>> 2fe5baf (treatment phase)
		authRoute.POST("/logs", petApi.controller.InsertPetLog)
		authRoute.DELETE("/logs/:log_id", petApi.controller.DeletePetLog)
		authRoute.PUT("/logs/:log_id", petApi.controller.UpdatePetLog)

>>>>>>> 3835eb4 (update pet_schedule api)
=======
		// Pet summary
		Pet.GET("/summary/:pet_id", petApi.controller.GetPetProfileSummary)
>>>>>>> ffc9071 (AI suggestion)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
	}

}
