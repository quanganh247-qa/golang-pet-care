package connection

import (
	"context"
	"fmt"
	"log"
<<<<<<< HEAD
	"time"
=======
>>>>>>> 6610455 (feat: redis queue)

<<<<<<< HEAD
	"github.com/fatih/color"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/mail"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/minio"
=======
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> 272832d (redis cache)
=======
>>>>>>> 6610455 (feat: redis queue)
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type Connection struct {
	Close func()
}

func Init(config util.Config) (*Connection, error) {
	if _, err := token.NewJWTMaker(config.SymmetricKey); err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}

<<<<<<< HEAD
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
=======
	// Initialize database connection pool
	connPool, err := pgxpool.New(context.Background(), config.DATABASE_URL)
>>>>>>> 33fcf96 (Big update)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

<<<<<<< HEAD
	// Verify the connection
	if err := connPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}
=======
	_ = asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	err = redis.InitRedis(config.RedisAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to redis: %w", err)
	}

<<<<<<< HEAD
	// // Initialize RabbitMQ client
	// clientRabbitMQ := rabbitmq.Init(config.RabbitMQAddress)
	// sender := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
>>>>>>> 272832d (redis cache)

	// Khởi tạo Redis
	go redis.InitRedis(config.RedisAddress)
	go runTaskProcessor(&config, asynq.RedisClientOpt{Addr: config.RedisAddress}, db.InitStore(connPool))
	go initMinio(&config)
=======
	// _ = db.InitStore(connPool)

	DB := db.InitStore(connPool)
	fmt.Println("DB")
	go runTaskProcessor(&config, asynq.RedisClientOpt{Addr: config.RedisAddress}, DB)
>>>>>>> 6610455 (feat: redis queue)

	conn := &Connection{
		Close: func() {
			connPool.Close()
		},
	}
	return conn, nil
}

func runTaskProcessor(config *util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	log.Println("Starting task processor...")
=======
	// log.Println("Starting task processor...")
>>>>>>> 1f24c18 (feat: OTP with redis)

>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> eb8d761 (updated pet schedule)
	// Kiểm tra mailer
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	if mailer == nil {
		log.Fatal("Failed to create mailer")
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	log.Println("Mailer initialized")
>>>>>>> 6610455 (feat: redis queue)
=======
	// log.Println("Mailer initialized")
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
>>>>>>> eb8d761 (updated pet schedule)

	// Khởi tạo task processor
	taskProcessor := worker.NewRedisTaskProccessor(redisOpt, store, mailer)
	if taskProcessor == nil {
		log.Fatal("Failed to create task processor")
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	log.Println("Task processor initialized")
>>>>>>> 6610455 (feat: redis queue)
=======
	// log.Println("Task processor initialized")
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
>>>>>>> eb8d761 (updated pet schedule)

	// Bắt đầu task processor
	err := taskProcessor.Start()
	if err != nil {
		log.Fatalf("Failed to start task processor: %v", err)
	}
}
<<<<<<< HEAD

func initMinio(config *util.Config) {
	client, err := minio.NewMinIOClient(
		config.MinIOEndpoint,
		config.MinIOAccessKey,
		config.MinIOSecretKey,
		config.MinIOSSL,
	)
	if err != nil {
		log.Println(color.RedString("Failed to initialize MinIO client: %v", err))
	}

	if err := client.CheckConnection(context.Background()); err != nil {
		log.Println(color.RedString("Failed to check connection: %v", err))
	}

}
=======
>>>>>>> 6610455 (feat: redis queue)
