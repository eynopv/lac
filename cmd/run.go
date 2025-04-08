package cmd

import (
	"fmt"
	"os"

	"github.com/eynopv/lac/pkg/builder"
)

func runCommandFunction(bldr *builder.Builder) {
	req, err := bldr.BuildRequest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	auth, err := bldr.BuildAuth()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := bldr.BuildClient()

	result, err := c.Do(req, auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prntr := bldr.BuildPrinter()
	prntr.Print(result)
}
