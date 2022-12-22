// First Go program
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Sign struct {
	Value     int
	Name      string
	Beats     string
	BeatsChar string
}

var sign_map = make(map[string]Sign)

// Main function
func main() {
	fmt.Println("Starting")
	generate_signs()
	lines, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	result1, result2 := calculate_outcome(lines, 0, 3, 6)
	fmt.Printf("P1: Total player score: %d\nP2: Total player score: %d", result1, result2)
}

func calculate_outcome(lines []string, lose_score int, draw_score int, win_score int) (int, int) {
	p1_total_score := 0
	p2_total_score := 0
	for _, line := range lines {
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