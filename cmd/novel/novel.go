package novel

import (
	"log"
	"os"
	"robbykansas/another-novel-scraper/cmd/flags"

	tea "github.com/charmbracelet/bubbletea"
)

type Novel struct {
	NovelTitle  string
	ChosenTitle flags.ChosenTitle
	Web         flags.Web
	Folder      string
	Exit        bool
}

func (n *Novel) ExitCLI(p *tea.Program) {
	if n.Exit {
		if err := p.ReleaseTerminal(); err != nil {
			log.Fatal(p.ReleaseTerminal().Error())
		}

		os.Exit(1)
	}
}
