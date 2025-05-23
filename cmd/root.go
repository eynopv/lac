package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/eynopv/lac/pkg/builder"
	"github.com/eynopv/lac/pkg/client"
	"github.com/eynopv/lac/pkg/request"
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
			ClientConfig.PrinterConfig.PrintResponseMeta = strings.Contains(PrintParameters, "m")
			ClientConfig.PrinterConfig.PrintRequestBody = strings.Contains(PrintParameters, "B")
			ClientConfig.PrinterConfig.PrintRequestHeaders = strings.Contains(PrintParameters, "H")
			ClientConfig.PrinterConfig.PrintRequestMeta = strings.Contains(PrintParameters, "M")

			if strings.Contains(PrintParameters, "r") {
				ClientConfig.PrinterConfig.PrintResponseBody = true
				ClientConfig.PrinterConfig.PrintResponseHeaders = true
				ClientConfig.PrinterConfig.PrintResponseMeta = true
			}

			if strings.Contains(PrintParameters, "R") {
				ClientConfig.PrinterConfig.PrintRequestBody = true
				ClientConfig.PrinterConfig.PrintRequestHeaders = true
				ClientConfig.PrinterConfig.PrintRequestMeta = true
			}

			if err := prepareVariables(); err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			b := builder.Builder{
				ClientConfig: ClientConfig,
				TemplatePath: args[0],
				Variables:    Variables,
				Headers:      defaultHeaders,
			}
			runCommandFunction(&b)
		},
	}

	VariablesInput           []string
	EnvironmentFilePathInput string
	PrintParameters          string
	ClientConfig             client.ClientConfig
	Variables                variables.Variables

	defaultHeaders = map[string]request.StringOrStringList{
		"User-Agent": []string{fmt.Sprintf("lac/%s", version)},
	}
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&VariablesInput, "vars", []string{}, "variables")
	rootCmd.PersistentFlags().StringVar(&EnvironmentFilePathInput, "env", ".env", "environment file")
	rootCmd.PersistentFlags().IntVarP(&ClientConfig.Timeout, "timeout", "t", 15, "request timeout")
	rootCmd.PersistentFlags().BoolVar(&ClientConfig.NoRedirects, "no-redirects", false, "do not follow redirects")
	rootCmd.PersistentFlags().StringVar(&PrintParameters, "print", "b",
		"what should be printed in output:\n"+
			" b - response body\n"+
			" h - response headers\n"+
			" m - response meta\n"+
			" r - response body, headers and meta\n"+
			" B - request body\n"+
			" H - request headers\n"+
			" M - request meta\n"+
			" R - request body, headers and meta\n",
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

func loadEnv() {
	err := godotenv.Load(EnvironmentFilePathInput)
	if err != nil && !os.IsNotExist(errors.Unwrap(err)) {
		fmt.Printf("Unable to load environment file: %v\n", err)
	}
}
