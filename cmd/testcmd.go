package cmd

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/eynopv/gorcli/internal"
	"github.com/eynopv/gorcli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute test",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		testCommandFunction(args, Variables, Headers)
	},
}

type TestFlow []struct {
	Id     string                `json:"id" yaml:"id"`
	Uses   string                `json:"uses" yaml:"uses"`
	With   map[string]string     `json:"with" yaml:"with"`
	Expect *internal.Expectation `json:"expect" yaml:"expect"`
}

func testCommandFunction(args []string, variables map[string]string, headers map[string]string) {
	var testFlow TestFlow
	testPath := args[0]

	if err := utils.LoadItem(testPath, &testFlow); err != nil {
		fmt.Printf("Unable to make test: %v\n", err)
		os.Exit(1)
	}

	newVariables := make(map[string]string)

	for _, item := range testFlow {
		requestVariables := make(map[string]string)
		maps.Copy(requestVariables, variables)

		if item.With != nil {
			withVars := make(map[string]string)
			for key, value := range item.With {
				withVars[key] = internal.ParseParam(value, variables)
			}
			maps.Copy(requestVariables, withVars)
		}

		usesPath := filepath.Join(filepath.Dir(testPath), item.Uses)
		request, err := internal.NewRequest(usesPath, headers, requestVariables)
		if err != nil {
			fmt.Printf("Unable to make request: %v\n", err)
			os.Exit(1)
		}

		result, err := internal.DoRequest(request)
		if err != nil {
			fmt.Printf("Unable to send request: %v\n", err)
			os.Exit(1)
		}

		if item.Expect != nil {
			err := item.Expect.Check(result)
			if err != nil {
				utils.PrintPrettyJson(result.Body)
				fmt.Println(err)
				os.Exit(1)
			}
		}

		result.Print(Verbose)

		if item.Id != "" {
			maps.Copy(newVariables, utils.FlattenMap(result.Body, item.Id))
		}
	}
}
