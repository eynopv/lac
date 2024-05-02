package cmd

import (
	"fmt"
	"os"

	"github.com/eynopv/lac/internal/client"
	"github.com/eynopv/lac/internal/param"
	"github.com/eynopv/lac/internal/request"
	"github.com/eynopv/lac/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCommandFunction(args, Variables, Headers, Timeout)
	},
}

func runCommandFunction(
	args []string,
	variables map[string]string,
	headers map[string]string,
	timeout int,
) {
	resolvedHeaders := map[string]string{}
	for key, value := range headers {
		resolvedHeaders[key] = param.Param(value).Resolve(variables)
	}

	req, err := request.LoadRequest(args[0])
	req.Headers = utils.CombineMaps(resolvedHeaders, req.Headers)

	if err != nil {
		fmt.Printf("Unable to make request: %v\n", err)
		os.Exit(1)
	}
	req.ResolveParameters(variables)

	requestClient := client.NewClient(timeout)
	result, err := requestClient.Do(req)

	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	result.Print(Verbose)
}
