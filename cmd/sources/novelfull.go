package sources

import (
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/models"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var NovelfullInfo = models.WebInfo{
	WebName:   "Novelfull",
	Host:      "https://novelfull.net",
	SearchUrl: "https://novelfull.net/search?keyword=%s",
}

func NovelfullSearch(searchTitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(NovelfullInfo.SearchUrl, searchTitle)

	c := colly.NewCollector()
	var novels []models.NovelData

	c.OnHTML(".list .row", func(e *colly.HTMLElement) {
		Title := e.ChildText(".truyen-title")
		Url := e.ChildAttrs("a", "href")
		LatestChapter := e.ChildText(".text-info")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &models.NovelData{
				WebName:          string(NovelfullInfo.WebName),
				Title:            Title,
				Url:              Url[0],
				AvailableChapter: fmt.Sprintf("<= %s", LatestChapter),
			}

			novels = append(novels, *novel)
		}
	})

	err := c.Visit(path)
	if err != nil {
		chErr <- fmt.Errorf("%s %s", NovelfullInfo.WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func NovelfullGetContent(params models.ListChapter, wg *sync.WaitGroup, ch chan<- models.ListChapter) {
	defer wg.Done()
	c := colly.NewCollector()
	path := params.Url
	var content string

	c.OnHTML("div#chr-content", func(e *colly.HTMLElement) {
		e.DOM.Each(func(_ int, s *goquery.Selection) {
			h, _ := s.Html()
			content = fmt.Sprintf("%s \n", h)
		})
	})

	err := c.Visit(path)
	if err != nil {
		log.Fatalf("Error while getting content with error: %v", err)
	}

	params.Content = content

	ch <- params
}

func init() {
	var WebName = string(NovelfullInfo.WebName)
	models.MapSearch[WebName] = NovelfullSearch
	models.MapContent[WebName] = NovelfullGetContent
}
