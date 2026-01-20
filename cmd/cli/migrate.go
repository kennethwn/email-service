package cli

import (
	"worker-service/config"
	"worker-service/infrastructure"
	"worker-service/infrastructure/migrations"

	"github.com/spf13/cobra"
)

func NewMigrate() *cobra.Command {
	return &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"migrate"},
		Short:   "run migrations for database",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.New()
			db := infrastructure.InitializeDBConnection(*cfg)
			migrations.MigrateAll(db)
		},
	}
}
