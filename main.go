package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"go.uber.org/zap"
)

// main is the entry point for the 1View portal API server.
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

	// Initialize Elasticsearch
	es, err := elasticsearch.NewESService(*config)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to create elasticsearch client: %v\n", err))
	}

	es.CreateIndices()

	server := runGinServer(*config, taskDistributor, es)

	defer func() {
		server.Connection.Close()
	}()

}

func runGinServer(config util.Config, taskDistributor worker.TaskDistributor, es *elasticsearch.ESService) *api.Server {
	server, err := api.NewServer(config, taskDistributor, es)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to create server: %v\n", err))
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		fmt.Printf(color.RedString("❌ ERROR: Failed to start server: %v\n", err))
	}

	return server
}
