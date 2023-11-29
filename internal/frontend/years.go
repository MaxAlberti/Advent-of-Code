package frontend

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

type aocFrontendItem struct {
	Title      string
	Intro      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
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
	getAoCYears()

	FV.YearIndex = map[string][]string{
		"":     {"Welcome", "2022", "2023"},
		"2022": {"2022.01"},
		"2023": {"2023.01"},
		//"2022": {"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7", "Day 8", "Day 9", "Day 10", "Day 11", "Day 12", "Day 13", "Day 14", "Day 15", "Day 16", "Day 17", "Day 18", "Day 19", "Day 20", "Day 21", "Day 22", "Day 23", "Day 24", "Day 25"},
		//"2023": {"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7", "Day 8", "Day 9", "Day 10", "Day 11", "Day 12", "Day 13", "Day 14", "Day 15", "Day 16", "Day 17", "Day 18", "Day 19", "Day 20", "Day 21", "Day 22", "Day 23", "Day 24", "Day 25"},
	}

	FV.Years = map[string]aocFrontendItem{
		"Welcome": {Title: "Welcome", Intro: "", View: welcomeScreen, SupportWeb: false},
		"2022":    {Title: "2022", Intro: "Problems of the year 2022", View: nil, SupportWeb: false},
		"2023":    {Title: "2023", Intro: "Problems of the year 2023", View: nil, SupportWeb: false},
		"2022.01": {Title: "Day 1", Intro: "Problems of the year 2023", View: nil, SupportWeb: false},
		"2023.01": {Title: "Day 1", Intro: "Problems of the year 2023", View: nil, SupportWeb: false},
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
				cur_year = aocYear{FolderPath: path, aocFrontendItem: aocFrontendItem{Title: d.Name()}}

			case root+"/"+cur_year.Title+"/"+d.Name() == path:
				// day folder
				if cur_year.Title == "" {
					fmt.Println("Warning - Can't parse day before year, Skipped")
					return nil
				}
				day := aocDay{FolderPath: path, aocFrontendItem: aocFrontendItem{Title: d.Name()}}
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

	return fmt.Sprintf("Year:\n\tTitle: %v\n\tFolderPath: %v\n\tDays: %v\n", y.Title, y.FolderPath, dstr)
}

func (d aocDay) String() string {
	return fmt.Sprintf("Day: Title: %v, Folder: %v, GoFile: %v, SoFile: %v", d.Title, d.FolderPath, d.GoFilePath, d.SoFilePath)
}
