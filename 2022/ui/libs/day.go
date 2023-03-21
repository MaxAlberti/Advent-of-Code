package libs

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Day struct {
	app.Compo

	DayStr      string
	DayUrl      string
	Description string
	Status      string
	Input       string
	OutputLog   string
}

func (c *Day) OnNav(ctx app.Context) {
	c.Input = ""
	c.OutputLog = ""
	c.Status = "-NotRunYet-"
}

func (c *Day) GetLandingPageTable() app.UI {
	return app.Tr().Body(
		app.Td().Text(c.DayStr),
		app.Td().Text(c.Description),
		app.Td().Body(
			app.A().Href(c.DayUrl).Text("Link"),
		),
	)
}

func (c *Day) Render() app.UI {
	return app.Div().Body(
		c.generateTop(),
		app.Div().Body(
			app.Button().Text("Run").ID("RunButton").OnClick(c.runBtnClick),
			app.A().Text("\t"+c.Status),
		),
		c.generateBottom(),
	)
}

func (c *Day) generateTop() app.UI {
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
	)
}

func (c *Day) generateBottom() app.UI {
	return app.Div().Body(
		app.A().Text("Output:"),
		app.Div().Body(
			app.Textarea().Style("width", "1000px").Style("height", "500px").Cols(60).Rows(10).Text(c.OutputLog),
		),
	)
}

func (c *Day) onFileUpload(ctx app.Context, e app.Event) {
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

func (c *Day) runBtnClick(ctx app.Context, e app.Event) {
	c.Status = "Running..."
	c.log("Your output could be here")
	c.Status = "Done!"
	c.Update()
}

func (c *Day) log(msg string) {
	c.OutputLog += "\n" + msg
}
