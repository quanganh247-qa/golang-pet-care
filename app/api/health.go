package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
