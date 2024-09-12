package sources

import (
	"context"
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/flags"
	"sort"
	"time"

	"sync"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

var FirstKissNovelInfo = flags.WebInfo{
	WebName:   "1stKissNovel",
	SearchUrl: "https://1stkissnovel.org/?s=%s&post_type=wp-manga",
}

func FirstKissNovelSearch(searchTitle string, webInfo flags.WebInfo, wg *sync.WaitGroup, ch chan<- []flags.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(string(webInfo.SearchUrl), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData
	WebName := "1stKissNovel"

	c.OnHTML(".c-tabs-item__content", func(e *colly.HTMLElement) {
		Title := e.ChildText(".post-title")
		Url := e.ChildAttr("a", "href")
		LatestChapter := e.ChildText(".latest-chap")

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
		chErr <- fmt.Errorf("%s %s", WebName, err.Error())
	}

	ch <- novels
	chErr <- nil
}

func FirstKissNovelContent(path string, title string) *NovelInfo {
	var wg sync.WaitGroup
	var channelContent = make(chan ListChapter, 10)
	Target := path
	c := colly.NewCollector()
	var list []ListChapter
	var getAllContent []ListChapter
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

	for _, content := range list {
		wg.Add(1)
		time.Sleep(10 * time.Millisecond)
		go FirstKissNovelGetContent(content, &wg, channelContent)
	}

	go func() {
		wg.Wait()
		close(channelContent)
	}()

	for c := range channelContent {
		getAllContent = append(getAllContent, c)
	}

	res := &NovelInfo{
		Title:    title,
		Image:    Image,
		Author:   Author,
		Synopsis: Synopsis,
		Data:     getAllContent,
	}

	return res
}

func FirstKissNovelList(url string) []ListChapter {
	Target := url
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// var htmlContent string
	var nodes []*cdp.Node
	var listChapter []ListChapter

	err := chromedp.Run(ctx,
		chromedp.Navigate(Target),     // Navigate to the page
		chromedp.Sleep(5*time.Second), // Wait for AJAX content to load
		chromedp.Click(".chapter-readmore", chromedp.ByQuery),
		chromedp.Nodes(".wp-manga-chapter a", &nodes, chromedp.NodeVisible, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Fatal(err)
	}

	order := len(nodes) + 1
	for _, n := range nodes {
		order -= 1
		chapter := &ListChapter{
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

func FirstKissNovelGetContent(params ListChapter, wg *sync.WaitGroup, ch chan<- ListChapter) {
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
