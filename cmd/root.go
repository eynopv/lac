package cmd

import (
	"github.com/eynopv/gorcli/internal"
	"github.com/eynopv/gorcli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gorcli",
		Version: "0.1.1",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			internal.LoadDotEnv()

			if len(VariablesFilePath) > 0 {
				if err := utils.LoadItem(VariablesFilePath, &Variables); err != nil {
					return err
				}
			}

			if len(HeadersFilePath) > 0 {
				if err := utils.LoadItem(HeadersFilePath, &Headers); err != nil {
					return err
				}
			}

			return nil
		},
	}

	VariablesFilePath string
	HeadersFilePath   string
	Verbose           bool
	Variables         map[string]string
	Headers           map[string]string
	Timeout           int
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&VariablesFilePath, "vars", "", "variables file path")
	rootCmd.PersistentFlags().StringVar(&HeadersFilePath, "headers", "", "headers file path")
	rootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 15, "request timeout")
}
