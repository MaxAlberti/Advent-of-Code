package pck_day_01

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Elf struct {
	Number   int
	Calories int
	Indexes  []int
}

var inpFile = "input.txt"

// Main function
func run(c *Day01) {
	lines, err := get_file_lines(c.InputPath)
	if err != nil {
		c.OutputLog += err.Error()
		return
	}
	elfes, err := group_lines_to_elfes(lines)
	if err != nil {
		c.OutputLog += err.Error()
		return
	}
	elfes = sort_elfes_by_cals(elfes)

	if len(elfes) > 0 {
		c.OutputLog += fmt.Sprintf("\nFound %s elfes!\nTop:\n\t-Number:\t%s\n\t-Calories:\t%s\n\tIndexes:\t%s\nBottom:\n\t-Number:\t%s\n\t-Calories:\t%s\n\tIndexes:\t%s\n",
			fmt.Sprint(len(elfes)),
			fmt.Sprint(elfes[0].Number),
			fmt.Sprint(elfes[0].Calories),
			fmt.Sprint(elfes[0].Indexes),
			fmt.Sprint(elfes[len(elfes)-1].Number),
			fmt.Sprint(elfes[len(elfes)-1].Calories),
			fmt.Sprint(elfes[len(elfes)-1].Indexes))
	} else {
		c.OutputLog += "\nNo elfes fonund..."
	}

	top_three, err := get_top_three_cals(elfes)
	if err != nil {
		c.OutputLog += err.Error()
		return
	}
	c.OutputLog += fmt.Sprintf("\nThe top3 elfes are carrying %s cals!", fmt.Sprint(top_three))
}

func get_file_lines(filepath string) ([]string, error) {
	// open file
	arr := []string{}
	f, err := os.Open(filepath)
	if err != nil {
		return arr, err
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

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
