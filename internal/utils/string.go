package utils

import "strings"

func TruncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")] + " ..."
}
