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
	return ParseStringParam(p.Value, nil)
}

func ParseStringParam(p string, replacements map[string]string) string {
	re := regexp.MustCompile(`\${([^}]+)}`)

	replaced := re.ReplaceAllStringFunc(p, func(match string) string {
		placeholder := match[2 : len(match)-1]

		if replacements != nil {
			if value, ok := replacements[placeholder]; ok {
				return value
			}
		}

		if value, ok := os.LookupEnv(placeholder); ok {
			return value
		}

		return match
	})

	return replaced
}
