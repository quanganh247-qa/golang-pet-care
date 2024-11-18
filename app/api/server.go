package api

import (
	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
=======
>>>>>>> 6610455 (feat: redis queue)
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/connection"
)

type Server struct {
	Router          *gin.Engine
	Connection      *connection.Connection
	taskDistributor worker.TaskDistributor
<<<<<<< HEAD
	es              *elasticsearch.ESService
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) (*Server, error) {
=======
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
>>>>>>> 6610455 (feat: redis queue)
	conn, err := connection.Init(config)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Router:          gin.Default(),
		Connection:      conn,
		taskDistributor: taskDistributor,
<<<<<<< HEAD
		es:              es,
	}
	server.SetupRoutes(taskDistributor, config, es)
=======
	}
	server.SetupRoutes(taskDistributor)
>>>>>>> 6610455 (feat: redis queue)

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
