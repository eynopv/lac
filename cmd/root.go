package cmd

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/utils"
)

var (
	version = "0.3.0"

	rootCmd = &cobra.Command{
		Use:     "lac",
		Version: version,
		Args:    cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			godotenv.Load(EnvironmentFilePathInput)

			if err := prepareVariables(); err != nil {
				return err
			}

			if err := prepareHeaders(); err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCommandFunction(args, Variables, Headers, &ClientConfig)
		},
	}

	VariablesInput           []string
	HeadersInput             []string
	EnvironmentFilePathInput string

	ClientConfig client.ClientConfig
	Variables    = map[string]string{}
	Headers      = map[string]string{
		"user-agent": fmt.Sprintf("lac/%s", version),
	}
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&VariablesInput, "vars", []string{}, "variables")
	rootCmd.PersistentFlags().StringSliceVar(&HeadersInput, "headers", []string{}, "headers")
	rootCmd.PersistentFlags().StringVar(&EnvironmentFilePathInput, "env", ".env", "environment file")
	rootCmd.PersistentFlags().IntVarP(&ClientConfig.Timeout, "timeout", "t", 15, "request timeout")
	rootCmd.PersistentFlags().BoolVar(&ClientConfig.NoRedirects, "no-redirects", false, "do not follow redirects")
}

func Execute() error {
	return rootCmd.Execute()
}

func prepareVariables() error {
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
	return nil
}

func prepareHeaders() error {
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
}
