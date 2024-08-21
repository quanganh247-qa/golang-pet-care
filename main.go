package main

import (
	"log"

	"github.com/quanganh247-qa/go-blog-be/app/api"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	server := runGinServer(*config)
	defer func() {
		server.Connection.Close()
	}()

}

func runGinServer(config util.Config) *api.Server {
	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	return server
}
