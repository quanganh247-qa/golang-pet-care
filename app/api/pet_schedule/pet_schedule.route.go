package petschedule

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
	petSchedule := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(petSchedule)
	// PetSchedule.Use(middleware.IPbasedRateLimitingMiddleware())

	// Apply cache middleware to GET endpoints - 5 minute cache duration
	petSchedule.Use(middleware.CacheMiddleware(time.Minute*5, "pet_schedules", []string{"GET"}))

	// Khoi tao api
	petScheduleApi := &PetScheduleApi{
		&PetScheduleController{
			service: &PetScheduleService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				config:  config,
				redis:   redis.Client,
			},
		},
	}

	{
		// Pet schedule management
		authRoute.POST("/schedules", petScheduleApi.controller.createPetSchedule)
		authRoute.GET("/schedules/pet/:petid", petScheduleApi.controller.getAllSchedulesByPet)
		authRoute.GET("/schedules", petScheduleApi.controller.listPetSchedulesByUsername)
		authRoute.PUT("/schedules/:schedule_id", petScheduleApi.controller.updatePetScheduleService)
		authRoute.DELETE("/schedules/:schedule_id", petScheduleApi.controller.deletePetSchedule)

		// Schedule status management
		authRoute.PUT("/schedules/:schedule_id/activate", petScheduleApi.controller.activePetSchedule)

		// Schedule suggestions
		petSchedule.POST("/schedules/suggestion", petScheduleApi.controller.generateScheduleSuggestion)
	}

}
