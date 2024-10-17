package sources

import (
	"context"
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/models"
	"sort"

	"sync"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

var FirstKissNovelInfo = models.WebInfo{
	WebName:   "1stKissNovel",
	SearchUrl: "https://1stkissnovel.org/?s=%s&post_type=wp-manga",
}

func FirstKissNovelSearch(searchTitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(FirstKissNovelInfo.SearchUrl, searchTitle)

	c := colly.NewCollector()
	var novels []models.NovelData
	WebName := "1stKissNovel"

	c.OnHTML(".c-tabs-item__content", func(e *colly.HTMLElement) {
		Title := e.ChildText(".post-title")
		Url := e.ChildAttr("a", "href")
		LatestChapter := e.ChildText(".latest-chap")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
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
		chErr <- fmt.Errorf("%s %s", WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func FirstKissNovelContent(path string, title string) *models.NovelInfo {
	var list []models.ListChapter
	Target := path
	c := colly.NewCollector()
	Author := ""
	Image := ""
	Synopsis := ""

	c.OnHTML(".summary_image", func(e *colly.HTMLElement) {
		Image = e.ChildAttr("img", "src")
	})

	c.OnHTML(".author-content", func(e *colly.HTMLElement) {
		Author = e.ChildText("a")
	})

	err := c.Visit(Target)
	if err != nil {
		fmt.Println(err.Error(), "<<<<< error visit first")
	}

	list = FirstKissNovelList(Target)

	res := &models.NovelInfo{
		Title:    title,
		Image:    Image,
		Author:   Author,
		Synopsis: Synopsis,
		Data:     list,
	}

	return res
}

func FirstKissNovelList(url string) []models.ListChapter {
	Target := url
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// var htmlContent string
	var nodes []*cdp.Node
	var listChapter []models.ListChapter

	err := chromedp.Run(ctx,
		chromedp.Navigate(Target), // Navigate to the page
		// chromedp.Sleep(5*time.Second), // Wait for AJAX content to load
		chromedp.Click(".chapter-readmore", chromedp.ByQuery),
		chromedp.WaitReady(".wp-manga-chapter a"),
		chromedp.Nodes(".wp-manga-chapter a", &nodes, chromedp.NodeVisible, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Fatal(err)
	}

	order := len(nodes) + 1
	for _, n := range nodes {
		order -= 1
		chapter := &models.ListChapter{
			Order: order,
			Title: strings.TrimSpace(n.Children[0].NodeValue),
			Url:   n.AttributeValue("href"),
		}

		listChapter = append(listChapter, *chapter)
	}

	sort.Slice(listChapter, func(i, j int) bool {
		return listChapter[i].Order < listChapter[j].Order
	})

	return listChapter
}

func FirstKissNovelGetContent(params models.ListChapter, wg *sync.WaitGroup, ch chan<- models.ListChapter) {
	defer wg.Done()
	c := colly.NewCollector()
	var content string

	c.OnHTML("div.text-left", func(e *colly.HTMLElement) {
		e.DOM.Each(func(_ int, s *goquery.Selection) {
			h, _ := s.Html()
			content = fmt.Sprintf("%s \n", h)
		})
	})

	err := c.Visit(params.Url)
	if err != nil {
		log.Fatal(err.Error())
	}

	params.Content = content

	ch <- params
}

func init() {
	fmt.Println("<<<<<<<<<< access this")
	models.MapSearch[string(FirstKissNovelInfo.WebName)] = FirstKissNovelSearch
	models.MapToc[string(FirstKissNovelInfo.WebName)] = FirstKissNovelContent
	models.MapContent[string(FirstKissNovelInfo.WebName)] = FirstKissNovelGetContent
}
