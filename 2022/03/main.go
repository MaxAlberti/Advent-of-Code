// First Go program
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Item struct {
	Priority int
	Name     string
	Count    int
}

type Rucksack struct {
	Comp1 map[string]Item
	Comp2 map[string]Item
}

// Main function
func main() {
	fmt.Println("Starting")
	_, err := get_file_lines("test.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit((1))
	}
	create_empty_rucksack()
}

func create_empty_rucksack() Rucksack {
	var r Rucksack
	// create empty compartments
	r.Comp1 = make(map[string]Item)
	r.Comp2 = make(map[string]Item)
	i := 97
	item_prio := 1
	for i <= 122 {
		var item Item
		item.Count = 0
		item.Name = string(byte(i))
		item.Priority = item_prio
		r.Comp1[string(byte(i))] = item
		r.Comp2[string(byte(i))] = item
		item_prio += 1
		i += 1
	}
	i = 65
	for i <= 90 {
		var item Item
		item.Count = 0
		item.Name = string(byte(i))
		item.Priority = item_prio
		r.Comp1[string(byte(i))] = item
		r.Comp2[string(byte(i))] = item
		item_prio += 1
		i += 1
	}

	return r
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
