package utils

import "strings"

// RemoveWhitepaces removes common space dividers like spaces, tabs and newlines
func RemoveWhitepaces(str string) string {
	replacer := strings.NewReplacer(
		" ", "",
		"\t", "",
		"\r", "",
		"\n", "")

	return replacer.Replace(str)
}
