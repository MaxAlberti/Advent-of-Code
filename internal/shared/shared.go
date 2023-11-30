package shared

import (
	"os"
	"strings"
)

func GetStringLines(input string) []string {
	return strings.Split(input, "\n")
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
