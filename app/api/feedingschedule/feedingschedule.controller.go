package feedingschedule

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedingScheduleController struct {
	service FeedingScheduleServiceInterface
}

func (c *FeedingScheduleController) CreateFeedingSchedule(ctx *gin.Context) {
	var req createFeedingScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.CreateFeedingSchedule(ctx, req.PetID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *FeedingScheduleController) GetFeedingScheduleByPetID(ctx *gin.Context) {
	petID, err := strconv.ParseInt(ctx.Param("pet_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	res, err := c.service.GetFeedingScheduleByPetID(ctx, petID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *FeedingScheduleController) ListActiveFeedingSchedules(ctx *gin.Context) {
	res, err := c.service.ListActiveFeedingSchedules(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *FeedingScheduleController) UpdateFeedingSchedule(ctx *gin.Context) {
	feedingScheduleID, err := strconv.ParseInt(ctx.Param("feeding_schedule_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feeding schedule ID"})
		return
	}

	var req updateFeedingScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.FeedingScheduleID = feedingScheduleID

	if err := c.service.UpdateFeedingSchedule(ctx, feedingScheduleID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Feeding schedule updated successfully"})
}

func (c *FeedingScheduleController) DeleteFeedingSchedule(ctx *gin.Context) {
	feedingScheduleID, err := strconv.ParseInt(ctx.Param("feeding_schedule_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feeding schedule ID"})
		return
	}

	if err := c.service.DeleteFeedingSchedule(ctx, feedingScheduleID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Feeding schedule deleted successfully"})
}
