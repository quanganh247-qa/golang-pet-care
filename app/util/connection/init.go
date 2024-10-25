package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
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
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	// // Initialize RabbitMQ client
	// clientRabbitMQ := rabbitmq.Init(config.RabbitMQAddress)
	// sender := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	// err = clientRabbitMQ.Email.ConsumeMessage(sender)

	// if err != nil {
	// 	log.Println("Error consuming RabbitMQ :", err)

	// }
	// Initialize the database store
	db.InitStore(connPool)

	conn := &Connection{
		Close: func() {
			// Close resources when `Close` is called
			connPool.Close()
			// if err := clientRabbitMQ.Conn.Close(); err != nil {
			// 	log.Println("Error closing RabbitMQ connection:", err)
			// }
			// if err := clientRabbitMQ.Channel.Close(); err != nil {
			// 	log.Println("Error closing RabbitMQ channel:", err)
			// }
			// fmt.Println("Database and RabbitMQ connections closed.")
		},
	}
	return conn, nil
}
