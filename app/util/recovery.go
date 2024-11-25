package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryResponse represents the structure of the error response
type RecoveryResponse struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
	Stack     string    `json:"stack,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	RequestID string    `json:"request_id,omitempty"`
}

// Recover returns a Gin middleware for handling panics during request processing.
func Recover(logger *zap.Logger, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}

				// Log the error details
				logger.Error("panic recovered",
					zap.Error(err),
					// zap.String("url", c.Request.URL.String()),
					zap.String("method", c.Request.Method),
					zap.String("client_ip", c.ClientIP()),
				)

				// Provide response to client
				if debug {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal Server Error",
						"details": err.Error(),
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Internal Server Error",
					})
				}

				// Abort to stop further middleware processing
				c.Abort()
			}
		}()

		// Continue processing the request
		c.Next()
	}
}
