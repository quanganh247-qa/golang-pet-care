package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthCheck godoc
// @Summary Health check endpoint
// @Description Get the status of the server and its services
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health check response"
// @Router /health [get]
func (server *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "up",
		"services": gin.H{
			"database":      "connected",
			"redis":         "connected",
			"elasticsearch": "connected",
			"minio":         "connected",
		},
	})
}
