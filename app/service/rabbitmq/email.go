package rabbitmq

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/mail"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/rabbitmq/amqp091-go"
)

const (
	EmailQueueName  = "dhqanh.email_queue"
	EmailExchange   = "dhqanh.email_exchange"
	EmailRoutingKey = "dhqanh.email_routing_key"

	EmailLogFailed  = "failed"
	EmailLogSuccess = "success"
)

type EmailQueue struct {
	client     *ClientMQType
	queueName  string
	routingKey string
	exchange   string
	store      db.Store
}

type PayloadVerifyEmail struct {
	Username string
}

func (e *EmailQueue) init(c *ClientMQType) {
	e.client = c
	e.exchange = EmailExchange
	e.routingKey = EmailRoutingKey
	e.queueName = EmailQueueName
	e.queueDeclare()
	e.exchangeDeclare()
	e.bindQueue()
}

func (e *EmailQueue) exchangeDeclare() {
	err := e.client.Channel.ExchangeDeclare(e.exchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error creating exchange job:", err)
	}
}

func (e *EmailQueue) queueDeclare() {
	_, err := e.client.Channel.QueueDeclare(e.queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error creating queue job:", err)
	}
}

func (e *EmailQueue) bindQueue() {
	err := e.client.Channel.QueueBind(e.queueName, e.routingKey, e.exchange, false, nil)
	if err != nil {
		log.Fatal("Error binding queue job:", err)
	}
}

func (e *EmailQueue) PublishEmail(data PayloadVerifyEmail) error {
	bodyString, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error marshalling message: %v", err)
	}

	err = e.client.Channel.PublishWithContext(
		e.client.ctx,
		e.exchange,
		e.routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(bodyString),
		})
	if err != nil {
		return fmt.Errorf("Error publishing message: %v", err)
	}
	return nil
}

func (e *EmailQueue) ConsumeEmail() error {
	msgs, err := e.client.Channel.ConsumeWithContext(
		e.client.ctx,
		e.queueName,
		"", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error consuming message: %v", err)
	}
	go func() {
		for d := range msgs {
			err = d.Reject(false)
			if err != nil {
				log.Println("Error nacking message:", err)
			}
			continue

		}

	}()
	return nil
}

func (e *EmailQueue) ConsumeMessage(mailer mail.EmailSender) error {
	msgs, err := e.client.Channel.ConsumeWithContext(
		e.client.ctx,
		e.queueName,
		"",    // consumer tag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("Error consuming message: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var payload PayloadVerifyEmail
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				// Nack and requeue the message
				d.Nack(false, true)
				continue
			}
			log.Printf("COntent a message: %s", payload.Username)

			// Process the task (e.g., send verification email)
			err = e.ProccessTaskSendVerifyEmail(e.client.ctx, payload, mailer)
			if err != nil {
				log.Printf("Failed to process task: %v", err)
				// Nack and requeue the message
				d.Nack(false, true)
				continue
			}

			// Acknowledge successful processing
			d.Ack(false)
		}
	}()

	return nil
}

func (e *EmailQueue) ProccessTaskSendVerifyEmail(ctx context.Context, payload PayloadVerifyEmail, mailer mail.EmailSender) error {
	// var payload PayloadVerifyEmail
	// if err := json.Unmarshal(task.Payload(), &payload); err != nil {
	// 	return fmt.Errorf("failed to unmarshal payload: %w", err)
	// }
	log.Printf("Processing task for user: %s", payload.Username)

	user, err := db.StoreDB.GetUser(ctx, payload.Username)
	if err != nil {
		fmt.Println("User")

		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exists: %w", err)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	log.Printf("User retrieved successfully: %+v", user)

	verifyEmail, err := db.StoreDB.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}
	subject := "Welcome to Simple Bank"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}
