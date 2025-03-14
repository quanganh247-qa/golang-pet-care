package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
// Creating to interact with Redis queue
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> ada3717 (Docker file)
type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) *RedisTaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
