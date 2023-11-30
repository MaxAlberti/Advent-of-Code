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
	h.(func(a ...any))("Hello!")

	i, err := p.Lookup("I")
	if err != nil {
		panic(err)
	}
	ch := make(chan string)
	go i.(func(ch chan string))(ch)
	for msg := range ch {
		fmt.Println(msg)
	}
}
