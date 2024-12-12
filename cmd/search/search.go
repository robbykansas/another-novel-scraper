package search

import (
	"log"
	"robbykansas/another-novel-scraper/cmd/models"
	"robbykansas/another-novel-scraper/cmd/ui/progressbar"
	"time"

	"sync"

	// Call this to trigger the init from package sources
	_ "robbykansas/another-novel-scraper/cmd/sources"

	tea "github.com/charmbracelet/bubbletea"
)

func SearchTitle(title string) (map[string][]models.NovelData, error) {
	var wg sync.WaitGroup
	var channelRes = make(chan []models.NovelData, 5)
	var channelErr = make(chan error)
	groupedTitle := make(map[string][]models.NovelData)

	for _, search := range models.MapSearch {
		wg.Add(1)
		go search(title, &wg, channelRes, channelErr)
	}

	progressbarModel := progressbar.InitialModel(len(models.MapSearch))
	p := tea.NewProgram(progressbarModel)

	go func() {
		for {
			// Search for channel error because channel will send error even if its a nil
			_, ok := <-channelErr
			if ok {
				// Animated progress based on open received channel
				p.Send(progressbar.ProgressMsg{})
			} else {
				// When channel is closed, it will check animated progress, will quit when its finished,
				// delay is needed because there is on going animated progress
				time.Sleep(1 * time.Second)
				p.Send(tea.Quit())
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		close(channelErr)
		close(channelRes)
	}()

	if _, err := p.Run(); err != nil {
		log.Fatalf("error running progressbar message: %v", err)
	}

	// mapped channel result and grouped it based on title
	for res := range channelRes {
		if len(res) > 0 {
			for _, g := range res {
				groupedTitle[g.Title] = append(groupedTitle[g.Title], g)
			}
		}
	}

	return groupedTitle, nil
}
