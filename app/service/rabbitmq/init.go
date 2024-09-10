package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClientMQType struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	ctx     context.Context
	Email   *EmailQueue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

//"amqp://guest:guest@localhost:5672/"

var Client = &ClientMQType{}

func Init(url string) *ClientMQType {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	Client.Conn = conn
	Client.Channel = ch
	Client.ctx = context.WithValue(context.Background(), "appId", uuid.New().String())

	Client.Email = &EmailQueue{}
	Client.Email.init(Client)

	return Client

}

func handleErrorConsumer(d *amqp.Delivery, err error) {
	log.Println(err)
	err = d.Nack(false, true)
	if err != nil {
		log.Println("Error nacking message:", err)
	}
}

func (c *EmailQueue) ConsumerByQueueAndRouting(queueName string, routingKey string, handler func(PayloadVerifyEmail) (interface{}, error)) error {
	msgs, err := c.client.Channel.ConsumeWithContext(c.client.ctx, queueName, routingKey, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error consuming message: %v", err)
	}
	go func() {
		for d := range msgs {
			var out PayloadVerifyEmail
			if err := json.Unmarshal(d.Body, &out); err != nil {
				log.Println("Error unmarshalling message:", err)
				err = d.Reject(false)
				if err != nil {
					log.Println("Error nacking message:", err)
				}
				continue
			}
			_, err = handler(out)
			if err != nil {
				continue
			}
		}
	}()
	log.Println("AMQP Listening consumer job")
	return nil
}
