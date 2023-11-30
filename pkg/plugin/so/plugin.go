package main

import "fmt"

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
		println(s)
	}
}

func I(ch chan string) {
	defer close(ch)

	ch <- "Channel works!"
}
