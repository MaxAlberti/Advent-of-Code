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
