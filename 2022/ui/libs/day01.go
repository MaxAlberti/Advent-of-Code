package libs

import (
	pck_day_01 "adventofcode/2022/01"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var O_Day01 = CreateDay01()

type Day01 struct {
	Day
}

func CreateDay01() Day01 {
	d := Day01{}
	d.DayStr = "01"
	d.Description = "Description for Day 01"
	d.DayUrl = "/day01"

	return d
}

func (c *Day01) Render() app.UI {
	return app.Div().Body(
		c.generateTop(),
		app.Div().Body(
			app.Button().Text("Run").ID("RunButton").OnClick(c.runBtnClick),
			app.A().Text("\t"+c.Status),
		),
		c.generateBottom(),
	)
}

func (c *Day01) runBtnClick(ctx app.Context, e app.Event) {
	c.Status = "Running..."
	c.Update()
	c_msg := make(chan string)
	go pck_day_01.Run(c.Input, c_msg)
	for msg := range c_msg {
		c.log(msg)
	}
	c.Status = "Done!"
	c.Update()
}
