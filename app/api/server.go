package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/connection"
)

type Server struct {
	Router 	*gin.Engine
	Connection *connection.Connection

}

func NewServer(config util.Config) (*Server, error) {
	conn, err := connection.Init(config)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Router: gin.Default(),
	}
	server.SetupRoutes()
	server.Connection = conn
	return server, nil

}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
