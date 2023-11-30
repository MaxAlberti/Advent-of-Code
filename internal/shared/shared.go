package shared

import "strings"

func GetStringLines(input string) []string {
	return strings.Split(input, "\n")
}
