package content

import (
	"log"
	"robbykansas/another-novel-scraper/cmd/epub"
	"robbykansas/another-novel-scraper/cmd/models"
	"robbykansas/another-novel-scraper/cmd/ui/progressbar"
	"robbykansas/another-novel-scraper/cmd/ui/spinner"
	"sort"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func GetContent(content string, folder string, title string) {
	var wg sync.WaitGroup
	var channelContent = make(chan models.ListChapter, 10)
	var AllContent []models.ListChapter
	dataContent := strings.Split(content, ",")
	WebName := dataContent[0]
	Url := dataContent[1]

	spinnerModel := spinner.InitialModel()
	s := tea.NewProgram(spinnerModel)

	go func() {
		if _, err := s.Run(); err != nil {
			cobra.CheckErr(err)
		}
	}()

	listData := models.MapToc[WebName](Url, title)

	s.Send(tea.QuitMsg{})

	progressbarModel := progressbar.InitialModel(len(listData.Data))
	p := tea.NewProgram(progressbarModel)

	for _, content := range listData.Data {
		wg.Add(1)
		time.Sleep(10 * time.Millisecond)
		go models.MapContent[WebName](content, &wg, channelContent)
	}

	go func() {
		for {
			content, ok := <-channelContent
			if ok {
				AllContent = append(AllContent, content)
				p.Send(progressbar.ProgressMsg{})
			} else {
				time.Sleep(500 * time.Millisecond)
				p.Send(tea.Quit())
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		close(channelContent)
	}()

	if _, err := p.Run(); err != nil {
		log.Fatalf("error running progressbar message: %v", err)
	}

	sort.Slice(AllContent, func(i, j int) bool {
		return AllContent[i].Order < AllContent[j].Order
	})

	listData.Data = AllContent

	epub.SetEpub(folder, listData)
}
