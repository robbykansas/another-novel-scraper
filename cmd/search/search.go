package search

import (
	"log"
	"robbykansas/another-novel-scraper/cmd/flags"
	"robbykansas/another-novel-scraper/cmd/sources"
	"robbykansas/another-novel-scraper/cmd/ui/progressbar"
	"time"

	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

var AllWebInfo = []flags.WebInfo{
	sources.NovelhallInfo,
	sources.FirstKissNovelInfo,
}

var AllSources = map[string]func(string, flags.WebInfo, *sync.WaitGroup, chan<- []flags.NovelData, chan<- error){
	"Novelhall":    sources.NovelhallSearch,
	"1stKissNovel": sources.FirstKissNovelSearch,
}

func SearchTitle(title string) (map[string][]flags.NovelData, error) {
	var wg sync.WaitGroup
	var channelRes = make(chan []flags.NovelData, 5)
	var channelErr = make(chan error)
	groupedTitle := make(map[string][]flags.NovelData)

	for _, search := range AllWebInfo {
		wg.Add(1)
		go AllSources[search.WebName.String()](title, search, &wg, channelRes, channelErr)
	}

	progressbarModel := progressbar.InitialModel(len(AllSources))
	p := tea.NewProgram(progressbarModel)

	go func() {
		for {
			// Search for channel error because channel will send error even if its a nil
			err, ok := <-channelErr
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

			if err != nil {
				log.Fatalf("error search content message: %v", err)
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
