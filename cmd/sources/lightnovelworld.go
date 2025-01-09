package sources

import (
	"context"
	"fmt"
	"log"
	"robbykansas/another-novel-scraper/cmd/models"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var LightNovelWorld = models.WebInfo{
	WebName:   "LightNovelWorld",
	SearchUrl: "https://www.lightnovelworld.co/search",
}

func LightNovelWorldSearch(searchtitle string, wg *sync.WaitGroup, ch chan<- []models.NovelData, chErr chan<- error) {
	defer wg.Done()

	Target := LightNovelWorld.SearchUrl

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var nodes []*cdp.Node
	fmt.Println(Target)
	err := chromedp.Run(ctx,
		chromedp.Navigate(Target),
		// chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible("#inputContent", chromedp.ByID),
		chromedp.SendKeys("#inputContent", searchtitle, chromedp.ByID),
		// chromedp.WaitVisible("#novelListBase ul", chromedp.ByID),
		chromedp.Nodes(".novel-list .novel-item a", &nodes, chromedp.NodeEnabled, chromedp.ByQueryAll),
	)

	if err != nil {
		fmt.Println(err, "<<< error lightnovelworld")
		log.Fatal(err)
	}

	for _, n := range nodes {
		fmt.Println(n, "<<<< range nodes")
	}

	fmt.Println("<<<<<<<<< done")
}

// func init() {
// 	WebName := string(LightNovelWorld.WebName)
// 	models.MapSearch[WebName] = LightNovelWorldSearch
// }
