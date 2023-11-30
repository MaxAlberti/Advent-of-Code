package main

import (
	"fmt"
	"plugin"
)

// go run ./pkg/plugin/go/main.go
func main() {
	p, err := plugin.Open("pkg/plugin/so/plugin.so")
	if err != nil {
		panic(err)
	}
	v, err := p.Lookup("V")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("F")
	if err != nil {
		panic(err)
	}
	*v.(*int) = 7
	f.(func())() // prints "Hello, number 7"

	g, err := p.Lookup("G")
	if err != nil {
		panic(err)
	}
	g.(func(s string))("Hi!")

	h, err := p.Lookup("H")
	if err != nil {
		panic(err)
	}
	h.(func(a ...any))("Hello!", 123, []int{1, 2, 3})

	i, err := p.Lookup("I")
	if err != nil {
		panic(err)
	}
	ch := make(chan string)
	go i.(func(ch chan string))(ch)
	for msg := range ch {
		fmt.Println(msg)
	}

	run, err := p.Lookup("Run")
	if err != nil {
		panic(err)
	}
	com := make(chan any)
	out := make(chan string)
	go run.(func(ch chan any))(com)
	for msg := range com {
		switch msg {
		case "GetOut":
			com <- out
		case "GetInp":
			com <- "My Input"
		case "GetAss":
			com <- [2]string{"Inp1", "Out1"}
			com <- [2]string{"Inp2", "Out2"}
			close(com)
		default:
			fmt.Println("Error - Unhandled command in com channel, closing")
			close(com)
		}
	}
	for msg := range out {
		fmt.Println(msg)
	}
}
