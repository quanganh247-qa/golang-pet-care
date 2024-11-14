package petschedule
<<<<<<< HEAD

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
	petSchedule := routerGroup.RouterDefault.Group("/pet-schedule")
	authRoute := routerGroup.RouterAuth(petSchedule)
	// PetSchedule.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	petScheduleApi := &PetScheduleApi{
		&PetScheduleController{
			service: &PetScheduleService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				config:  config,
			},
		},
	}

	{
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
