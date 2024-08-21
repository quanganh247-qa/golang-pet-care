package api

import "github.com/gin-gonic/gin"


func (server *Server) SetupRoutes() {
	server.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}