package test

import (
	"errors"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/eynopv/gorcli/internal"
	"github.com/eynopv/gorcli/internal/utils"
)

type TestFlow []struct {
	Id     string                `json:"id" yaml:"id"`
	Uses   string                `json:"uses" yaml:"uses"`
	With   map[string]string     `json:"with" yaml:"with"`
	Expect *internal.Expectation `json:"expect" yaml:"expect"`
}

func ExecuteTestCmd() error {
	cmd := flag.NewFlagSet("test", flag.ExitOnError)
	cmd.Usage = printUsage

	cmd.Parse(os.Args[2:])

	args := cmd.Args()

	internal.LoadDotEnv()

	config, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	if len(args) < 1 {
		cmd.Usage()
		os.Exit(0)
	}

	var testFlow TestFlow
	testPath := args[0]
	err = utils.LoadItem(testPath, &testFlow)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to parse test %s: %v\n", testPath, err))
	}

	variables := make(map[string]string)

	for _, item := range testFlow {
		requestVariables := make(map[string]string)
		maps.Copy(requestVariables, config.Variables)

		if item.With != nil {
			withVars := make(map[string]string)
			for key, value := range item.With {
				withVars[key] = internal.ParseParam(value, variables)
			}
			maps.Copy(requestVariables, withVars)
		}

		usesPath := filepath.Join(filepath.Dir(testPath), item.Uses)
		request, err := internal.NewRequest(usesPath, config.Headers, requestVariables)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to make request: %v\n", err))
		}

		result, err := internal.DoRequest(request)
		if err != nil {
			return errors.New(fmt.Sprintf("Error sending request: %v\n", err))
		}

		if item.Expect != nil {
			err := item.Expect.Check(result)
			if err != nil {
				utils.PrintPrettyJson(result.Body)
				return errors.New(err.Error())
			}
		}

		result.Print(config.ShowHeaders)

		if item.Id != "" {
			maps.Copy(variables, utils.FlattenMap(result.Body, item.Id))
		}
	}

	return nil
}

func printUsage() {
	fmt.Println(fmt.Sprintf("Usage: %s test [flags] <test name>", os.Args[0]))
}
