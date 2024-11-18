package main

import (
	"log"
	"os"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
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
