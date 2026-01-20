package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"worker-service/config"
	"worker-service/internal/delivery/workers"
	tasks "worker-service/internal/dto"
	"worker-service/internal/pkg/redis"
	"worker-service/internal/services"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tag = "run-worker"

func NewWorker() *cobra.Command {
	return &cobra.Command{
		Use:     tag,
		Aliases: []string{"worker"},
		Short:   "Run worker",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("Running worker...")

			// Init
			appConfig := config.New()
			redisClient := redis.NewRedisClient[tasks.EmailTask](*appConfig, "email_queue", 0)
			emailService := services.NewEmailService(*appConfig)
			w := workers.NewEmailWorker(redisClient, emailService)

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				<-sigChan
				logrus.Info("Received signal, canceling root context")
				cancel()
				os.Exit(1)
			}()

			if err := w.Run(ctx); err != nil {
				logrus.Error(err)
			}
		},
	}
}
