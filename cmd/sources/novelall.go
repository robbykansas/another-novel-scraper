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

func NovelAllContent(path string, title string) *models.NovelInfo {
	var list []models.ListChapter
	Target := path
	c := colly.NewCollector()
	Author := ""
	Image := ""
	// Synopsis := ""

	c.OnHTML(".detail-info", func(e *colly.HTMLElement) {
		Author = e.ChildText("p:nth-child(2)")
	})

	c.OnHTML(".manga-detailtop img", func(e *colly.HTMLElement) {
		Image = e.Attr("src")
	})

	err := c.Visit(Target)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	list = NovelAllList(Target)

	res := &models.NovelInfo{
		Title:  title,
		Image:  Image,
		Author: Author,
		Data:   list,
	}

	return res
}

func NovelAllList(url string) []models.ListChapter {
	c := colly.NewCollector()
	var list []models.ListChapter
	Order := 0

	c.OnHTML("ul.detail-chlist > li", func(e *colly.HTMLElement) {
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

	maxOrder := len(list)

	if len(list) == 0 {
		log.Fatal("Error: No chapters found")
	}

	for i := 0; i < maxOrder; i++ {
		list[i].Order = maxOrder - i
	}

	return list
}

func NovelAllGetContent(params *models.ListChapter, wp *models.WorkerPoolContent) {
	defer wp.Wg.Done()
	c := colly.NewCollector()
	var content string

	c.OnHTML(".reading-box", func(e *colly.HTMLElement) {
		e.DOM.Each(func(_ int, s *goquery.Selection) {
			h, _ := s.Html()
			content = fmt.Sprintf("%s \n", h)
		})
	})

	err := c.Visit(params.Url)
	if err != nil {
		log.Fatal("Error while getting content")
	}

	res := wp.Pool.Get().(*models.ListChapter)
	res.Title = params.Title
	res.Order = params.Order
	res.Content = content

	wp.Res <- res
}

func init() {
	// WebName := string(NovelAllInfo.WebName)
	// models.MapSearch[WebName] = NovelAllSearch
	// models.MapToc[WebName] = NovelAllContent
	// models.MapContent[WebName] = NovelAllGetContent
}
