package sources

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/flags"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var NovelhallInfo = flags.WebInfo{
	WebName:   "Novelhall",
	SearchUrl: "https://www.novelhall.com/index.php?s=so&module=book&keyword=%s",
}

func NovelhallSearch(searchTitle string, webInfo flags.WebInfo, wg *sync.WaitGroup, ch chan<- []flags.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(string(webInfo.SearchUrl), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData
	WebName := "Novelhall"

	c.OnHTML(".section3 table tbody tr", func(e *colly.HTMLElement) {
		Title := e.ChildText("td:nth-child(2)")
		Url := e.ChildAttrs("a", "href")
		LatestChapter := e.ChildText("td:nth-child(3)")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &flags.NovelData{
				WebName:          WebName,
				Title:            Title,
				Url:              Url[1],
				AvailableChapter: fmt.Sprintf("<= %s", LatestChapter),
			}

			novels = append(novels, *novel)
		}
	})

	err := c.Visit(path)
	if err != nil {
		chErr <- fmt.Errorf("%s %s", WebName, err.Error())
	}
	fmt.Println(novels, "<<<<< novels")
	ch <- novels
}
