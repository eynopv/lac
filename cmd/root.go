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
		Version: "0.2.5",
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

			Headers = map[string]string{}
			for _, headersInput := range HeadersInput {
				err := utils.LoadItem(headersInput, &Headers)
				if err != nil {
					keyValue := strings.Split(headersInput, "=")
					if len(keyValue) != 2 {
						return fmt.Errorf("Invalid headers input: %v", headersInput)
					}
					Headers[keyValue[0]] = keyValue[1]
				}
			}

			return nil
		},
	}

	VariablesInput []string
	HeadersInput   []string
	Verbose        bool
	Variables      map[string]string
	Headers        map[string]string
	Timeout        int
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringSliceVar(&VariablesInput, "vars", []string{}, "variables")
	rootCmd.PersistentFlags().StringSliceVar(&HeadersInput, "headers", []string{}, "headers")
	rootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 15, "request timeout")
}
