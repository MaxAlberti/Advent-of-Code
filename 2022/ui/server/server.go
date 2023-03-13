package main

import (
	libs "adventofcode/2022/ui/libs"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func route_days() {
	app.Route("/day01", &libs.Day01)
}

func main() {
	app.Route("/", &libs.LandingPage{})
	route_days()
	app.RunWhenOnBrowser()
	http.Handle("/", &app.Handler{
		Name:        "AdventOfCode2022",
		Description: "Small UI for my Advent of Code 2022",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
