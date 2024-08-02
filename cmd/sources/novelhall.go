package sources

import (
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/flags"

	"github.com/gocolly/colly/v2"
)

func NovelhallSearch(searchTitle string) ([]flags.NovelData, error) {
	path := fmt.Sprintf(string(flags.Novelhall), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData

	c.OnHTML(".section3 table tbody tr", func(e *colly.HTMLElement) {
		Title := e.ChildText("td:nth-child(2)")
		Url := e.ChildAttrs("a", "href")
		LastChapter := e.ChildText("td:nth-child(3)")

		WebName := "Novelhall"
		// if strings.Contains(strings.ToLower(Title), strings.ToLower(searchTitle)) {
		// 	novel := &flags.NovelData{
		// 		WebName:          WebName,
		// 		Title:            Title,
		// 		Url:              Url[1],
		// 		AvailableChapter: fmt.Sprintf("<= %s", LastChapter),
		// 	}

		// 	novels = append(novels, *novel)
		// }

		novel := &flags.NovelData{
			WebName:          WebName,
			Title:            Title,
			Url:              Url[1],
			AvailableChapter: fmt.Sprintf("<= %s", LastChapter),
		}

		novels = append(novels, *novel)
	})

	err := c.Visit(path)
	if err != nil {
		log.Fatal(err)
	}

	return novels, nil
}
