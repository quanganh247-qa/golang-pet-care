package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type Server struct {
	Router 	*gin.Engine
}

func NewServer(config util.Config) (*Server, error) {
	server := &Server{
		Router: gin.Default(),
	}
	server.SetupRoutes()
	return server, nil

}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
