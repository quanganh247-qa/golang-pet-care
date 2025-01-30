package petschedule
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e01abc5 (pet schedule api)

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
=======
)

func Routes(routerGroup middleware.RouterGroup) {
>>>>>>> e01abc5 (pet schedule api)
	petSchedule := routerGroup.RouterDefault.Group("/pet-schedule")
	authRoute := routerGroup.RouterAuth(petSchedule)
	// PetSchedule.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	petScheduleApi := &PetScheduleApi{
		&PetScheduleController{
			service: &PetScheduleService{
				storeDB: db.StoreDB, // This should refer to the actual instance
<<<<<<< HEAD
				config:  config,
=======
>>>>>>> e01abc5 (pet schedule api)
			},
		},
	}

	{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.POST("/pet/:pet_id", petScheduleApi.controller.createPetSchedule)
		authRoute.GET("/pet/:pet_id", petScheduleApi.controller.getAllSchedulesByPet)
		authRoute.PUT("/active/:schedule_id", petScheduleApi.controller.activePetSchedule)
		authRoute.GET("/", petScheduleApi.controller.listPetSchedulesByUsername)
		authRoute.DELETE("/:schedule_id", petScheduleApi.controller.deletePetSchedule)
		authRoute.PUT("/:schedule_id", petScheduleApi.controller.updatePetScheduleService)
		petSchedule.POST("/suggestion", petScheduleApi.controller.generateScheduleSuggestion)
	}

}
=======
>>>>>>> 272832d (redis cache)
=======
		authRoute.POST("/create", petScheduleApi.controller.createPetSchedule)
		authRoute.GET("/pet", petScheduleApi.controller.getAllSchedulesByPet)
=======
		authRoute.POST("/pet/:petid", petScheduleApi.controller.createPetSchedule)
		authRoute.GET("/pet/:petid", petScheduleApi.controller.getAllSchedulesByPet)
<<<<<<< HEAD
>>>>>>> 6610455 (feat: redis queue)
		// authRoute.GET("/list", petScheduleApi.controller.ListPetSchedules)
=======
=======
		authRoute.POST("/pet/:pet_id", petScheduleApi.controller.createPetSchedule)
		authRoute.GET("/pet/:pet_id", petScheduleApi.controller.getAllSchedulesByPet)
>>>>>>> 2fe5baf (treatment phase)
		authRoute.PUT("/active/:schedule_id", petScheduleApi.controller.activePetSchedule)
>>>>>>> eb8d761 (updated pet schedule)
		authRoute.GET("/", petScheduleApi.controller.listPetSchedulesByUsername)
		authRoute.DELETE("/:schedule_id", petScheduleApi.controller.deletePetSchedule)
		authRoute.PUT("/:schedule_id", petScheduleApi.controller.updatePetScheduleService)
	}

}
>>>>>>> e01abc5 (pet schedule api)
