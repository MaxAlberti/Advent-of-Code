package main

import (
	pck_day_01 "adventofcode/2022/01"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var d01 = pck_day_01.CreateDay01()

type landingPage struct {
	app.Compo
}

func (c *landingPage) Render() app.UI {
	return app.Div().Body(
		app.H1().
			Class("title").
			Text("Advent of Code 2022"),
		app.P().
			Class("text").
			Text("Select one of the days bellow to see more:"),
		app.Div().Body(
			app.Table().Style("width", "100%").Style("border", "1px solid black").Body(
				app.TBody().Body(
					app.Tr().Style("border", "1px solid black").Body(
						app.Td().Style("border", "1px solid black").Text("Day"),
						app.Td().Style("border", "1px solid black").Text("Description"),
						app.Td().Style("border", "1px solid black").Text("Link"),
					),
					d01.GetLandingPageTable(),
				),
			),
		),
	)
}

func route_days() {
	app.Route("/day01", &d01)
}

func main() {
	app.Route("/", &landingPage{})
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
