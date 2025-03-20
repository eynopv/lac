package param

import (
	"os"
	"regexp"
)

type Param string

func (p Param) Resolve(replacements map[string]string, useEnv bool) string {
	re := regexp.MustCompile(`\${([^}]+)}`)

	replaced := re.ReplaceAllStringFunc(string(p), func(match string) string {
		placeholder := match[2 : len(match)-1]

		if replacements != nil {
			if value, ok := replacements[placeholder]; ok {
				return value
			}
		}

		if useEnv {
			if value, ok := os.LookupEnv(placeholder); ok {
				return value
			}
		}

		return match
	})

	return replaced
}
