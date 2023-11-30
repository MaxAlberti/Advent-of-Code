package frontend

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/MaxAlberti/Advent-of-Code/internal/aoc"
	"github.com/MaxAlberti/Advent-of-Code/internal/shared"
)

type aocDay struct {
	aocFrontendItem
	YearName      string
	FolderPath    string
	InputFilePath string
	Input         string
	SoFilePath    string
	AssFilePath   string
	Plugin        *plugin.Plugin
	Assertions    []Assertion
}

type Assertion struct {
	aoc.Assertion
	IB binding.String
	OB binding.String
}

func (d aocDay) String() string {
	return fmt.Sprintf("Day: Title: %v, UID: %v, Folder: %v, InputFile: %v, SoFile: %v", d.Title, d.UID, d.FolderPath, d.InputFilePath, d.SoFilePath)
}

func (d aocDay) GenerateView() aocDay {
	d.View = func(w fyne.Window) fyne.CanvasObject {
		//head := widget.NewLabelWithStyle("Advent of Code - "+d.YearName+" Day "+d.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

		tabs := container.NewAppTabs(
			container.NewTabItem("Input", makeDayInputView(d)),
			container.NewTabItem("Asserts", makeDayAssertsView(d)),
			container.NewTabItem("Run", makeDayRunView(d)),
		)

		tail := widget.NewLabel("")

		return container.NewBorder(nil, tail, nil, nil, tabs)
	}
	return d
}

func (d aocDay) LoadDayAsserts() aocDay {
	var jsonstr string
	if shared.FileExists(d.AssFilePath) {
		content, err := os.ReadFile(d.AssFilePath)
		if err != nil {
			fyne.LogError("Error loading file "+d.AssFilePath, err)
			return d
		} else {
			jsonstr = string(content)
		}
	}

	d.Assertions = append(d.Assertions, Assertion{Assertion: aoc.Assertion{Input: "T1", Output: "T2"}})

	// Add bindings
	for j := range d.Assertions {
		d.Assertions[j].IB = binding.BindString(&d.Assertions[j].Input)
		d.Assertions[j].OB = binding.BindString(&d.Assertions[j].Output)
	}
	return d
}

func getDayData(dir string) (gofile string, sofile string, inputfile string, assfile string) {
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
			case info.Name() == "asserts.json":
				assfile = path
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

func makeDayInputView(d aocDay) fyne.CanvasObject {
	headbind := binding.NewString()
	headbind.Set(d.InputFilePath)
	head := container.NewHBox(
		widget.NewLabel("Input file: "),
		widget.NewLabelWithData(headbind),
	)

	textbinding := binding.NewString()
	textbinding.Set(d.Input[:200])
	textbox := widget.NewLabelWithData(textbinding)

	content := container.NewScroll(textbox)

	tail := container.NewVBox(
		widget.NewLabelWithStyle("(input text display limited to 200 rows)", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
		widget.NewButton("Select input file (.txt)", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, FV.Window)
					return
				}
				if reader == nil {
					log.Println("Cancelled")
					return
				}
				if reader == nil {
					log.Println("Cancelled")
					return
				}
				defer reader.Close()
				data, err := io.ReadAll(reader)
				if err != nil {
					fyne.LogError("Failed to load text data", err)
					return
				}
				d.Input = string(data)
				d.InputFilePath = reader.URI().Path()

				headbind.Set(d.InputFilePath)
				textbinding.Set(d.Input[:200])

				FV.Window.Content().Refresh()
			}, FV.Window)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
			if FV.WorkingDir != "" {
				uri := storage.NewFileURI(FV.WorkingDir + "/pkg/aocyear")
				lsuri, err := storage.ListerForURI(uri)
				if err != nil {
					fyne.LogError("Error - Unable to create uri to initial folder", err)
					return
				}
				fd.SetLocation(lsuri)
			}

			fd.Show()
		}),
	)

	return container.NewBorder(head, tail, nil, nil, content)
}

func makeDayAssertsView(d aocDay) fyne.CanvasObject {
	head := widget.NewLabel("Add assertions here!")

	content := container.NewDocTabs()
	i := 0
	// make view
	fn_makeEmptyTab := func() *container.TabItem {
		i++
		ass := Assertion{}
		d.Assertions = append(d.Assertions, ass)

		d.Assertions[len(d.Assertions)-1].IB = binding.BindString(&d.Assertions[len(d.Assertions)-1].Input)
		d.Assertions[len(d.Assertions)-1].OB = binding.BindString(&d.Assertions[len(d.Assertions)-1].Output)

		return makeDayAssertsTabItem(i, &d.Assertions[len(d.Assertions)-1])
	}
	content.CreateTab = fn_makeEmptyTab
	// Add existing asserts
	for j := range d.Assertions {
		i++

		content.Append(makeDayAssertsTabItem(i, &d.Assertions[j]))
	}

	tail := widget.NewButton("Save Assertions", func() {
		fmt.Println(d.Assertions)
	})

	return container.NewBorder(head, tail, nil, nil, content)
}

func makeDayAssertsTabItem(i int, ass *Assertion) *container.TabItem {
	ientry := widget.NewMultiLineEntry()
	ientry.Wrapping = fyne.TextWrapOff
	ientry.Scroll = container.ScrollBoth
	ientry.Bind(ass.IB)

	oentry := widget.NewMultiLineEntry()
	oentry.Wrapping = fyne.TextWrapOff
	oentry.Scroll = container.ScrollBoth
	oentry.Bind(ass.OB)

	return container.NewTabItem(
		fmt.Sprintf("Assertion %d", i),
		container.NewHSplit(
			container.NewBorder(
				widget.NewLabel("Input"),
				nil,
				nil,
				nil,
				ientry,
			),
			container.NewBorder(
				widget.NewLabel("Output"),
				nil,
				nil,
				nil,
				oentry,
			),
		),
	)
}

func makeDayRunView(d aocDay) fyne.CanvasObject {
	var output string
	var output_binding = binding.NewString()
	output_binding.Set(output)
	head := widget.NewButton("RUN", func() {
		run, err := d.Plugin.Lookup("Run")
		if err != nil {
			dialog.ShowError(err, FV.Window)
			return
		}
		com := make(chan any)
		out := make(chan string)
		go run.(func(ch chan any))(com)
		for msg := range com {
			switch msg {
			case "GetOut":
				com <- out
			case "GetInp":
				com <- d.Input
			case "GetAss":
				for _, ass := range d.Assertions {
					com <- ass.Assertion
				}
				close(com)
			default:
				fmt.Println("Error - Unhandled command in com channel, closing")
				close(com)
			}
		}
		for msg := range out {
			output += msg
			output_binding.Set(output)
			FV.Window.Content().Refresh()
		}
	})

	content := container.NewScroll(
		widget.NewLabelWithData(output_binding),
	)

	tail := container.NewVBox(
		widget.NewLabel("Assertion Results:"),
		widget.NewLabel("TBD"),
	)

	return container.NewBorder(head, tail, nil, nil, content)
}
