package main

import (
	"fmt"

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

type Sign struct {
	Value     int
	Name      string
	Beats     string
	BeatsChar string
}

var sign_map = make(map[string]Sign)

// Main function
func run(input string, asserts []aoc.Assertion) {
	print("Starting\n")
	generate_signs()
	lines := shared.GetStringLines(input)

	result1, result2 := calculate_outcome(lines, 0, 3, 6)
	println(fmt.Sprintf("P1: Total player score: %d\nP2: Total player score: %d", result1, result2))
}

func calculate_outcome(lines []string, lose_score int, draw_score int, win_score int) (int, int) {
	p1_total_score := 0
	p2_total_score := 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		elf_char := string(line[0:1])
		player_char := string(line[2:3])
		// Convert to signs
		elf_sign := sign_map[elf_char]
		player_sign := sign_map[player_char]
		// P1 - Calculate round result
		round_res := player_sign.Value
		if elf_sign.Name == player_sign.Name {
			round_res += draw_score
		} else if player_sign.Beats == elf_sign.Name {
			round_res += win_score
		} else {
			round_res += lose_score
		}
		p1_total_score += round_res
		// P2 - Calculate round result
		round_res = 0
		if player_char == "X" {
			// Neet to loose
			player_sign = sign_map[elf_sign.BeatsChar]
			round_res += player_sign.Value
			round_res += lose_score
		} else if player_char == "Y" {
			// Need draw
			player_sign = sign_map[elf_char]
			round_res += player_sign.Value
			round_res += draw_score
		} else {
			// Need to win
			player_sign = sign_map[sign_map[elf_sign.BeatsChar].BeatsChar]
			round_res += player_sign.Value
			round_res += win_score
		}
		p2_total_score += round_res
	}

	return p1_total_score, p2_total_score
}

func generate_signs() {
	var rock Sign
	rock.Value = 1
	rock.Beats = "Scissors"
	rock.Name = "Rock"
	rock.BeatsChar = "C"
	var paper Sign
	paper.Value = 2
	paper.Beats = "Rock"
	paper.Name = "Paper"
	paper.BeatsChar = "A"
	var scissors Sign
	scissors.Value = 3
	scissors.Beats = "Paper"
	scissors.Name = "Scissors"
	scissors.BeatsChar = "B"
	sign_map["A"] = rock
	sign_map["B"] = paper
	sign_map["C"] = scissors
	sign_map["X"] = rock
	sign_map["Y"] = paper
	sign_map["Z"] = scissors
}
