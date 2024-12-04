package sources

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/models"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var FoxaholicInfo = models.WebInfo{
	WebName:   "Foxaholic",
	SearchUrl: "https://www.foxaholic.com/?s=%s&post_type=wp-manga",
}

func FoxaholicSearch(searchTitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(FoxaholicInfo.SearchUrl, searchTitle)

	c := colly.NewCollector()
	var novels []models.NovelData
	WebName := "Foxaholic"
	fmt.Println(path, "<<<<<<< access foxaholic")
	c.OnHTML(".c-tabs-item .c-tabs-item__content", func(e *colly.HTMLElement) {
		fmt.Println("<<<<<< go this")
		Title := e.ChildText(".post-title")
		Url := e.ChildAttr(".post-title h3 a", "href")
		LatestChapter := e.ChildText(".latest-chap .chapter a")
		fmt.Println(Title, Url, LatestChapter)
		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			fmt.Println(Title, Url, LatestChapter)
			novel := &models.NovelData{
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
		fmt.Println(err, "<<<<< error")
		chErr <- fmt.Errorf("%s %s", WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func init() {
	WebName := string(FoxaholicInfo.WebName)
	models.MapSearch[WebName] = FoxaholicSearch
}
