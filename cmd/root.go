package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/utils"
	"github.com/eynopv/lac/pkg/variables"
)

var (
	rootCmd = &cobra.Command{
		Use:     "lac",
		Version: version,
		Args:    cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			loadEnv()

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
	Variables    variables.Variables
	Headers      = map[string]string{
		"User-Agent": fmt.Sprintf("lac/%s", version),
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

			value := keyValue[1]
			if value == "true" {
				Variables[keyValue[0]] = true
			} else if value == "false" {
				Variables[keyValue[0]] = false
			} else if parsedInt, err := strconv.ParseInt(value, 10, 64); err == nil {
				Variables[keyValue[0]] = parsedInt
			} else if parsedFloat, err := strconv.ParseFloat(value, 64); err == nil {
				Variables[keyValue[0]] = parsedFloat
			} else {
				Variables[keyValue[0]] = value
			}
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

			Headers[http.CanonicalHeaderKey(keyValue[0])] = keyValue[1]
		}
	}

	return nil
}

func loadEnv() {
	err := godotenv.Load(EnvironmentFilePathInput)
	if err != nil && !os.IsNotExist(errors.Unwrap(err)) {
		fmt.Printf("Unable to load environment file: %v\n", err)
	}
}
