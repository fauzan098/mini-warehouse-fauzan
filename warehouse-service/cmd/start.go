package cmd

import (
	// "micro-warehouse/user-service/app"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		// Start the HTTP server logic here
		// app.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}