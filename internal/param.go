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
	return ParseStringParam(p.Value)
}

func ParseStringParam(p string) string {
	re := regexp.MustCompile(`\${([^}]+)}`)

	config := GetContext().Config

	replaced := re.ReplaceAllStringFunc(p, func(match string) string {
		placeholder := match[2 : len(match)-1]

		if value, ok := config.Variables[placeholder]; ok {
			return value
		}

		if value, ok := os.LookupEnv(placeholder); ok {
			return value
		}

		return match
	})

	return replaced
}
