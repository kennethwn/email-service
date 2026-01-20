package workers

import (
	"context"
	tasks "worker-service/internal/dto"
	"worker-service/internal/pkg/redis"
	"worker-service/internal/services"

	"github.com/sirupsen/logrus"
)

type EmailWorker struct {
	queue        *redis.RedisClient[tasks.EmailTask]
	emailService services.EmailService
}

func NewEmailWorker(q *redis.RedisClient[tasks.EmailTask], emailService services.EmailService) *EmailWorker {
	return &EmailWorker{
		queue:        q,
		emailService: emailService,
	}
}

func (w *EmailWorker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logrus.Println("email worker shutting down...")
			return nil
		default:
			task, err := w.queue.Dequeue(ctx)
			if err != nil {
				logrus.Error("error dequeuing task: ", err)
				continue
			}
			if err := w.emailService.SendEmail(ctx, task); err != nil {
				logrus.Error("error sending email: ", err)
				continue
			}
		}
	}
}
