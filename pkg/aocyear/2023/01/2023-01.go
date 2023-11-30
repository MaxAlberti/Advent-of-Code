package main

import "fmt"

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
	out, _, _ := getData(inp)
	defer close(out)
}
