package printer

import "fmt"

type Color = string

const (
	ColorReset   Color = "\033[0m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
)

func Colorize(s string, c Color) string {
	return fmt.Sprintf("%v%s%v", c, s, ColorReset)
}

func Red(s string) string {
	return ColorRed + s + ColorReset
}

func Green(s string) string {
	return ColorGreen + s + ColorReset
}

func Yellow(s string) string {
	return ColorYellow + s + ColorReset
}

func Blue(s string) string {
	return ColorBlue + s + ColorReset
}

func Magenta(s string) string {
	return ColorMagenta + s + ColorReset
}

func Cyan(s string) string {
	return ColorCyan + s + ColorReset
}
