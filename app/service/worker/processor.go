package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	mail "github.com/quanganh247-qa/go-blog-be/app/service/mail"

	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProccessor interface {
	Start() error
	ProccessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProccessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProccessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProccessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)
	return &RedisTaskProccessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

func (proccessor *RedisTaskProccessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, proccessor.ProccessTaskSendVerifyEmail)

	return proccessor.server.Start(mux)
}
