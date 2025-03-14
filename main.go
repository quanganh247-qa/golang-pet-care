package main

import (
	"fmt"
	"log"

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
	"github.com/fatih/color"
	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
<<<<<<< HEAD
=======
	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"go.uber.org/zap"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	logger, _ := zap.NewProduction()
	defer logger.Sync()
=======
	server := runGinServer(*config, taskDistributor)
>>>>>>> 6610455 (feat: redis queue)
=======
	// Initialize Elasticsearch
	es, err := elasticsearch.NewESService(*config)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to create elasticsearch client: %v\n", err))
	}

	es.CreateIndices()

	server := runGinServer(*config, taskDistributor, es)
>>>>>>> e859654 (Elastic search)

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
<<<<<<< HEAD
func runGinServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) *api.Server {
	server, err := api.NewServer(config, taskDistributor, es)
=======
func runGinServer(config util.Config, taskDistributor worker.TaskDistributor) *api.Server {
	server, err := api.NewServer(config, taskDistributor)
>>>>>>> 6610455 (feat: redis queue)
=======
func runGinServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) *api.Server {
	server, err := api.NewServer(config, taskDistributor, es)
>>>>>>> e859654 (Elastic search)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to create server: %v\n", err))
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to start server: %v\n", err))
	}

	return server
}
