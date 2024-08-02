package steps

import "robbykansas/another-novel-scraper/cmd/flags"

// A StepSchema contains the data that is used
// for an individual step of the CLI
type StepSchema struct {
	StepName  string // The name of a given step
	ListTitle []Item // The slice of each option for a given step
	Headers   string // The title displayed at the top of a given step
	Field     string
}

// Steps contains a slice of steps
type Steps struct {
	Steps map[string]StepSchema
}

// An Item contains the data for each option
// in a StepSchema.Options
type Item struct {
	Flag, WebName, AvailableChapter string
}

// InitSteps initializes and returns the *Steps to be used in the CLI program
func InitSteps(projectType flags.Web) *Steps {
	steps := &Steps{
		map[string]StepSchema{
			"web": {
				StepName: "web for scraping",
				ListTitle: []Item{
					{
						WebName:          "Syosetu",
						AvailableChapter: "1-100",
					},
					{
						WebName:          "NovelFull",
						AvailableChapter: "1-50",
					},
				},
				Headers: "Where you want the novel scraping from?",
				Field:   projectType.String(),
			},
		},
	}

	return steps
}
