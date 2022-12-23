// First Go program
package main

import (
	"bufio"
	"errors"
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
	lines, err := get_file_lines("test.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit((1))
	}
	dups, err := sort_bakckpacks(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit((1))
	}
	print_result_assert_one_per_bck(dups)
}

func print_result_assert_one_per_bck(dups [][]string) error {
	var res_str string
	var res []string
	for _, arr := range dups {
		if len(arr) != 1 {
			return errors.New("more than one dup found")
		}
		res_str += arr[0]
		res = append(res, arr[0])
	}
	fmt.Printf("Assuming one per sack:\n\t- As str: %s\n\t- As arr: ", res_str)
	fmt.Println(res)
	return nil
}

func sort_bakckpacks(lines []string) ([][]string, error) {
	var global_dups [][]string
	for _, line := range lines {
		// create rucksack
		rs := create_empty_rucksack()
		// fill rucksack
		if len(line)%2 != 0 {
			return [][]string{}, errors.New("uneven number of items in backpack")
		}
		var comp_size int = len(line) / 2
		for i, char := range line {
			if i < comp_size {
				item := rs.Comp1[string(char)]
				item.Count += 1
				rs.Comp1[string(char)] = item
			} else {
				item := rs.Comp2[string(char)]
				item.Count += 1
				rs.Comp2[string(char)] = item
			}
		}
		// check for duplicates
		var duplicates []string
		for _, item1 := range rs.Comp1 {
			item2 := rs.Comp2[item1.Name]
			if item1.Count > 0 && item2.Count > 0 {
				duplicates = append(duplicates, item1.Name)
			}
		}
		global_dups = append(global_dups, duplicates)
	}
	return global_dups, nil
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
