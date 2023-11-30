package main

import (
	"fmt"
	"reflect"
)

// go build -buildmode=plugin -o ./pkg/plugin/so/plugin.so ./pkg/plugin/so/plugin.go

var V int

func F() {
	fmt.Printf("Hello, number %d\n", V)
}

func G(s string) {
	fmt.Println(s)
}

func H(a ...any) {
	for _, s := range a {
		reflect.TypeOf(s)
		switch reflect.TypeOf(s) {
		case reflect.TypeOf(""):
			fmt.Printf("Type: string, Val: %v\n", s)
		case reflect.TypeOf(1):
			fmt.Printf("Type: int, Val: %v\n", s)
		default:
			fmt.Printf("Type: unhandeled, Val: %v\n", s)
		}
	}
}

func I(ch chan string) {
	defer close(ch)

	ch <- "Channel works!"
}

func Run(inp chan any) {
	out, input, ass := getData(inp)
	defer close(out)
	out <- "Data recieved!"
	out <- fmt.Sprintf("Input: %s\n", input)
	for _, a := range ass {
		out <- fmt.Sprintln(a)
	}
}

func getData(inp chan any) (chan string, string, []Assertion) {
	//defer close(inp)
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

func (ass Assertion) String() string {
	return fmt.Sprintf("Assert: %v == %v", ass.Input, ass.Output)
}

type Assertion struct {
	Input  string
	Output string
}
