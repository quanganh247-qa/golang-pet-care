package connection

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/mail"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type Connection struct {
	Close func()
}

func Init(config util.Config) (*Connection, error) {
	// Initialize JWT token maker
	_, err := token.NewJWTMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}

	// Initialize database connection pool
	connPool, err := pgxpool.New(context.Background(), config.DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	_ = asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	err = redis.InitRedis(config.RedisAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to redis: %w", err)
	}

	DB := db.InitStore(connPool)
	go runTaskProcessor(&config, asynq.RedisClientOpt{Addr: config.RedisAddress}, DB)

	conn := &Connection{
		Close: func() {
			// Close resources when `Close` is called
			connPool.Close()
		},
	}
	return conn, nil
}

func runTaskProcessor(config *util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	// Kiểm tra mailer
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	if mailer == nil {
		log.Fatal("Failed to create mailer")
	}

	// Khởi tạo task processor
	taskProcessor := worker.NewRedisTaskProccessor(redisOpt, store, mailer)
	if taskProcessor == nil {
		log.Fatal("Failed to create task processor")
	}

	// Bắt đầu task processor
	err := taskProcessor.Start()
	if err != nil {
		log.Fatalf("Failed to start task processor: %v", err)
	}
}
