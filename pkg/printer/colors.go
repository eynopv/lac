package printer

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
	return c + s + ColorReset
}

func Red(s string) string {
	return Colorize(s, ColorRed)
}

func Green(s string) string {
	return Colorize(s, ColorGreen)
}

func Yellow(s string) string {
	return Colorize(s, ColorYellow)
}

func Blue(s string) string {
	return Colorize(s, ColorBlue)
}

func Magenta(s string) string {
	return Colorize(s, ColorMagenta)
}

func Cyan(s string) string {
	return Colorize(s, ColorCyan)
}
