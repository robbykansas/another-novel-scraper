package sources

import (
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/flags"
	"strings"

	"github.com/gocolly/colly/v2"
)

func NovelhallSearch(searchTitle string, webInfo *flags.NovelInfo) ([]flags.NovelData, error) {
	// originSearchTitle := searchTitle
	// searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(string(webInfo.SearchUrl), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData

	c.OnHTML(".section3 table tbody tr", func(e *colly.HTMLElement) {
		Title := e.ChildText("td:nth-child(2)")
		Url := e.ChildAttrs("a", "href")
		LatestChapter := e.ChildText("td:nth-child(3)")

		WebName := "Novelhall"
		if strings.Contains(strings.ToLower(Title), strings.ToLower(searchTitle)) {
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
		log.Fatal(err)
	}

	return novels, nil
}