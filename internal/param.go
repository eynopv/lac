package internal

import (
	"os"
	"strings"
)

type Param struct {
	Name  string
	Value string
}

func (p Param) ParseValue() string {
	if strings.HasPrefix(p.Value, "$env") {
		parts := strings.Split(p.Value, ".")
		return os.Getenv(strings.Join(parts[1:], "."))
	}
	return p.Value
}
