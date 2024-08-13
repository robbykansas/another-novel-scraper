package sources

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/flags"

	"sync"

	"strings"

	"github.com/gocolly/colly/v2"
)

var FirstKissNovelInfo = flags.WebInfo{
	WebName:   "1stKissNovel",
	SearchUrl: "https://1stkissnovel.org/?s=%s&post_type=wp-manga",
}

func FirstKissNovelSearch(searchTitle string, webInfo flags.WebInfo, wg *sync.WaitGroup, ch chan<- []flags.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(string(webInfo.SearchUrl), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData
	WebName := "1stKissNovel"

	c.OnHTML(".c-tabs-item__content", func(e *colly.HTMLElement) {
		Title := e.ChildText(".post-title")
		Url := e.ChildAttr("a", "href")
		LatestChapter := e.ChildText(".latest-chap")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &flags.NovelData{
				WebName:          WebName,
				Title:            Title,
				Url:              Url,
				AvailableChapter: fmt.Sprintf("<= %s", LatestChapter),
			}

			novels = append(novels, *novel)
		}
	})

	err := c.Visit(path)
	if err != nil {
		chErr <- fmt.Errorf("%s %s", WebName, err.Error())
	}
	fmt.Println(novels, "<<<<<<<< novels fist")
	ch <- novels
}
