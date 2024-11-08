package activitylog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ActivityLogControllerInterface interface {
	CreateActivityLog(ctx *gin.Context)
	GetActivityLogByID(ctx *gin.Context)
	ListActivityLogs(ctx *gin.Context)
	UpdateActivityLog(ctx *gin.Context)
	DeleteActivityLog(ctx *gin.Context)
}

func (c *ActivityLogController) CreateActivityLog(ctx *gin.Context) {
	var req createActivityLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.CreateActivityLog(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *ActivityLogController) GetActivityLogByID(ctx *gin.Context) {
	logIDStr := ctx.Param("logid")
	logID, err := strconv.ParseInt(logIDStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid log ID"})
		return
	}

	res, err := c.service.GetActivityLogByID(ctx, logID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *ActivityLogController) UpdateActivityLog(ctx *gin.Context) {
	logID, err := strconv.ParseInt(ctx.Param("logid"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	var req updateActivityLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdateActivityLog(ctx, logID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity log updated successfully"})
}

func (c *ActivityLogController) DeleteActivityLog(ctx *gin.Context) {
	logID, err := strconv.ParseInt(ctx.Param("logid"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	err = c.service.DeleteActivityLog(ctx, logID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity log deleted successfully"})
}

func (c *ActivityLogController) ListActivityLogs(ctx *gin.Context) {
	petIDStr := ctx.Query("petid")
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	res, err := c.service.ListActivityLogs(ctx, petID, int32(limit), int32(offset))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
