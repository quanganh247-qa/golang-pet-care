package api

import (
	"github.com/gin-gonic/gin"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
=======
>>>>>>> 6610455 (feat: redis queue)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 6610455 (feat: redis queue)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/connection"
)

type Server struct {
	Router          *gin.Engine
	Connection      *connection.Connection
	taskDistributor worker.TaskDistributor
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	es              *elasticsearch.ESService
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) (*Server, error) {
=======
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
>>>>>>> 6610455 (feat: redis queue)
=======
	es              *elasticsearch.ESService
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) (*Server, error) {
>>>>>>> e859654 (Elastic search)
=======
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
>>>>>>> 6610455 (feat: redis queue)
=======
	es              *elasticsearch.ESService
}

func NewServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) (*Server, error) {
>>>>>>> e859654 (Elastic search)
	conn, err := connection.Init(config)
	if err != nil {
		return nil, err
	}
	server := &Server{
		Router:          gin.Default(),
		Connection:      conn,
		taskDistributor: taskDistributor,
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		es:              es,
	}
	server.SetupRoutes(taskDistributor, config, es)
=======
	}
<<<<<<< HEAD
	server.SetupRoutes(taskDistributor)
>>>>>>> 6610455 (feat: redis queue)
=======
	server.SetupRoutes(taskDistributor, config)
>>>>>>> 1a9e82a (reset password api)
=======
		es:              es,
	}
	server.SetupRoutes(taskDistributor, config, es)
>>>>>>> e859654 (Elastic search)
=======
	}
<<<<<<< HEAD
	server.SetupRoutes(taskDistributor)
>>>>>>> 6610455 (feat: redis queue)
=======
	server.SetupRoutes(taskDistributor, config)
>>>>>>> 1a9e82a (reset password api)
=======
		es:              es,
	}
	server.SetupRoutes(taskDistributor, config, es)
>>>>>>> e859654 (Elastic search)

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
