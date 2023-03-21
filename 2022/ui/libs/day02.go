package libs

import (
	pck_day_02 "adventofcode/2022/02"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var O_Day02 = CreateDay02()

type Day02 struct {
	Day
}

func CreateDay02() Day02 {
	d := Day02{}
	d.DayStr = "Day02"
	d.Description = "Description for Day 02"
	d.DayUrl = "/day02"

	return d
}

func (c *Day02) Render() app.UI {
	return app.Div().Body(
		c.generateTop(),
		app.Div().Body(
			app.Button().Text("Run").OnClick(c.runBtnClick),
			app.A().Text("\t"+c.Status),
		),
		c.generateBottom(),
	)
}

func (c *Day02) runBtnClick(ctx app.Context, e app.Event) {
	c.Status = "Running..."
	c.Update()
	c_msg := make(chan string)
	go pck_day_02.Run(c.Input, c_msg)
	for msg := range c_msg {
		c.log(msg)
	}
	c.Status = "Done!"
	c.Update()
}
