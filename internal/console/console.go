package console

import (
	"fmt"

	"github.com/fatih/color"
)

// TODO: Refactor this to just print directly instead
// of doing fmt.Println(console.Hihlight()) and so on...
func Higlight(format string, a ...any) string {
	magenta := color.New(color.FgHiWhite).Add(color.BgBlack).SprintFunc()
	return magenta(fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(fmt.Sprintf(format, a...))
}

func Info(format string, a ...any) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	return yellow(fmt.Sprintf(format, a...))
}

func Success(format string, a ...any) string {
	green := color.New(color.FgHiGreen).SprintFunc()
	return green(fmt.Sprintf(format, a...))
}

func Normal(format string, a ...any) string {
  white := color.New(color.FgWhite).SprintFunc()
  return white(fmt.Sprintf(format, a...))
}
