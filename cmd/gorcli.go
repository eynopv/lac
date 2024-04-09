package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/eynopv/gorcli/cmd/run"
	"github.com/eynopv/gorcli/cmd/test"
	"github.com/eynopv/gorcli/internal/utils"
)

var version string
var showVersion bool

func init() {
	version = "0.1.1"

	flag.BoolVar(&showVersion, "v", false, "Print the version")
	flag.Usage = printUsage
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}

	isGorcliDirectoryOrFail()

	var err error

	if len(os.Args) < 2 || (os.Args[1] != "run" && os.Args[1] != "test") {
		flag.Usage()
		return
	}

	if os.Args[1] == "run" {
		err = run.ExecuteRunCmd()
	} else if os.Args[1] == "test" {
		err = test.ExecuteTestCmd()
	}

	if err != nil {
		handleError(err.Error())
	}
}

func isGorcliDirectoryOrFail() {
	fullPath, _ := utils.FullPath("./.gorcli")
	isGorcliDirectory := utils.FileExists(fullPath)
	if !isGorcliDirectory {
		handleError("Not a gorcli directory")
	}
}

func handleError(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func printUsage() {
	fmt.Println(fmt.Sprintf("Usage: %s [flags] <command>", os.Args[0]))
	fmt.Println("Commands:")
	fmt.Println("  run - Perform request")
	fmt.Println("  test - Perform test")
	fmt.Println("Flags:")
	fmt.Println("  -v, --version - Show version")
}
