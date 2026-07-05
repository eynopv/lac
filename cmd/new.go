package cmd

import (
	"github.com/spf13/cobra"

	"github.com/eynopv/lac/internal/app"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new request",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.New(args[0])
	},
	Args: cobra.ExactArgs(1),
}
