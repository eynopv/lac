package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/eynopv/gorcli/internal"
	"github.com/eynopv/gorcli/internal/utils"
)

func main() {
	isGorcliDirectoryOrFail()

	showHeaders := flag.Bool("sh", false, "Show response headers")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: gorcli <request> [flags]")
		os.Exit(1)
	}

	internal.LoadDotEnv()

	config, err := internal.LoadConfig()
	if err != nil {
		fmt.Println("Unable to load config")
		os.Exit(1)
	}

	if *showHeaders {
		config.ShowHeaders = *showHeaders
	}

	requestName := args[0]

	if strings.HasPrefix(requestName, "tests/") {
		runTestFlow(requestName, config)
		return
	}

	request, err := internal.NewRequest(requestName, config.Headers, config.Variables)
	if err != nil {
		fmt.Println("Unable to make request:", err)
		return
	}

	result, err := internal.DoRequest(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println("Status:", result.Status)
	fmt.Println("Elapsed Time:", result.ElapsedTime)

	if config.ShowHeaders {
		utils.PrintPrettyJson(result.Headers)
	}

	if result.Body != nil {
		utils.PrintPrettyJson(result.Body)
	}
}

type TestFlow []struct {
	Id     string                `json:"id"`
	Uses   string                `json:"uses"`
	With   map[string]string     `json:"with"`
	Expect *internal.Expectation `json:"expect"`
}

func runTestFlow(name string, config *internal.Config) {
	filePath := fmt.Sprintf("./.gorcli/%s.json", name)
	testFlowContent := utils.LoadFile(filePath)

	if testFlowContent == nil {
		os.Exit(1)
	}

	var testFlow TestFlow
	err := json.Unmarshal(*testFlowContent, &testFlow)

	if err != nil {
		fmt.Printf("Failed to parse test flow %s: %v\n", name, err)
		return
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

		request, err := internal.NewRequest(item.Uses, config.Headers, requestVariables)
		if err != nil {
			fmt.Println("Unable to make request:", err)
			return
		}

		result, err := internal.DoRequest(request)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		if item.Expect != nil {
			err := item.Expect.Check(result)
			if err != nil {
				fmt.Println(err.Error())
				utils.PrintPrettyJson(result.Body)
				os.Exit(1)
			}
		}

		fmt.Println("Status:", result.Status)
		fmt.Println("Elapsed Time:", result.ElapsedTime)

		if config.ShowHeaders {
			utils.PrintPrettyJson(result.Headers)
		}

		if result.Body != nil {
			utils.PrintPrettyJson(result.Body)
		}

		if item.Id != "" {
			maps.Copy(variables, utils.FlattenMap(result.Body, item.Id))
		}
	}
}

func isGorcliDirectoryOrFail() {
	fullPath, _ := utils.FullPath("./.gorcli")
	isGorcliDirectory := utils.FileExists(fullPath)
	if isGorcliDirectory != true {
		fmt.Println("Not a gorcli directory")
		os.Exit(1)
	}
}
