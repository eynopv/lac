package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "lac",
		Version: version,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
