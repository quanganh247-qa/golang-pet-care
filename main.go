package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// @title 1View Blog Portal
// @version 1.0
// @description 1View Blog Portal Rest API documentation
// @termsOfService https://swagger.io/terms/

// @contact.name Vivek Singh
// @contact.url https://github.com/san-data-systems
// @contact.email vbhadauriya@redcloudcomputing.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8088
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme.
// @schemes http https

// main is the entry point for the 1View portal API server.
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	server := runGinServer(*config, taskDistributor)

	defer func() {
		server.Connection.Close()
	}()

}

func runGinServer(config util.Config, taskDistributor worker.TaskDistributor) *api.Server {
	server, err := api.NewServer(config, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	return server
}
