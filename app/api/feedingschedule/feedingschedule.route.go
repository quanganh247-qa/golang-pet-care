package feedingschedule

import (
	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

func Routes(router *gin.Engine, store db.Store) {
	feedingScheduleController := &FeedingScheduleController{
		service: &FeedingScheduleService{
			storeDB: store,
		},
	}

	feedingSchedule := router.Group("/feeding-schedule")
	{
		feedingSchedule.POST("/create", feedingScheduleController.CreateFeedingSchedule)
		feedingSchedule.GET("/pet/:pet_id", feedingScheduleController.GetFeedingScheduleByPetID)
		feedingSchedule.GET("/active", feedingScheduleController.ListActiveFeedingSchedules)
		feedingSchedule.PUT("/:feeding_schedule_id", feedingScheduleController.UpdateFeedingSchedule)
		feedingSchedule.DELETE("/:feeding_schedule_id", feedingScheduleController.DeleteFeedingSchedule)
	}
}
