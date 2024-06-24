package cmd

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/eynopv/lac/pkg/utils"
)

var (
	rootCmd = &cobra.Command{
		Use:     "lac",
		Version: "0.2.5",
		Args:    cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			godotenv.Load(EnvironmentFilePath)

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
					Headers[strings.ToLower(keyValue[0])] = keyValue[1]
				}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCommandFunction(args, Variables, Headers, Timeout)
		},
	}

	VariablesInput      []string
	HeadersInput        []string
	Variables           map[string]string
	Headers             map[string]string
	Timeout             int
	EnvironmentFilePath string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&VariablesInput, "vars", []string{}, "variables")
	rootCmd.PersistentFlags().StringSliceVar(&HeadersInput, "headers", []string{}, "headers")
	rootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 15, "request timeout")
	rootCmd.PersistentFlags().StringVar(&EnvironmentFilePath, "env", ".env", "environment file")
}
