package queue

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueueControllerInterface interface {
	GetQueue(ctx *gin.Context)
	UpdateQueueItemStatus(ctx *gin.Context)
}

func (c *QueueController) GetQueue(ctx *gin.Context) {
	queueItems, err := c.service.GetQueueService(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, queueItems)
}

func (c *QueueController) UpdateQueueItemStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdateQueueItemStatusService(ctx, id, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
