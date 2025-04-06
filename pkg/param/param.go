package param

import (
	"fmt"
	"os"
	"regexp"
)

type Param string

func (p Param) Resolve(replacements map[string]any, useEnv bool) string {
	resolve := func(placeholder string, quoted bool) string {
		if replacements != nil {
			if value, ok := replacements[placeholder]; ok {
				switch v := value.(type) {
				case string:
					if quoted {
						return fmt.Sprintf(`"%v"`, v)
					}

					return v

				case nil:
					return "null"

				default:
					return fmt.Sprintf("%v", v)
				}
			}
		}

		if useEnv {
			if value, ok := os.LookupEnv(placeholder); ok {
				if quoted {
					return fmt.Sprintf(`"%v"`, value)
				}

				return value
			}
		}

		if quoted {
			return fmt.Sprintf(`"${%v}"`, placeholder)
		}

		return fmt.Sprintf("${%v}", placeholder)
	}

	reWithQuotes := regexp.MustCompile(`"\$\{([^}]+)\}"`)
	replaced := reWithQuotes.ReplaceAllStringFunc(string(p), func(match string) string {
		placeholder := match[3 : len(match)-2]
		return resolve(placeholder, true)
	})

	reWithoutQuotes := regexp.MustCompile(`\$\{([^}]+)\}`)
	replaced = reWithoutQuotes.ReplaceAllStringFunc(replaced, func(match string) string {
		placeholder := match[2 : len(match)-1]
		return resolve(placeholder, false)
	})

	return replaced
}
