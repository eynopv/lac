package internal

import (
	"os"
	"regexp"
)

type Param struct {
	Name  string
	Value string
}

func (p Param) ParseValue() string {
	re := regexp.MustCompile(`\${([^}]+)}`)

	replaced := re.ReplaceAllStringFunc(p.Value, func(match string) string {
		placeholder := match[2 : len(match)-1]

		if value, ok := os.LookupEnv(placeholder); ok {
			return value
		}

		return match
	})

	return replaced
}
