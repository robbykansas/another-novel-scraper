package sources

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/flags"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var Host = "https://www.novelhall.com/"

var NovelhallInfo = flags.WebInfo{
	WebName:   "Novelhall",
	SearchUrl: "https://www.novelhall.com/index.php?s=so&module=book&keyword=%s",
}

func NovelhallSearch(searchTitle string, webInfo flags.WebInfo, wg *sync.WaitGroup, ch chan<- []flags.NovelData, chErr chan<- error) {
	defer wg.Done()
	originSearchTitle := searchTitle
	searchTitle = strings.ReplaceAll(searchTitle, " ", "+")
	path := fmt.Sprintf(string(webInfo.SearchUrl), searchTitle)

	c := colly.NewCollector()
	var novels []flags.NovelData
	WebName := "Novelhall"

	c.OnHTML(".section3 table tbody tr", func(e *colly.HTMLElement) {
		Title := e.ChildText("td:nth-child(2)")
		Url := e.ChildAttrs("a", "href")
		LatestChapter := e.ChildText("td:nth-child(3)")

		if strings.Contains(strings.ToLower(Title), strings.ToLower(originSearchTitle)) {
			novel := &flags.NovelData{
				WebName:          WebName,
				Title:            Title,
				Url:              Url[1],
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
}

func NovelhallContent(path string, title string) *NovelInfo {
	var wg sync.WaitGroup
	var channelContent = make(chan ListChapter, 10)
	Target := fmt.Sprintf("%s%s", Host, path)
	c := colly.NewCollector()
	var list []ListChapter
	var getAllContent []ListChapter
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

		info := &ListChapter{
			Order: Order,
			Title: Title,
			Url:   Url,
		}

		list = append(list, *info)
	})

	err := c.Visit(Target)
	if err != nil {
		fmt.Println(err.Error(), "<<<<< error content")
	}

	for _, content := range list {
		wg.Add(1)
		time.Sleep(10 * time.Millisecond)
		go NovelhallGetContent(content, &wg, channelContent)
	}

	go func() {
		wg.Wait()
		close(channelContent)
	}()

	for c := range channelContent {
		getAllContent = append(getAllContent, c)
	}

	sort.Slice(getAllContent, func(i, j int) bool {
		return getAllContent[i].Order < getAllContent[j].Order
	})

	res := &NovelInfo{
		Title:    title,
		Image:    Image,
		Author:   Author,
		Synopsis: Synopsis,
		Data:     getAllContent,
	}

	return res
}

func NovelhallGetContent(params ListChapter, wg *sync.WaitGroup, ch chan<- ListChapter) {
	defer wg.Done()
	c := colly.NewCollector()
	path := fmt.Sprintf("%s%s", Host, params.Url)
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
		fmt.Println(err.Error(), "<<<< error get content novelhall")
	}

	params.Content = content

	ch <- params
}
