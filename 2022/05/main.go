// First Go program
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Cargo struct {
	Stacks [][]string
}

type Operation struct {
	FromIndex int
	ToIndex   int
	NumCrates int
}

// Main function
func main() {
	fmt.Println("Starting")

	lines, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	cargo, ops, err := parse_input(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	cargo_p1, err := perform_ops_p1(cargo, ops, false)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	print_result(cargo_p1, 1)

	// P2
	cargo, ops, err = parse_input(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	cargo_p2, err := perform_ops_p2(cargo, ops, false)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	print_result(cargo_p2, 2)
}

func perform_ops_p1(cargo Cargo, ops []Operation, debug bool) (Cargo, error) {
	if debug {
		print_cargo(cargo)
		fmt.Println("---------------")
	}

	for _, op := range ops {
		// P1
		num := op.NumCrates
		for num > 0 {
			// - Perform pop
			stack, item, err := stack_pop(cargo.Stacks[op.FromIndex])
			if err != nil {
				return Cargo{}, err
			}
			cargo.Stacks[op.FromIndex] = stack
			// - Perform push
			stack = stack_push(cargo.Stacks[op.ToIndex], item)
			cargo.Stacks[op.ToIndex] = stack
			// - Decrement num
			num -= 1
		}
		// - Debug
		if debug {
			print_cargo(cargo)
			fmt.Println("---------------")
		}
		/*
			// P2
			// - Perform pop
			stack, items, err := stack_pop_range(cargo_p2.Stacks[op.FromIndex], num)
			if err != nil {
				return Cargo{}, Cargo{}, err
			}
			cargo_p2.Stacks[op.FromIndex] = stack
			// - Perform push
			stack = stack_push_range(cargo_p2.Stacks[op.ToIndex], items)
			cargo_p2.Stacks[op.ToIndex] = stack
			// - Debug
			if debug_p2 {
				print_cargo(cargo_p2)
				fmt.Println("---------------")
			}
		*/
	}

	return cargo, nil
}

func perform_ops_p2(cargo Cargo, ops []Operation, debug bool) (Cargo, error) {
	if debug {
		print_cargo(cargo)
		fmt.Println("---------------")
	}

	for _, op := range ops {
		// P2
		// - Perform pop
		stack, items, err := stack_pop_range(cargo.Stacks[op.FromIndex], op.NumCrates)
		if err != nil {
			return Cargo{}, err
		}
		cargo.Stacks[op.FromIndex] = make([]string, len(stack))
		copy(cargo.Stacks[op.FromIndex], stack)
		// - Perform push
		stack = stack_push_range(cargo.Stacks[op.ToIndex], items)
		cargo.Stacks[op.ToIndex] = make([]string, len(stack))
		copy(cargo.Stacks[op.ToIndex], stack)
		// - Debug
		if debug {
			print_cargo(cargo)
			fmt.Println("---------------")
		}

		// I fucking hate pointers or whatever
	}

	return cargo, nil
}

func parse_input(lines []string) (Cargo, []Operation, error) {
	var cargo Cargo
	var ops []Operation
	parse_seperation := -1
	for index, line := range lines {
		if line == "" {
			parse_seperation = index
			break
		}
	}
	cargo, err := parse_stack(lines[0:parse_seperation], cargo)
	if err != nil {
		return Cargo{}, []Operation{}, nil
	}
	ops, err = parse_moves(lines[parse_seperation+1:], ops)
	if err != nil {
		return Cargo{}, []Operation{}, nil
	}

	return cargo, ops, nil
}

func parse_stack(lines []string, cargo Cargo) (Cargo, error) {
	for index, line := range lines {
		if (len(line)+1)%4 != 0 {
			return Cargo{}, errors.New("can't parse stack, not multiple of three cols")
		}
		if index == len(lines)-1 {
			break
		}
		if len(cargo.Stacks) == 0 {
			// Create Stacks
			i := 0
			for i < (len(line)+1)/4 {
				cargo.Stacks = append(cargo.Stacks, []string{})
				i += 1
			}
		}

		// Fill Stacks
		max_index := 3
		for max_index < len(line)+1 {
			item := line[max_index-2 : max_index-1]
			if item != " " {
				stack_id := (max_index - 3) / 4
				stack := cargo.Stacks[stack_id]
				stack = append(stack, item) // Backwards because parse
				cargo.Stacks[stack_id] = stack
			}
			max_index += 4
		}
	}

	return cargo, nil
}

func parse_moves(lines []string, ops []Operation) ([]Operation, error) {
	re := regexp.MustCompile(`(?m)^move (?P<num>[0-9]+) from (?P<from>[0-9]+) to (?P<to>[0-9]+)`)
	for _, line := range lines {
		result := make(map[string]string)
		match := re.FindStringSubmatch(line)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		var op Operation
		op.FromIndex, _ = strconv.Atoi(result["from"])
		op.ToIndex, _ = strconv.Atoi(result["to"])
		op.NumCrates, _ = strconv.Atoi(result["num"])
		// Make indexes 0-based
		op.FromIndex -= 1
		op.ToIndex -= 1
		ops = append(ops, op)
	}

	return ops, nil
}

func stack_push(stack []string, item string) []string {
	res := []string{item}
	res = append(res, stack...)

	return res
}

func stack_pop(stack []string) ([]string, string, error) {
	if len(stack) == 0 {
		return []string{}, "", errors.New("can't pop, stack is empty")
	}
	item := stack[0]
	return stack[1:], item, nil
}

func stack_peek(stack []string) (string, error) {
	if len(stack) == 0 {
		return "", errors.New("cant peek on empty stack")
	}
	return stack[0], nil
}

func stack_push_range(stack []string, r []string) []string {
	return append(r, stack...)
}

func stack_pop_range(stack []string, num_items int) ([]string, []string, error) {
	if len(stack) < num_items {
		return []string{}, []string{}, errors.New("can't pop, stack is too small")
	}
	res_i := stack[0:num_items]
	res_s := stack[num_items:]
	return res_s, res_i, nil
}

/*
	func stack_reverse(stack []string) []string {
		res := []string{}
		i := len(stack) - 1
		for i >= 0 {
			res = append(res, stack[i])
			i -= 1
		}

		return res
	}
*/

/*
func copy_cargo(cargo Cargo) Cargo {
	var res Cargo
	res.Stacks = append(res.Stacks, cargo.Stacks...)
	return res
}
*/

func print_cargo(cargo Cargo) {
	for i, stack := range cargo.Stacks {
		fmt.Printf("%d: [ ", i+1)
		i := len(stack) - 1
		for i >= 0 {
			fmt.Printf(" %s ", stack[i])
			i -= 1
		}
		fmt.Println(" ]")
	}
}

func print_result(cargo Cargo, part int) {
	res := ""
	for _, stack := range cargo.Stacks {
		peek, err := stack_peek(stack)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		res += peek
	}
	fmt.Printf("P%d: %s\n", part, res)
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