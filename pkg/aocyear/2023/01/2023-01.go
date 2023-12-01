package main

import (
	"fmt"
	"unicode"

	"github.com/MaxAlberti/Advent-of-Code/internal/aoc"
	"github.com/MaxAlberti/Advent-of-Code/internal/shared"
)

// ----------------------

var Intro string = "TBD - Intro"

var out chan string

func Run(inp chan any) {
	o, i, a := aoc.GetData(inp)
	out = o
	defer close(out)
	run(i, a)
}

func print(m string) {
	out <- m
}

func println(m string) {
	out <- m + "\n"
}

// ^^^^ TEMPLATE ^^^^^^^^

// Main function
func run(input string, asserts []aoc.Assertion) {
	lines := shared.GetStringLines(input)

	println(fmt.Sprint(len(lines)))

	// Part 1
	// Assert 1
	inp_ass1 := []string{"1abc2", "pqr3stu8vwx", "a1b2c3d4e5f", "treb7uchet"}
	exp_ass1 := 142
	res_ass1 := getPt1(inp_ass1)
	equ_ass1 := (exp_ass1 == res_ass1)
	println(fmt.Sprintf("Assert1: All digits combined are %v %v %v", res_ass1, equ_ass1, exp_ass1))
	res_p1 := getPt1(lines)
	println(fmt.Sprintf("Part1: All digits combined are %v", res_p1))

	println(fmt.Sprint(len(lines)))

	// Part 2
	// Assert 2
	inp_ass2 := []string{"two1nine", "eightwothree", "abcone2threexyz", "xtwone3four", "4nineeightseven2", "zoneight234", "7pqrstsixteen"}
	exp_ass2 := 281
	res_ass2 := getPt2(inp_ass2)
	equ_ass2 := (exp_ass2 == res_ass2)
	println(fmt.Sprintf("Assert2: All digits combined are %v %v %v", res_ass2, equ_ass2, exp_ass2))
	res_p2 := getPt2(lines)
	println(fmt.Sprintf("Part2: All digits combined are %v", res_p2))
}

func getPt1(lines []string) int {
	var sum int
	for _, line := range lines {
		digits := getDigits(line)
		first, last := getFirstAndLast(digits)
		combined := 10*first + last
		sum += combined
	}
	return sum
}

func getDigits(input string) []int {
	var digits []int
	for _, char := range input {
		if unicode.IsDigit(char) {
			digits = append(digits, int(char-'0'))
		}
	}
	return digits
}

func getFirstAndLast(input []int) (first int, last int) {
	if len(input) > 0 {
		first = input[0]
		last = input[len(input)-1]
	}
	return
}

func getPt2(lines []string) int {
	var sum int
	for _, line := range lines {
		if line == "" {
			continue
		}
		digits := getNumAndStringDigits(line)
		first, last := getFirstAndLast(digits)
		combined := 10*first + last
		sum += combined
	}
	return sum
}

func getNumAndStringDigits(input string) []int {
	var digits []int
	for i := 0; i < len(input); i++ {
		//println(fmt.Sprint(i))
		char := rune(input[i])
		if unicode.IsDigit(char) {
			digits = append(digits, int(char-'0'))
		} else if len(input)-2-i > 0 {
			// one, two, six,
			// four, five, nine
			// three, seven, eight
			n3 := input[i : i+3]
			switch n3 {
			case "one":
				digits = append(digits, 1)
			case "two":
				digits = append(digits, 2)
			case "six":
				digits = append(digits, 6)
			default:
				if len(input)-3-i > 0 {
					n4 := input[i : i+4]
					switch n4 {
					case "four":
						digits = append(digits, 4)
					case "five":
						digits = append(digits, 5)
					case "nine":
						digits = append(digits, 9)
					default:
						if len(input)-4-i > 0 {
							n5 := input[i : i+5]
							switch n5 {
							case "three":
								digits = append(digits, 3)
							case "seven":
								digits = append(digits, 7)
							case "eight":
								digits = append(digits, 8)
							}
						}
					}
				}
			}
		}
	}
	//println(fmt.Sprintf("Inp: %v, dig: %v", input, digits))
	return digits
}
