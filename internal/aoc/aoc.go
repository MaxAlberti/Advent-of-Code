package aoc

import (
	"fmt"
)

type Assertion struct {
	Input  string
	Output string
}

func (ass Assertion) String() string {
	return fmt.Sprintf("Assert: %v == %v", ass.Input, ass.Output)
}

func GetData(inp chan any) (chan string, string, []Assertion) {
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
		if ass, ok := anyValue.(Assertion); ok {
			assertions = append(assertions, ass)
		} else {
			fmt.Println("Error - Could not resolve output channel, ass")
		}
	}

	return out, input, assertions
}
