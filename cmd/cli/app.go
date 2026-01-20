package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"worker-service/config"
	"worker-service/infrastructure"
	"worker-service/internal/delivery/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewApp() *cobra.Command {
	return &cobra.Command{
		Use:     "run-app",
		Aliases: []string{"app"},
		Short:   "Run App",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("Running app...")

			// Init
			appConfig := config.New()
			db := infrastructure.InitializeDBConnection(*appConfig)
			router := http.InitRoutes(db)

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			_, cancel := context.WithCancel(context.Background())
			go func() {
				<-sigChan
				logrus.Info("Received signal, canceling root context")
				cancel()
				os.Exit(1)
			}()

			if err := router.Run(":8080"); err != nil {
				logrus.Error(err)
			}
		},
	}
}
