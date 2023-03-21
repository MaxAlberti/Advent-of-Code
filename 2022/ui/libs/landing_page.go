package libs

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type LandingPage struct {
	app.Compo
}

func (c *LandingPage) Render() app.UI {
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
					O_Day01.GetLandingPageTable(),
					O_Day02.GetLandingPageTable(),
				),
			),
		),
	)
}
