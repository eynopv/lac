package cmd

import (
	"fmt"
	"os"

	"github.com/eynopv/gorcli/internal"
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
	requestName := args[0]

	request, err := internal.NewRequest(requestName, headers, variables)

	if err != nil {
		fmt.Printf("Unable to make request: %v\n", err)
		os.Exit(1)
	}

	result, err := internal.DoRequest(request, timeout)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}

	result.Print(Verbose)
}
