package variables

import (
	"fmt"
	"regexp"
)

type Resolver struct {
	Vars map[string]any
}

func NewResolver(vars map[string]any) *Resolver {
	if vars == nil {
		vars = map[string]any{}
	}

	return &Resolver{
		Vars: vars,
	}
}

func (r *Resolver) Resolve(input string) string {
	resolve := func(placeholder string, quoted bool) string {
		if value, ok := r.Vars[placeholder]; ok {
			switch v := value.(type) {
			case string:
				return v

			case nil:
				return "null"

			default:
				return fmt.Sprintf("%v", v)
			}
		}

		if quoted {
			return fmt.Sprintf(`"${%v}"`, placeholder)
		}

		return fmt.Sprintf("${%v}", placeholder)
	}

	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	replaced := re.ReplaceAllStringFunc(input, func(match string) string {
		placeholder := match[2 : len(match)-1]
		return resolve(placeholder, false)
	})

	return replaced
}
