package frontend

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"plugin"
	"strings"

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
	YearName      string
	FolderPath    string
	InputFilePath string
	SoFilePath    string
	Plugin        *plugin.Plugin
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
						Title:      d.Name(),
						Intro:      intro,
						View:       nil,
						SupportWeb: false,
						UID:        cur_year.Title + "." + d.Name(),
					},
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

func (y aocYear) String() string {
	var dstr string
	for _, d := range y.Days {
		dstr += fmt.Sprintf("\n\t\t%v", d)
	}

	return fmt.Sprintf("Year:\n\tTitle: %v\n\tUID: %v\n\tFolderPath: %v\n\tDays: %v\n", y.Title, y.UID, y.FolderPath, dstr)
}

func (d aocDay) String() string {
	return fmt.Sprintf("Day: Title: %v, UID: %v, Folder: %v, InputFile: %v, SoFile: %v", d.Title, d.UID, d.FolderPath, d.InputFilePath, d.SoFilePath)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func getDayData(dir string) (gofile string, sofile string, inputfile string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			switch {
			case info.Name() == "main.go":
				gofile = path
			case strings.HasSuffix(info.Name(), ".go") && gofile == "":
				gofile = path
			case info.Name() == "plugin.so":
				sofile = path
			case strings.HasSuffix(info.Name(), ".so") && sofile == "":
				sofile = path
			case info.Name() == "input.txt":
				inputfile = path
			}
		}
		return nil
	})
	if err != nil {
		fyne.LogError("Error - Get files in day folder "+dir, err)
	}

	return
}

func loadDayPlugin(name string, sofile string, inputfile string) (intro string, p *plugin.Plugin) {
	p, err := plugin.Open(sofile)
	if err != nil {
		panic(err)
	}
	intro_ptr, err := p.Lookup("Intro")
	if err == nil {
		intro = *intro_ptr.(*string)
	} else {
		fyne.LogError("Error - Unable to load intro from "+name, err)
	}

	return
}

func (y aocYear) GenerateView() aocYear {
	y.View = func(w fyne.Window) fyne.CanvasObject {
		var logo *canvas.Image
		mydir, err := os.Getwd()
		if err != nil {
			fyne.LogError("Error - Unable to locate working dir", err)
			logo = &canvas.Image{}
		} else if fileExists(mydir + "/assets/aoc_" + y.Title + ".png") {
			logo = canvas.NewImageFromFile(mydir + "/assets/aoc_" + y.Title + ".png")
		} else if fileExists(mydir + "/assets/aoc.png") {
			logo = canvas.NewImageFromFile(mydir + "/assets/aoc.png")
		} else {
			fyne.LogError("Error - Unable to load year image "+mydir+"/assets/aoc_"+y.Title+".png)", err)
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

func (d aocDay) GenerateView() aocDay {
	d.View = func(w fyne.Window) fyne.CanvasObject {
		//head := widget.NewLabelWithStyle("Advent of Code - "+d.YearName+" Day "+d.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

		tabs := container.NewAppTabs(
			container.NewTabItem("Input", makeDayInputView(d)),
			container.NewTabItem("Asserts", widget.NewLabel("Content of tab 2")),
			container.NewTabItem("Run", widget.NewLabel("Content of tab 3")),
		)

		tail := widget.NewLabel("")

		return container.NewBorder(nil, tail, nil, nil, tabs)
	}
	return d
}

func makeDayInputView(d aocDay) fyne.CanvasObject {
	head := widget.NewLabel("Input file:")
	return container.NewBorder(head, nil, nil, nil, nil)
}
