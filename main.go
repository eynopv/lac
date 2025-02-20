package main

import (
	"fmt"

	"github.com/eynopv/lac/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
