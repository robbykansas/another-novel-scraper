package content

import (
	"context"
	"log"
	"robbykansas/another-novel-scraper/cmd/epub"
	"robbykansas/another-novel-scraper/cmd/models"
	"robbykansas/another-novel-scraper/cmd/ui/progressbar"
	"robbykansas/another-novel-scraper/cmd/ui/spinner"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func GetContent(content string, folder string, title string) {
	list := models.ListContent(content, title)
	wp := models.WorkerPoolContent{
		List:        list.Data,
		Concurrency: 10,
	}
	var SortContent []models.ListChapter

	//  Spinner UI
	spinnerModel := spinner.InitialModel()
	s := tea.NewProgram(spinnerModel)

	go func() {
		if _, err := s.Run(); err != nil {
			cobra.CheckErr(err)
		}
	}()

	s.Send(tea.QuitMsg{})

	// Progressbar UI
	progressbarModel := progressbar.InitialModel(len(wp.List))
	p := tea.NewProgram(progressbarModel)

	go func() {
		for {
			content, ok := <-wp.Res
			if ok {
				copy := *content
				SortContent = append(SortContent, copy)
				content.Reset()
				wp.Pool.Put(content)
				p.Send(progressbar.ProgressMsg{})
			} else {
				time.Sleep(1000 * time.Millisecond)
				p.Send(tea.Quit())
				break
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, err := p.Run(); err != nil {
					log.Fatalf("error running progressbar message: %v", err)
				}
			}
		}
	}(ctx)

	wp.Run(content, title)

	cancel()

	sort.Slice(SortContent, func(i, j int) bool {
		return SortContent[i].Order < SortContent[j].Order
	})

	list.Data = SortContent

	epub.SetEpub(folder, list)
}
