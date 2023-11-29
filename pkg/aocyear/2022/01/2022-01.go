package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Elf struct {
	Number   int
	Calories int
	Indexes  []int
}

// Main function
func Run(input string, c_msg chan string) {
	lines := get_file_lines(input)

	elfes, err := group_lines_to_elfes(lines)
	if err != nil {
		c_msg <- err.Error()
		close(c_msg)
		return
	}
	elfes = sort_elfes_by_cals(elfes)

	if len(elfes) > 0 {
		c_msg <- fmt.Sprintf("Found %s elfes!\nTop:\n\t-Number:\t%s\n\t-Calories:\t%s\n\tIndexes:\t%s\nBottom:\n\t-Number:\t%s\n\t-Calories:\t%s\n\tIndexes:\t%s\n",
			fmt.Sprint(len(elfes)),
			fmt.Sprint(elfes[0].Number),
			fmt.Sprint(elfes[0].Calories),
			fmt.Sprint(elfes[0].Indexes),
			fmt.Sprint(elfes[len(elfes)-1].Number),
			fmt.Sprint(elfes[len(elfes)-1].Calories),
			fmt.Sprint(elfes[len(elfes)-1].Indexes))
	} else {
		c_msg <- "No elfes fonund..."
	}

	top_three, err := get_top_three_cals(elfes)
	if err != nil {
		c_msg <- err.Error()
		close(c_msg)
		return
	}
	c_msg <- fmt.Sprintf("\nThe top3 elfes are carrying %s cals!", fmt.Sprint(top_three))
	close(c_msg)
}

func get_file_lines(input string) []string {
	return strings.Split(input, "\n")
}

func group_lines_to_elfes(lines []string) ([]Elf, error) {
	elfes := []Elf{}
	elf_num := 1
	// create new elf
	var elf Elf
	elf.Calories = 0
	elf.Indexes = []int{}
	elf.Number = elf_num
	// add elf to slice
	elfes = append(elfes, elf)
	elf_num = elf_num + 1
	for index, line := range lines {
		if line != "" {
			// add index
			elfes[len(elfes)-1].Indexes = append(elfes[len(elfes)-1].Indexes, index)
			// add calories
			i, err := strconv.Atoi(line)
			if err != nil {
				return []Elf{}, err
			}
			elfes[len(elfes)-1].Calories += i
		} else {
			// create new elf
			var elf Elf
			elf.Calories = 0
			elf.Indexes = []int{}
			elf.Number = elf_num
			// add elf to slice
			elfes = append(elfes, elf)
			elf_num = elf_num + 1
		}
	}

	return elfes, nil
}

func sort_elfes_by_cals(elfes []Elf) []Elf {
	sort.Slice(elfes, func(i, j int) bool {
		return elfes[i].Calories > elfes[j].Calories
	})
	return elfes
}

func get_top_three_cals(elfes []Elf) (int, error) {
	if len(elfes) < 3 {
		return 0, errors.New("not enough elfes")
	}

	elfes = sort_elfes_by_cals(elfes)

	return (elfes[0].Calories + elfes[1].Calories + elfes[2].Calories), nil
}
