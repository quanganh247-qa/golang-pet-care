package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()

		// Execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := c.Request.Method

		// Request route
		reqUri := c.Request.RequestURI

		// Status code
		statusCode := c.Writer.Status()

		// Request IP
		clientIP := c.ClientIP()

		// Get cache status from header (if available)
		cacheStatus := c.Writer.Header().Get("X-Cache")
		if cacheStatus == "" {
			cacheStatus = "N/A"
		}

		// Get more detailed cache info from context if available
		cacheDetails := ""
		if status, exists := c.Get("cache_status"); exists {
			source, _ := c.Get("cache_source")
			cacheDetails = fmt.Sprintf(" [%s from %s]", status, source)
		}

		log.Printf("| %3d | %13v | %15s | %6s | %s | CACHE: %s%s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			cacheStatus,
			cacheDetails,
		)
	}
}
