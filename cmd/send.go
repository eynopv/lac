package cmd

import (
	"github.com/spf13/cobra"

	"github.com/eynopv/lac/internal/app"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.Send(args[0])
	},
}
