package frontend

import (
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentYear = "Welcome"

type frontendVars struct {
	App       fyne.App
	Window    fyne.Window
	Years     map[string]aocFrontendItem
	YearIndex map[string][]string
}

var FV frontendVars

func OpenMainWindow() {
	a := app.NewWithID("max.aoc.app")
	w := a.NewWindow("Advent of Code")

	FV = frontendVars{}
	FV.App = a
	FV.Window = w
	makeNavigation()

	w.SetContent(makeWindow())
	w.SetMaster()
	w.Resize(fyne.NewSize(800, 650))

	w.ShowAndRun()
}

func makeWindow() fyne.CanvasObject {
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setYear := func(y aocFrontendItem) {
		title.SetText(y.Title)
		intro.SetText(y.Intro)

		content.Objects = []fyne.CanvasObject{y.View(FV.Window)}
		content.Refresh()
	}

	main := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)

	return container.NewHSplit(makeNav(setYear, true), main)
}

func makeNav(setYear func(year aocFrontendItem), loadPrevious bool) fyne.CanvasObject {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return FV.YearIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := FV.YearIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := FV.Years[uid]
			if !ok {
				fyne.LogError("Missing year panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if t, ok := FV.Years[uid]; ok {
				FV.App.Preferences().SetString(preferenceCurrentYear, uid)
				setYear(t) //Closure
			}
		},
	}

	if loadPrevious {
		currentPref := FV.App.Preferences().StringWithFallback(preferenceCurrentYear, "Welcome")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	var logo *canvas.Image
	mydir, err := os.Getwd()
	if err != nil {
		fyne.LogError("", err)
		logo = &canvas.Image{}
	} else {
		logo = canvas.NewImageFromFile(mydir + "/assets/aoc.png")
	}
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(500, 250))
	//logo.ScaleMode = canvas.ImageScaleFastest
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to my Advent of Code Repo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewCenter(
			container.NewHBox(
				widget.NewHyperlink("GitHub Repo", parseURL("https://github.com/MaxAlberti/Advent-of-Code")),
				widget.NewLabel("-"),
				widget.NewHyperlink("Advent of Code", parseURL("https://adventofcode.com/")),
			),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
