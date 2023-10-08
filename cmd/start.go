package cmd

import (
	"github.com/spf13/cobra"
)

// startCmd starts this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the governor api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startServer() error {
	logger.Debug("Hello!")
	return nil
}
