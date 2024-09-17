package content

import (
	"log"
	"robbykansas/another-novel-scraper/cmd/epub"
	"robbykansas/another-novel-scraper/cmd/sources"
	"robbykansas/another-novel-scraper/cmd/ui/progressbar"
	"sort"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var ListContent = map[string]interface{}{
	"Novelhall":    sources.NovelhallContent,
	"1stKissNovel": sources.FirstKissNovelContent,
}

var GetAllContent = map[string]func(sources.ListChapter, *sync.WaitGroup, chan<- sources.ListChapter){
	"Novelhall":    sources.NovelhallGetContent,
	"1stKissNovel": sources.FirstKissNovelGetContent,
}

func GetContent(content string, folder string, title string) {
	var wg sync.WaitGroup
	var channelContent = make(chan sources.ListChapter, 10)
	var AllContent []sources.ListChapter
	dataContent := strings.Split(content, ",")
	WebName := dataContent[0]
	Url := dataContent[1]

	listData := ListContent[WebName].(func(string, string) *sources.NovelInfo)(Url, title)

	progressbarModel := progressbar.InitialModel(len(listData.Data))
	p := tea.NewProgram(progressbarModel)

	for _, content := range listData.Data {
		wg.Add(1)
		time.Sleep(10 * time.Millisecond)
		go GetAllContent[WebName](content, &wg, channelContent)
	}

	go func() {
		for {
			content, ok := <-channelContent
			if ok {
				AllContent = append(AllContent, content)
				p.Send(progressbar.ProgressMsg{})
			} else {
				time.Sleep(500 * time.Millisecond)
				for {
					p.Send(progressbar.ProgressMsg{})
				}
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

	// for c := range channelContent {
	// 	AllContent = append(AllContent, c)
	// }

	sort.Slice(AllContent, func(i, j int) bool {
		return AllContent[i].Order < AllContent[j].Order
	})

	listData.Data = AllContent

	epub.SetEpub(folder, listData)
}
