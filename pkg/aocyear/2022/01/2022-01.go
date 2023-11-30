package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var Intro string = "TBD - Intro"

type Elf struct {
	Number   int
	Calories int
	Indexes  []int
}
type Assertion struct {
	Input  string
	Output string
}

func (ass Assertion) String() string {
	return fmt.Sprintf("Assert: %v == %v", ass.Input, ass.Output)
}

func getData(inp chan any) (chan string, string, []Assertion) {
	var out chan string
	var input string
	var assertions []Assertion

	// Get out chan
	inp <- "GetOut"
	var anyValue interface{} = <-inp
	if channelValue, ok := anyValue.(chan string); ok {
		out = channelValue
	} else {
		fmt.Println("Error - Could not resolve output channel, out")
	}

	// Get input
	inp <- "GetInp"
	anyValue = <-inp
	if strValue, ok := anyValue.(string); ok {
		input = strValue
	} else {
		fmt.Println("Error - Could not resolve output channel, inp")
	}

	// Get Asserts
	inp <- "GetAss"
	for resp := range inp {
		anyValue = resp
		if arrValue, ok := anyValue.([2]string); ok {
			assertions = append(assertions, Assertion{Input: arrValue[0], Output: arrValue[1]})
		} else {
			fmt.Println("Error - Could not resolve output channel, ass")
		}
	}

	return out, input, assertions
}

// Main function
func Run(inp chan any) {
	out, input, _ := getData(inp)
	defer close(out)
	lines := get_file_lines(input)

	elfes, err := group_lines_to_elfes(lines)
	if err != nil {
		out <- err.Error()
		return
	}
	elfes = sort_elfes_by_cals(elfes)

	if len(elfes) > 0 {
		out <- fmt.Sprintf("Found %s elfes!\nTop:\n\t-Number:\t%s\n\t-Calories:\t%s\n\tIndexes:\t%s\nBottom:\n\t-Number:\t%s\n\t-Calories:\t%s\n\tIndexes:\t%s\n",
			fmt.Sprint(len(elfes)),
			fmt.Sprint(elfes[0].Number),
			fmt.Sprint(elfes[0].Calories),
			fmt.Sprint(elfes[0].Indexes),
			fmt.Sprint(elfes[len(elfes)-1].Number),
			fmt.Sprint(elfes[len(elfes)-1].Calories),
			fmt.Sprint(elfes[len(elfes)-1].Indexes))
	} else {
		out <- "No elfes fonund..."
	}

	top_three, err := get_top_three_cals(elfes)
	if err != nil {
		out <- err.Error()
		close(out)
		return
	}
	out <- fmt.Sprintf("\nThe top3 elfes are carrying %s cals!", fmt.Sprint(top_three))
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
