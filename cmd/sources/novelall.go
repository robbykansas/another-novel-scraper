package sources

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/models"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var NovelAllInfo = models.WebInfo{
	WebName:   "NovelAll",
	SearchUrl: "https://www.novelall.com/search/?name=%s",
}

func NovelAllSearch(searchTitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(NovelAllInfo.SearchUrl, searchTitle)

	c := colly.NewCollector()
	var novels []models.NovelData
	WebName := "NovelAll"
	c.OnHTML(".cover-info p.title", func(e *colly.HTMLElement) {
		Title := e.Text
		Url := e.ChildAttr("a", "href")
		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &models.NovelData{
				WebName:          WebName,
				Title:            Title,
				Url:              Url,
				AvailableChapter: "",
			}

			novels = append(novels, *novel)
		}
	})

	err := c.Visit(path)
	if err != nil {
		chErr <- fmt.Errorf("%s %s", WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func init() {
	WebName := string(NovelAllInfo.WebName)
	models.MapSearch[WebName] = NovelAllSearch
}
