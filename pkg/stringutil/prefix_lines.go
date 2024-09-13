package stringutil

import "strings"

func PrefixLines(text string, prefix string) string {
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
