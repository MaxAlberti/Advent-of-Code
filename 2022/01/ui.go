package pck_day_01

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Day01 struct {
	app.Compo

	DayStr      string
	Description string
	Status      string
	InputPath   string
	OutputLog   string
}

var dayStr = "01"
var desc = "Description for Day 01"

func CreateDay01() Day01 {
	d := Day01{}
	d.DayStr = dayStr
	d.Description = desc

	return d
}

func (c *Day01) OnNav(ctx app.Context) {
	c.InputPath = inpFile
	c.OutputLog = ""
	c.Status = "-NotRunYet-"
}

func (c *Day01) GetLandingPageTable() app.UI {
	return app.Tr().Body(
		app.Td().Text(c.DayStr),
		app.Td().Text(c.Description),
		app.Td().Body(
			app.A().Href("/day01").Text("Link"),
		),
	)
}

func (c *Day01) Render() app.UI {
	return app.Div().Body(
		app.H1().
			Class("title").
			Text("Advent of Code 2022 - Day01"),
		app.Div().Body(
			app.A().Text("Solution for Day 01. Click "),
			app.A().Href("/").Text("here"),
			app.A().Text(" to return."),
		),
		app.Div().Body(
			app.A().Text("Input File:"),
			app.Input().
				Type("text").
				Value(c.InputPath).
				AutoFocus(true).
				OnChange(c.ValueTo(&c.InputPath)),
		),
		app.Div().Body(
			app.Button().Text("Run").OnClick(c.runBtnClick),
			app.A().Text("\t"+c.Status),
		),
		app.Div().Body(
			app.A().Text("Output:"),
			app.Div().Body(
				app.Textarea().Cols(60).Rows(10).Text(c.OutputLog),
			),
		),
	)
}

func (c *Day01) runBtnClick(ctx app.Context, e app.Event) {
	c.Status = "Running..."
	c.Update()
	run(c)
	c.Status = "Done!"
	c.Update()
}
