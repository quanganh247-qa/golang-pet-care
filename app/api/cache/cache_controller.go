package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

// ClearAllCache handles the request to clear all cache entries
// @Summary Clear all Redis cache
// @Description Removes all keys from the Redis cache
// @Tags cache
// @Produce json
// @Success 200 {object} map[string]interface{} "Cache cleared successfully"
// @Failure 500 {object} map[string]interface{} "Error clearing cache"
// @Router /api/v1/cache/clear [post]
func ClearAllCache(c *gin.Context) {
	if redis.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis client is not initialized",
		})
		return
	}

	err := redis.Client.ClearAllCache()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clear cache",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache cleared successfully",
	})
}
