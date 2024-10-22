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

var NovelhallInfo = models.WebInfo{
	WebName:   "Novelhall",
	Host:      "https://www.novelhall.com",
	SearchUrl: "https://www.novelhall.com/index.php?s=so&module=book&keyword=%s",
}

func NovelhallSearch(searchTitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(NovelhallInfo.SearchUrl, searchTitle)

	c := colly.NewCollector()
	var novels []models.NovelData

	c.OnHTML(".section3 table tbody tr", func(e *colly.HTMLElement) {
		Title := e.ChildText("td:nth-child(2)")
		Url := e.ChildAttrs("a", "href")
		LatestChapter := e.ChildText("td:nth-child(3)")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &models.NovelData{
				WebName:          string(NovelhallInfo.WebName),
				Title:            Title,
				Url:              Url[1],
				AvailableChapter: fmt.Sprintf("<= %s", LatestChapter),
			}

			novels = append(novels, *novel)
		}
	})

	err := c.Visit(path)
	if err != nil {
		chErr <- fmt.Errorf("%s %s", NovelhallInfo.WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func NovelhallContent(path string, title string) *models.NovelInfo {
	Target := fmt.Sprintf("%s%s", NovelhallInfo.Host, path)
	c := colly.NewCollector()
	var list []models.ListChapter
	Order := 0
	Author := ""
	Image := ""
	Synopsis := ""

	c.OnHTML(".book-img", func(e *colly.HTMLElement) {
		Image = e.ChildAttr("img", "src")
	})

	c.OnHTML("div.book-info div.total.booktag", func(e *colly.HTMLElement) {
		Author = e.ChildText("span:first-child")
	})

	c.OnHTML("div#morelist.book-catalog ul li", func(e *colly.HTMLElement) {
		Title := e.ChildText("a")
		Url := e.ChildAttr("a", "href")
		Order += 1

		info := &models.ListChapter{
			Order: Order,
			Title: Title,
			Url:   Url,
		}

		list = append(list, *info)
	})

	c.OnHTML(".js-close-wrap", func(e *colly.HTMLElement) {
		Synopsis = e.Text
	})

	err := c.Visit(Target)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	res := &models.NovelInfo{
		Title:    title,
		Image:    Image,
		Author:   Author,
		Synopsis: Synopsis,
		Data:     list,
	}

	return res
}

func NovelhallGetContent(params models.ListChapter, wg *sync.WaitGroup, ch chan<- models.ListChapter) {
	defer wg.Done()
	c := colly.NewCollector()
	path := fmt.Sprintf("%s%s", NovelhallInfo.Host, params.Url)
	var content string

	c.OnHTML("div#htmlContent.entry-content", func(e *colly.HTMLElement) {
		// This case for <br> in html, text won't get <br>, it only makes string combined, so we need to use goquery selection
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
	WebName := string(NovelhallInfo.WebName)
	models.MapSearch[WebName] = NovelhallSearch
	models.MapToc[WebName] = NovelhallContent
	models.MapContent[WebName] = NovelhallGetContent
}
