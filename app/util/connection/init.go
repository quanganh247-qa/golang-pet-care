package connection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type Connection struct {
	Close func()
	Store db.Store
	DB    *pgxpool.Pool
}

func Init(config util.Config) (*Connection, error) {
	if _, err := token.NewJWTMaker(config.SymmetricKey); err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}

	// Configure connection pool with reasonable defaults
	poolConfig, err := pgxpool.ParseConfig(config.DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Set some reasonable pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	// Verify the connection
	if err := connPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}

	// Initialize store
	store := db.InitStore(connPool)

	// Initialize Redis with debug mode
	go redis.InitRedis(config)
	go runTaskProcessor(&config, asynq.RedisClientOpt{
		Addr: config.RedisAddress,
		// Username: config.RedisUsername,
		// Password: config.RedisPassword,
	}, store)
	// go initMinio(&config)

	conn := &Connection{
		Close: func() {
			connPool.Close()
		},
		Store: store,
		DB:    connPool,
	}
	return conn, nil
}

func runTaskProcessor(config *util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	// Khởi tạo task processor với cấu hình
	taskProcessor := worker.NewRedisTaskProccessor(redisOpt, store, *config)
	if taskProcessor == nil {
		log.Fatal("Failed to create task processor")
	}

	// Bắt đầu task processor
	err := taskProcessor.Start()
	if err != nil {
		log.Fatalf("Failed to start task processor: %v", err)
	}
}

// func initMinio(config *util.Config) {
// 	client, err := minio.NewMinIOClient(
// 		config.MinIOEndpoint,
// 		config.MinIOAccessKey,
// 		config.MinIOSecretKey,
// 		config.MinIOSSL,
// 	)
// 	if err != nil {
// 		log.Println(color.RedString("Failed to initialize MinIO client: %v", err))
// 	}

// 	if err := client.CheckConnection(context.Background()); err != nil {
// 		log.Println(color.RedString("Failed to check connection: %v", err))
// 	}

// }
