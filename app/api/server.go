package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/connection"
)

type Server struct {
	Router          *gin.Engine
	Connection      *connection.Connection
	taskDistributor worker.TaskDistributor
	es              *elasticsearch.ESService
	ws              *websocket.WSClientManager
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService, ws *websocket.WSClientManager) (*Server, error) {
	conn, err := connection.Init(config)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Router:          gin.Default(),
		Connection:      conn,
		taskDistributor: taskDistributor,
		es:              es,
		ws:              ws,
	}

	server.SetupRoutes(taskDistributor, config, es, ws)

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
