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

var NovelbinInfo = models.WebInfo{
	WebName:   "Novelbin",
	Host:      "https://novelbin.me",
	SearchUrl: "https://novelbin.me/search?keyword=%s",
}

func NovelbinSearch(searchTitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(NovelbinInfo.SearchUrl, searchTitle)

	c := colly.NewCollector()
	var novels []models.NovelData

	c.OnHTML(".list .row", func(e *colly.HTMLElement) {
		Title := e.ChildText(".novel-title")
		Url := e.ChildAttrs("a", "href")
		LatestChapter := e.ChildText(".text-info")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &models.NovelData{
				WebName:          string(NovelbinInfo.WebName),
				Title:            Title,
				Url:              Url[0],
				AvailableChapter: fmt.Sprintf("<= %s", LatestChapter),
			}

			novels = append(novels, *novel)
		}
	})

	err := c.Visit(path)
	if err != nil {
		chErr <- fmt.Errorf("%s %s", NovelbinInfo.WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func NovelbinContent(path string, title string) *models.NovelInfo {
	Target := fmt.Sprintf("%s#tab-chapters-title", path)
	c := colly.NewCollector()
	Author := ""
	Image := ""
	Synopsis := ""
	NovelId := ""

	c.OnHTML(".book img", func(e *colly.HTMLElement) {
		Image = e.Attr("data-src")
	})

	c.OnHTML(".info-meta", func(e *colly.HTMLElement) {
		Author = e.ChildText("li:nth-child(2)")
	})

	c.OnHTML("#rating", func(e *colly.HTMLElement) {
		NovelId = e.Attr("data-novel-id")
	})

	err := c.Visit(Target)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	ListUrl := fmt.Sprintf("%s/ajax/chapter-archive?novelId=%s", NovelbinInfo.Host, NovelId)
	list := NovelbinList(ListUrl)

	res := &models.NovelInfo{
		Title:    title,
		Image:    Image,
		Author:   Author,
		Synopsis: Synopsis,
		Data:     list,
	}

	return res
}

func NovelbinList(url string) []models.ListChapter {
	c := colly.NewCollector()
	var list []models.ListChapter
	Order := 0

	c.OnHTML("ul.list-chapter > li", func(e *colly.HTMLElement) {
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

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	return list
}

func NovelbinGetContent(params *models.ListChapter, wp *models.WorkerPoolContent) {
	defer wp.Wg.Done()
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

	res := wp.Pool.Get().(*models.ListChapter)
	res.Title = params.Title
	res.Order = params.Order
	res.Content = content

	wp.Res <- res
}

func init() {
	var WebName = string(NovelbinInfo.WebName)
	models.MapSearch[WebName] = NovelbinSearch
	models.MapToc[WebName] = NovelbinContent
	models.MapContent[WebName] = NovelbinGetContent
}
