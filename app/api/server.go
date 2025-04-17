package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type Server struct {
	Router          *gin.Engine
	taskDistributor worker.TaskDistributor
	ws              *websocket.WSClientManager
	store           db.Store
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor, ws *websocket.WSClientManager, store db.Store) (*Server, error) {
	server := &Server{
		Router:          gin.Default(),
		taskDistributor: taskDistributor,
		ws:              ws,
		store:           store,
	}

	server.SetupRoutes(taskDistributor, config, ws)

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
