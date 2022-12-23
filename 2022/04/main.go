// First Go program
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ElfPair struct {
	Elf1Min int
	Elf1Max int
	Elf2Min int
	Elf2Max int
}

// Main function
func main() {
	fmt.Println("Starting")
	pairs, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	redund_pairs, total_partial_reds, err := count_redundant_pairs(pairs)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("P1: Redundant pairs: %d\nP2: PArtially dedundant pairs: %d\n", redund_pairs, total_partial_reds)
}

func count_redundant_pairs(pairs []string) (int, int, error) {
	var total_reds int = 0
	var total_partial_reds int = 0
	for _, pair := range pairs {
		elf_pair, err := parse_elf_pair(pair)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		// Check redundancy
		redundant := is_pair_redundant(elf_pair)
		if redundant {
			total_reds += 1
		}
		// Check partial redundancy
		part_red := is_pair_part_redundant(elf_pair)
		if part_red {
			total_partial_reds += 1
		}
	}
	return total_reds, total_partial_reds, nil
}

func is_pair_redundant(elf_pair ElfPair) bool {
	// Check if elf1 is redundant
	if elf_pair.Elf1Max <= elf_pair.Elf2Max && elf_pair.Elf1Min >= elf_pair.Elf2Min {
		return true
	}
	// Check if elf2 is redundant
	if elf_pair.Elf2Max <= elf_pair.Elf1Max && elf_pair.Elf2Min >= elf_pair.Elf1Min {
		return true
	}

	return false
}

func is_pair_part_redundant(elf_pair ElfPair) bool {
	// {...[...]...}
	//cond1 := elf_pair.Elf1Min <= elf_pair.Elf2Min && elf_pair.Elf1Max >= elf_pair.Elf2Max
	// [...{...}...]
	//cond2 := elf_pair.Elf2Min <= elf_pair.Elf1Min && elf_pair.Elf2Max >= elf_pair.Elf1Max
	// {...[...}...]
	cond1 := elf_pair.Elf1Min <= elf_pair.Elf2Min && elf_pair.Elf1Max <= elf_pair.Elf2Max && elf_pair.Elf1Max >= elf_pair.Elf2Min
	// [...{...]...}
	cond2 := elf_pair.Elf2Min <= elf_pair.Elf1Min && elf_pair.Elf2Max <= elf_pair.Elf1Max && elf_pair.Elf2Max >= elf_pair.Elf1Min
	if cond1 || cond2 || is_pair_redundant(elf_pair) {
		return true
	}
	return false
}

func parse_elf_pair(pair string) (ElfPair, error) {
	pair_split := strings.Split(pair, ",")
	if len(pair_split) != 2 {
		return ElfPair{}, errors.New("could not split pair string ','")
	}
	elf1_split := strings.Split(pair_split[0], "-")
	elf2_split := strings.Split(pair_split[1], "-")
	if len(elf1_split) != 2 {
		return ElfPair{}, errors.New("could not split elf1 string '-'")
	}
	if len(elf2_split) != 2 {
		return ElfPair{}, errors.New("could not split elf1 string '-'")
	}
	var elfPair ElfPair
	i1, err := strconv.Atoi(elf1_split[0])
	if err != nil {
		return ElfPair{}, err
	}
	i2, err := strconv.Atoi(elf1_split[1])
	if err != nil {
		return ElfPair{}, err
	}
	if i1 >= i2 {
		elfPair.Elf1Max = i1
		elfPair.Elf1Min = i2
	} else {
		elfPair.Elf1Max = i2
		elfPair.Elf1Min = i1
	}
	i1, err = strconv.Atoi(elf2_split[0])
	if err != nil {
		return ElfPair{}, err
	}
	i2, err = strconv.Atoi(elf2_split[1])
	if err != nil {
		return ElfPair{}, err
	}
	if i1 >= i2 {
		elfPair.Elf2Max = i1
		elfPair.Elf2Min = i2
	} else {
		elfPair.Elf2Max = i2
		elfPair.Elf2Min = i1
	}

	return elfPair, nil
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
