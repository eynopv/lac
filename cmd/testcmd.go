package cmd

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/eynopv/lac/internal"
	"github.com/eynopv/lac/internal/utils"
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
		testCommandFunction(args, Variables, Headers, Timeout)
	},
}

type TestFlow []struct {
	Id     string                `json:"id" yaml:"id"`
	Uses   string                `json:"uses" yaml:"uses"`
	With   map[string]string     `json:"with" yaml:"with"`
	Expect *internal.Expectation `json:"expect" yaml:"expect"`
}

func testCommandFunction(
	args []string,
	variables map[string]string,
	headers map[string]string,
	timeout int,
) {
	testPath := args[0]

	matches, err := findMatchindFiles(testPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	hadErrors := false
	for _, match := range matches {
		if err := loadAndRunTest(match, variables, headers, timeout); err != nil {
			hadErrors = true
		}
	}

	if hadErrors {
		os.Exit(1)
	}
}

func loadAndRunTest(
	testPath string,
	variables map[string]string,
	headers map[string]string,
	timeout int,
) error {
	var testFlow TestFlow
	if err := utils.LoadItem(testPath, &testFlow); err != nil {
		return err
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
			return err
		}

		result, err := internal.DoRequest(request, timeout)
		if err != nil {
			return err
		}

		if item.Expect != nil {
			err := item.Expect.Check(result)
			if err != nil {
				fmt.Println(testPath)
				utils.PrintPrettyJson(result.Body)
				return err
			}
		}

		result.Print(Verbose)

		if item.Id != "" {
			maps.Copy(newVariables, utils.FlattenMap(result.Body, item.Id))
		}
	}

	return nil
}

func findMatchindFiles(pattern string) ([]string, error) {
	var (
		matches []string
		err     error
	)
	if strings.Contains(pattern, "**") {
		parts := strings.SplitN(pattern, "**", 2)
		basePattern := parts[0]
		matchPattern := parts[1]
		err = filepath.Walk(basePattern, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				matched, err := filepath.Match(filepath.Dir(path)+matchPattern, path)
				if err != nil {
					return err
				}
				if matched {
					matches = append(matches, path)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		matches, err = filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
	}
	return matches, nil
}
