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
	rootCmd = &cobra.Command{
		Use:     "lac",
		Version: version,
		Args:    cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_ = godotenv.Load(EnvironmentFilePathInput)

			ClientConfig.PrinterConfig.PrintResponseBody = strings.Contains(PrintParameters, "b")
			ClientConfig.PrinterConfig.PrintResponseHeaders = strings.Contains(PrintParameters, "h")
			ClientConfig.PrinterConfig.PrintRequestBody = strings.Contains(PrintParameters, "B")
			ClientConfig.PrinterConfig.PrintRequestHeaders = strings.Contains(PrintParameters, "H")

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
	PrintParameters          string

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
	rootCmd.PersistentFlags().StringVar(&PrintParameters, "print", "b",
		"what should be printed in output:\n"+
			" b - response body\n"+
			" h - response headers\n"+
			" B - request body\n"+
			" H - request headers\n",
	)
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
