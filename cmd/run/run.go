package run

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/eynopv/gorcli/internal"
)

var showHeaders bool

func ExecuteRunCmd() error {
	cmd := flag.NewFlagSet("run", flag.ExitOnError)
	cmd.Usage = printUsage

	cmd.BoolVar(&showHeaders, "sh", false, "Show response headers")
	cmd.Parse(os.Args[2:])

	args := cmd.Args()

	internal.LoadDotEnv()
	config, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	config.ShowHeaders = showHeaders

	if len(args) < 1 {
		cmd.Usage()
		os.Exit(0)
	}

	requestName := args[0]

	request, err := internal.NewRequest(requestName, config.Headers, config.Variables)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to make request: %v\n", err))
	}

	result, err := internal.DoRequest(request)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending request: %v\n", err))
	}

	result.Print(config.ShowHeaders)

	return nil
}

func printUsage() {
	fmt.Println(fmt.Sprintf("Usage: %s run [flags] <request>", os.Args[0]))
	fmt.Println("Flags:")
	fmt.Println("  -sh - Show response headers")
}
