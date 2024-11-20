package petschedule

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	petSchedule := routerGroup.RouterDefault.Group("/pet-schedule")
	authRoute := routerGroup.RouterAuth(petSchedule)
	// PetSchedule.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	petScheduleApi := &PetScheduleApi{
		&PetScheduleController{
			service: &PetScheduleService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
		authRoute.POST("/pet/:petid", petScheduleApi.controller.createPetSchedule)
		authRoute.GET("/pet/:petid", petScheduleApi.controller.getAllSchedulesByPet)
		authRoute.PUT("/active/:schedule_id", petScheduleApi.controller.activePetSchedule)
		authRoute.GET("/", petScheduleApi.controller.listPetSchedulesByUsername)
		authRoute.DELETE("/:schedule_id", petScheduleApi.controller.deletePetSchedule)
	}

}
