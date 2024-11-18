package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

<<<<<<< HEAD
	"github.com/fatih/color"
	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
=======
	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
>>>>>>> 6610455 (feat: redis queue)
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"go.uber.org/zap"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

// @title Pet Care Management System
// @version 1.0
// @description Pet care management system APIs built with Go using Gin framework
// @host localhost:8088
// @BasePath /api/v1
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

<<<<<<< HEAD
<<<<<<< HEAD
	logger, _ := zap.NewProduction()
	defer logger.Sync()
=======
	server := runGinServer(*config, taskDistributor)
>>>>>>> 6610455 (feat: redis queue)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// Initialize Elasticsearch
	es, err := elasticsearch.NewESService(*config)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to create elasticsearch client: %v\n", err))
	}

	es.CreateIndices()

	server := runGinServer(*config, taskDistributor, es)

=======
	server := runGinServer(*config)
<<<<<<< HEAD
	log.Fatal("run gin server")
>>>>>>> 9d28896 (image pet)
=======

>>>>>>> e01abc5 (pet schedule api)
	defer func() {
		server.Connection.Close()
	}()

}

<<<<<<< HEAD
func runGinServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) *api.Server {
	server, err := api.NewServer(config, taskDistributor, es)
=======
func runGinServer(config util.Config, taskDistributor worker.TaskDistributor) *api.Server {
	server, err := api.NewServer(config, taskDistributor)
>>>>>>> 6610455 (feat: redis queue)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to create server: %v\n", err))
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to start server: %v\n", err))
	}

	return server
}
