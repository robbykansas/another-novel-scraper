package sources

import (
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/flags"

	"strings"

	"github.com/gocolly/colly/v2"
)

func FirstKissNovelSearch(searchTitle string, webInfo *flags.NovelInfo) ([]flags.NovelData, error) {
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(string(webInfo.SearchUrl), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData

	c.OnHTML(".c-tabs-item__content", func(e *colly.HTMLElement) {
		Title := e.ChildText(".post-title")
		Url := e.ChildAttr("a", "href")
		LatestChapter := e.ChildText(".latest-chap")

		WebName := "1stKissNovel"
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
		log.Fatal(err)
	}

	return novels, nil
}
