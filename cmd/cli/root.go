package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "email-worker-service",
	Short: "Email Worker Service",
	Long:  "Email Worker Service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(NewWorker())
	rootCmd.AddCommand(NewApp())
	rootCmd.AddCommand(NewMigrate())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
