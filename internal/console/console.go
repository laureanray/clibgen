package console

import (
	"github.com/fatih/color"
)

func highlight(s string) string {
	magenta := color.New(color.FgHiWhite).Add(color.BgBlack).SprintFunc()
	return magenta(s)
}

func errorColor(s string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(s)
}

func infoColor(s string) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	return yellow(s)
}

func successColor(s string) string {
	green := color.New(color.FgHiGreen).SprintFunc()
	return green(s)
}

func white(s string) string {
  white := color.New(color.FgWhite).SprintFunc()
  return white(s)
}
