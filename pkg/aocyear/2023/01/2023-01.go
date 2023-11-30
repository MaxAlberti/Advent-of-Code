package main

import "github.com/MaxAlberti/Advent-of-Code/internal/aoc"

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

}
