package frontend

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

type aocDay struct {
	aocFrontendItem
	FolderPath string
	GoFilePath string
	SoFilePath string
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
	mydir, err := os.Getwd()
	if err != nil {
		fyne.LogError("", err)
		return res
	}
	root := mydir + "/pkg/aocyear"
	var cur_year aocYear
	cur_year.Title = ""
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root {
			//fmt.Printf("Path: %v - Name: %v\n", path, d.Name())

			switch {
			case root+"/"+d.Name() == path:
				// year folder
				if cur_year.Title != "" {
					res = append(res, cur_year)
				}
				cur_year = aocYear{
					FolderPath: path,
					aocFrontendItem: aocFrontendItem{
						Title:      d.Name(),
						Intro:      "Problems of the year " + d.Name(),
						View:       generateYearScreenFunc(d.Name()),
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
				day := aocDay{
					FolderPath: path,
					GoFilePath: "FixMe",
					SoFilePath: "FixMe",
					aocFrontendItem: aocFrontendItem{
						Title:      d.Name(),
						Intro:      "FixMe",
						View:       generateYearScreenFunc(d.Name()), //FixMe
						SupportWeb: false,
						UID:        cur_year.Title + "." + d.Name(),
					},
				}
				cur_year.Days = append(cur_year.Days, day)
			default:
				fmt.Printf("Error Parsing - Parts: %v\n", root+"/"+cur_year.Title+"/"+d.Name())
			}
		}
		return nil
	})
	if cur_year.Title != "" {
		res = append(res, cur_year)
	}

	if err != nil {
		fyne.LogError("Error scanning for days", err)
		return res
	}

	return res
}

func (y aocYear) String() string {
	var dstr string
	for _, d := range y.Days {
		dstr += fmt.Sprintf("\n\t\t%v", d)
	}

	return fmt.Sprintf("Year:\n\tTitle: %v\n\tUID: %v\n\tFolderPath: %v\n\tDays: %v\n", y.Title, y.UID, y.FolderPath, dstr)
}

func (d aocDay) String() string {
	return fmt.Sprintf("Day: Title: %v, UID: %v, Folder: %v, GoFile: %v, SoFile: %v", d.Title, d.UID, d.FolderPath, d.GoFilePath, d.SoFilePath)
}

func generateYearScreenFunc(year string) func(w fyne.Window) fyne.CanvasObject {
	return func(w fyne.Window) fyne.CanvasObject {
		var logo *canvas.Image
		mydir, err := os.Getwd()
		if err != nil {
			fyne.LogError("Error - Unable to locate working dir", err)
			logo = &canvas.Image{}
		} else if fileExists(mydir + "/assets/aoc_" + year + ".png") {
			logo = canvas.NewImageFromFile(mydir + "/assets/aoc_" + year + ".png")
		} else if fileExists(mydir + "/assets/aoc.png") {
			logo = canvas.NewImageFromFile(mydir + "/assets/aoc.png")
		} else {
			fyne.LogError("Error - Unable to load year image "+mydir+"/assets/aoc_"+year+".png)", err)
			logo = &canvas.Image{}
		}
		logo.FillMode = canvas.ImageFillContain
		logo.SetMinSize(fyne.NewSize(500, 250))
		//logo.ScaleMode = canvas.ImageScaleFastest
		return container.NewCenter(container.NewVBox(
			widget.NewLabelWithStyle("Advent of Code - Year "+year, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
