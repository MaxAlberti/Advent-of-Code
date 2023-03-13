package libs

import (
	pck_day_01 "adventofcode/2022/01"
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var O_Day01 = CreateDay01()

type Day01 struct {
	app.Compo

	DayStr      string
	Description string
	Status      string
	Input       string
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
	c.Input = ""
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
			app.Div().Body(
				app.Form().Method("post").EncType("multipart/form-data").Body(
					app.Div().Body(
						app.Label().For("my_input").Text("Choose file to upload: "),
						app.Input().Type("file").ID("my_input").Name("my_input").Accept(".txt").OnChange(c.onFileUpload),
					),
				),
			),
		),
		app.Div().Body(
			app.Button().Text("Run").OnClick(c.runBtnClick),
			app.A().Text("\t"+c.Status),
		),
		app.Div().Body(
			app.A().Text("Output:"),
			app.Div().Body(
				app.Textarea().Style("width", "1000px").Style("height", "500px").Cols(60).Rows(10).Text(c.OutputLog),
			),
		),
	)
}

func (c *Day01) onFileUpload(ctx app.Context, e app.Event) {
	c.OutputLog += "\nFile uploaded"
	files := ctx.JSSrc().Get("files")
	if !files.Truthy() || files.Get("length").Int() == 0 {
		fmt.Println("file not found")
		return
	}

	file := files.Index(0)
	var close func()

	onFileLoad := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		event := args[0]
		content := event.Get("target").Get("result")

		c.Input = content.String()

		close()
		return nil
	})

	onFileLoadError := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		// Your error handling...
		c.OutputLog += "\nError loading input file"

		close()
		return nil
	})

	// To release resources when callback are called.
	close = func() {
		onFileLoad.Release()
		onFileLoadError.Release()
	}

	reader := app.Window().Get("FileReader").New()
	reader.Set("onload", onFileLoad)
	reader.Set("onerror", onFileLoadError)
	reader.Call("readAsText", file, "UTF-8")
}

func (c *Day01) runBtnClick(ctx app.Context, e app.Event) {
	c.Status = "Running..."
	c.Update()
	c_msg := make(chan string)
	go pck_day_01.Run(c.Input, c_msg)
	for msg := range c_msg {
		c.OutputLog += msg
	}
	c.Status = "Done!"
	c.Update()
}
