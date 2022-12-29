// First Go program
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Main function
func main() {
	fmt.Println("Starting")

	lines, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	process_messages(lines)
}

func process_messages(messages []string) {
	for index, message := range messages {
		fmt.Printf("Line %d:%s\n", index+1, process_message(message))
	}
}

func process_message(message string) string {
	headers := find_message_headers(message)
	top_header := "-"
	if len(headers) > 0 {
		top_header = fmt.Sprint(headers[0])
	}
	return fmt.Sprintf("\n\t- Top Header: %s\n\t- Headers: %s", top_header, fmt.Sprint(headers))
}

func find_message_headers(message string) []int {
	headers := []int{}

	for i := 4; i < len(message); i++ {
		header := message[i-4 : i]
		if is_unique_header(header) {
			headers = append(headers, i)
		}
	}
	return headers
}

func is_unique_header(header string) bool {
	var chars = make(map[rune]bool)

	for _, char := range header {
		_, exists := chars[char]
		if exists {
			return false
		} else {
			chars[char] = true
		}
	}

	return true
}

func get_file_lines(filepath string) ([]string, error) {
	// open file
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	arr := []string{}
	for scanner.Scan() {
		// Read the line
		arr = append(arr, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	} else {
		return arr, err
	}
}
