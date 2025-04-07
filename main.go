package main

import (
	"github.com/eynopv/lac/cmd"
)

func main() {
	// cobra handles error
	_ = cmd.Execute()
}
