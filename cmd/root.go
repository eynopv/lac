package cmd

import (
	"fmt"
	"strings"

	"github.com/eynopv/lac/internal"
	"github.com/eynopv/lac/internal/utils"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "lac",
		Version: "0.2.2",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			internal.LoadDotEnv()

			Variables = map[string]string{}

			for _, variableInput := range VariablesInput {
				err := utils.LoadItem(variableInput, &Variables)
				if err != nil {
					keyValue := strings.Split(variableInput, "=")
					if len(keyValue) != 2 {
						return fmt.Errorf("Invalid variables input: %v", variableInput)
					}
					Variables[keyValue[0]] = keyValue[1]
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

	VariablesInput  []string
	HeadersFilePath string
	Verbose         bool
	Variables       map[string]string
	Headers         map[string]string
	Timeout         int
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringSliceVar(&VariablesInput, "vars", []string{}, "variables")
	rootCmd.PersistentFlags().StringVar(&HeadersFilePath, "headers", "", "headers file path")
	rootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 15, "request timeout")
}
