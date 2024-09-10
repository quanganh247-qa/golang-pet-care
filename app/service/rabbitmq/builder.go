package rabbitmq

import "fmt"

func (c *ClientMQType) CreateQueue(queueName string) error {
	_, err := c.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error creating queue: %v", err)
	}
	return nil
}

func (c *ClientMQType) CreateQueueAndBindExchange(queueName string) error {
	_, err := c.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error creating queue: %v", err)
	}
	err = c.Channel.QueueBind(queueName, queueName, EmailExchange, false, nil)
	if err != nil {
		return fmt.Errorf("Error binding queue: %v", err)
	}
	return nil
}

func (c *ClientMQType) DeleteQueue(queueName string) error {
	_, err := c.Channel.QueueDelete(queueName, false, false, false)
	if err != nil {
		return fmt.Errorf("Error deleting queue: %v", err)
	}
	return nil
}
