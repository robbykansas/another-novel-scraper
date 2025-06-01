package sources

import (
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/models"
	"sort"
	"strconv"
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

type NovelEachPage struct {
	Url  string
	Page int
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
				Url:              fmt.Sprintf("%s%s", NovelfullInfo.Host, Url[0]),
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

func NovelfullContent(path string, title string) *models.NovelInfo {
	c := colly.NewCollector()
	Author := ""
	Image := ""
	Synopsis := ""

	c.OnHTML(".book", func(e *colly.HTMLElement) {
		Image = fmt.Sprintf("%s%s", NovelfullInfo.Host, e.ChildAttr("img", "src"))
	})

	c.OnHTML("desc-text", func(e *colly.HTMLElement) {
		Synopsis = e.Text
	})

	err := c.Visit(path)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	list := NovelfullList(path)

	res := &models.NovelInfo{
		Title:    title,
		Image:    Image,
		Author:   Author,
		Synopsis: Synopsis,
		Data:     list,
	}

	return res
}

func NovelfullList(url string) []models.ListChapter {
	c := colly.NewCollector()
	var total int
	var wg sync.WaitGroup
	channelPage := make(chan []models.ListChapter, 10)
	var list []models.ListChapter

	c.OnHTML("li.last", func(e *colly.HTMLElement) {
		lastUrl := e.ChildAttr("a", "href")
		total, _ = strconv.Atoi(strings.Split(lastUrl, "=")[1])
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	for i := 1; i <= total; i++ {
		wg.Add(1)

		pageUrl := fmt.Sprintf("%s?page=%s", url, strconv.Itoa(i))
		payload := &NovelEachPage{
			Url:  pageUrl,
			Page: i,
		}

		go NovelfullEachPage(payload, &wg, channelPage)
	}

	go func() {
		wg.Wait()
		close(channelPage)
	}()

	// each channel page return list and combine them into a single list
	for chapter := range channelPage {
		list = append(list, chapter...)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})

	return list
}

func NovelfullEachPage(params *NovelEachPage, wg *sync.WaitGroup, list chan<- []models.ListChapter) {
	defer wg.Done()
	var listChapter []models.ListChapter

	Order := (params.Page - 1) * 50

	c := colly.NewCollector()

	c.OnHTML(".row .list-chapter", func(e *colly.HTMLElement) {
		// case where the list is basically just table ul and li where ul has class and li doesn't
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			Title := el.ChildText("a")
			Url := el.ChildAttr("a", "href")
			Order += 1

			info := &models.ListChapter{
				Order: Order,
				Title: Title,
				Url:   fmt.Sprintf("%s%s", NovelfullInfo.Host, Url),
			}

			listChapter = append(listChapter, *info)
		})
	})

	err := c.Visit(params.Url)
	if err != nil {
		log.Fatalf("Error while visiting url with error: %v", err)
	}

	list <- listChapter
}

func NovelfullGetContent(params *models.ListChapter, wp *models.WorkerPoolContent) {
	defer wp.Wg.Done()
	c := colly.NewCollector()
	path := params.Url
	var content string

	c.OnHTML("div#chapter-content", func(e *colly.HTMLElement) {
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
	var WebName = string(NovelfullInfo.WebName)
	models.MapSearch[WebName] = NovelfullSearch
	models.MapToc[WebName] = NovelfullContent
	models.MapContent[WebName] = NovelfullGetContent
}
