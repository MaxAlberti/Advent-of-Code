package frontend

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/MaxAlberti/Advent-of-Code/internal/shared"
)

type aocFrontendItem struct {
	Title      string
	Intro      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
	UID        string
}

type aocYear struct {
	aocFrontendItem
	Days       []aocDay
	FolderPath string
}

func (y aocYear) String() string {
	var dstr string
	for _, d := range y.Days {
		dstr += fmt.Sprintf("\n\t\t%v", d)
	}

	return fmt.Sprintf("Year:\n\tTitle: %v\n\tUID: %v\n\tFolderPath: %v\n\tDays: %v\n", y.Title, y.UID, y.FolderPath, dstr)
}

func makeNavigation() {
	FV.FIIndex = map[string][]string{
		"": {"Welcome"},
	}
	FV.FI = map[string]aocFrontendItem{
		"Welcome": {Title: "Welcome", Intro: "", View: welcomeScreen, SupportWeb: false},
	}

	years := getAoCYears()
	for _, y := range years {
		// Add year to Glabal list
		FV.FIIndex[""] = append(FV.FIIndex[""], y.UID)
		// Add year sub list
		FV.FIIndex[y.UID] = []string{}
		// Add to Nav
		FV.FI[y.UID] = y.aocFrontendItem

		// Add days to sub menus
		for _, d := range y.Days {
			// Add day to year menu
			FV.FIIndex[y.UID] = append(FV.FIIndex[y.UID], d.UID)
			// Add to Nav
			FV.FI[d.UID] = d.aocFrontendItem
		}
	}
}

func getAoCYears() []aocYear {
	res := []aocYear{}
	if FV.WorkingDir == "" {
		return res
	}
	root := FV.WorkingDir + "/pkg/aocyear"
	var cur_year aocYear
	cur_year.Title = ""
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root {
			//fmt.Printf("Path: %v - Name: %v\n", path, d.Name())

			switch {
			case root+"/"+d.Name() == path:
				// year folder
				if cur_year.Title != "" {
					cur_year = cur_year.GenerateView()
					res = append(res, cur_year)
				}
				cur_year = aocYear{
					FolderPath: path,
					aocFrontendItem: aocFrontendItem{
						Title:      d.Name(),
						Intro:      "Problems of the year " + d.Name(),
						View:       nil,
						SupportWeb: false,
						UID:        d.Name(),
					},
				}

			case root+"/"+cur_year.Title+"/"+d.Name() == path:
				// day folder
				if cur_year.Title == "" {
					fmt.Println("Warning - Can't parse day before year, Skipped")
					return nil
				}
				_, sofile, inputfile := getDayData(path)
				intro, p := loadDayPlugin(d.Name(), sofile, inputfile)
				day := aocDay{
					YearName:      cur_year.Title,
					FolderPath:    path,
					InputFilePath: inputfile,
					SoFilePath:    sofile,
					Plugin:        p,
					aocFrontendItem: aocFrontendItem{
						Title:      cur_year.Title + " Day " + d.Name(),
						Intro:      intro,
						View:       nil,
						SupportWeb: false,
						UID:        cur_year.Title + "." + d.Name(),
					},
				}
				if shared.FileExists(day.InputFilePath) {
					content, err := os.ReadFile(day.InputFilePath)
					if err != nil {
						fyne.LogError("Error loading file "+day.InputFilePath, err)
					} else {
						day.Input = string(content)
					}
				}
				day = day.GenerateView()

				cur_year.Days = append(cur_year.Days, day)
			default:
				fmt.Printf("Error Parsing - Parts: %v\n", root+"/"+cur_year.Title+"/"+d.Name())
			}
		}
		return nil
	})
	if cur_year.Title != "" {
		cur_year = cur_year.GenerateView()
		res = append(res, cur_year)
	}

	if err != nil {
		fyne.LogError("Error scanning for days", err)
		return res
	}

	return res
}

func (y aocYear) GenerateView() aocYear {
	y.View = func(w fyne.Window) fyne.CanvasObject {
		var logo *canvas.Image
		if FV.WorkingDir == "" {
			logo = &canvas.Image{}
		} else if shared.FileExists(FV.WorkingDir + "/assets/aoc_" + y.Title + ".png") {
			logo = canvas.NewImageFromFile(FV.WorkingDir + "/assets/aoc_" + y.Title + ".png")
		} else if shared.FileExists(FV.WorkingDir + "/assets/aoc.png") {
			logo = canvas.NewImageFromFile(FV.WorkingDir + "/assets/aoc.png")
		} else {
			fyne.LogError("Error - Unable to load year image "+FV.WorkingDir+"/assets/aoc*.png)", errors.New("no image files found"))
			logo = &canvas.Image{}
		}
		logo.FillMode = canvas.ImageFillContain
		logo.SetMinSize(fyne.NewSize(500, 250))
		//logo.ScaleMode = canvas.ImageScaleFastest
		return container.NewCenter(container.NewVBox(
			widget.NewLabelWithStyle("Advent of Code - Year "+y.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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

	return y
}
