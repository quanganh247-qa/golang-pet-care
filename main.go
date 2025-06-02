package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/hibiken/asynq"
	"github.com/quanganh247-qa/go-blog-be/app/api"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/connection"
	_ "github.com/quanganh247-qa/go-blog-be/docs" // Import swagger docs
	"go.uber.org/zap"
)

// @title           Pet Care API
// @version         1.0
// @description     API Server for Pet Care Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.petcare.io/support
// @contact.email  support@petcare.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize database connection
	conn, err := connection.Init(*config)
	if err != nil {
		log.Fatal("cannot initialize connection:", err)
	}
	defer conn.Close()

	// Access the store from the connection
	storeDB := conn.Store

	// Ensure global StoreDB is initialized for package-wide access
	if db.StoreDB == nil {
		db.InitStore(conn.DB)
	}

	// Initialize Redis Task Distributor
	redisOpt := asynq.RedisClientOpt{
		Addr:     config.RedisAddress,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// Initialize WebSocket
	ws := websocket.NewWSClientManager(storeDB)
	go ws.Run()

	// Start the server in a goroutine
	go func() {
		server, err := api.NewServer(*config, taskDistributor, ws, storeDB)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}

		fmt.Printf(color.GreenString("Starting server at %s\n", config.HTTPServerAddress))
		if err := server.Start(config.HTTPServerAddress); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println(color.YellowString("\nShutting down server..."))
	fmt.Println(color.GreenString("Server gracefully stopped"))
}
