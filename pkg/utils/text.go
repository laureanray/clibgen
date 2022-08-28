package text

import "github.com/fatih/color"

func Highlight(s string) string {
	magenta := color.New(color.FgHiWhite).Add(color.BgBlack).SprintFunc()
	return magenta(s)
}

func Error(s string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(s)
}

func Info(s string) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	return yellow(s)
}

func Success(s string) string {
	green := color.New(color.FgHiGreen).SprintFunc()
	return green(s)
}
