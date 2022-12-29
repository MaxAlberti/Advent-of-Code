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

	_, err := get_file_lines("test.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
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
