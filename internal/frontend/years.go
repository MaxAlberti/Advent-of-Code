package frontend

func makeYears() {
	FV.YearIndex = map[string][]string{
		"":     {"Welcome", "2022", "2023"},
		"2022": {"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7", "Day 8", "Day 9", "Day 10", "Day 11", "Day 12", "Day 13", "Day 14", "Day 15", "Day 16", "Day 17", "Day 18", "Day 19", "Day 20", "Day 21", "Day 22", "Day 23", "Day 24", "Day 25"},
		"2023": {"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7", "Day 8", "Day 9", "Day 10", "Day 11", "Day 12", "Day 13", "Day 14", "Day 15", "Day 16", "Day 17", "Day 18", "Day 19", "Day 20", "Day 21", "Day 22", "Day 23", "Day 24", "Day 25"},
	}

	FV.Years = map[string]aocYear{
		"Welcome": {Title: "Welcome", Intro: "", View: welcomeScreen, SupportWeb: true},
		"2022":    {Title: "2022", Intro: "Problems of the year 2022", View: nil, SupportWeb: true},
		"2023":    {Title: "2023", Intro: "Problems of the year 2023", View: nil, SupportWeb: true},
	}
}
