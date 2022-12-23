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
	lines, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit((1))
	}
	prio_sum_p1, prio_sum_p2, err := sort_bakckpacks(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit((1))
	}
	fmt.Printf("Summed up P1 priorities: %d\nSummed up P2 priorities: %d\n", prio_sum_p1, prio_sum_p2)
}

func sort_bakckpacks(lines []string) (int, int, error) {
	if len(lines)%3 != 0 {
		return 0, 0, errors.New("elfes not divideable in groups of 3")
	}
	var global_prios_p1 int
	var global_prios_p2 int
	var ruck_set [3]Rucksack
	for index, line := range lines {
		// create rucksack
		rs := create_empty_rucksack()
		// fill rucksack
		if len(line)%2 != 0 {
			return 0, 0, errors.New("uneven number of items in backpack")
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
		for _, item1 := range rs.Comp1 {
			item2 := rs.Comp2[item1.Name]
			if item1.Count > 0 && item2.Count > 0 {
				global_prios_p1 += item1.Priority
			}
		}
		// check group
		ruck_grup_index := index % 3
		ruck_set[ruck_grup_index] = rs
		if ruck_grup_index == 2 {
			// Detect group label
			for _, item := range rs.Comp1 {
				i_name := item.Name
				rs0_stat := ruck_set[0].Comp1[i_name].Count > 0 || ruck_set[0].Comp2[i_name].Count > 0
				rs1_stat := ruck_set[1].Comp1[i_name].Count > 0 || ruck_set[1].Comp2[i_name].Count > 0
				rs2_stat := ruck_set[2].Comp1[i_name].Count > 0 || ruck_set[2].Comp2[i_name].Count > 0
				if rs0_stat && rs1_stat && rs2_stat {
					global_prios_p2 += item.Priority
				}
			}
		}
	}
	return global_prios_p1, global_prios_p2, nil
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
