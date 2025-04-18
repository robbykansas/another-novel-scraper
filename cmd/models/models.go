package models

import (
	"strings"
	"sync"
	"time"
)

type ListChapter struct {
	Order   int
	Title   string
	Url     string
	Content string
}

type NovelInfo struct {
	Title    string
	Image    string
	Author   string
	Synopsis string
	Data     []ListChapter
}

type Web string

type WebInfo struct {
	WebName   Web
	Host      string
	SearchUrl string
}

type NovelData struct {
	WebName          string
	Title            string
	Url              string
	AvailableChapter string
}

type WorkerPoolContent struct {
	List        []ListChapter
	Concurrency int
	Wg          sync.WaitGroup
	ch          chan *ListChapter // give task
	Res         chan *ListChapter // get result
	Pool        *sync.Pool
}

func ListContent(content string, title string) *NovelInfo {
	split := strings.Split(content, ",")
	WebName := split[0]
	Url := split[1]

	list := MapToc[WebName](Url, title)
	return list
}

func (wp *WorkerPoolContent) worker(web string) {
	for task := range wp.ch {
		go MapContent[web](task, wp)
	}
}

func (wp *WorkerPoolContent) Run(content string, title string) {
	// initialize the tasks channel and pool
	wp.ch = make(chan *ListChapter, 10)
	wp.Res = make(chan *ListChapter, 10)
	wp.Pool = &sync.Pool{
		New: func() interface{} {
			return &ListChapter{}
		},
	}

	split := strings.Split(content, ",")
	WebName := split[0]

	list := ListContent(content, title)

	// start workers
	for i := 0; i < wp.Concurrency; i++ {
		go wp.worker(WebName)
	}

	// send tasks to the channel
	for i := range list.Data {
		wp.Wg.Add(1)
		wp.ch <- &list.Data[i]

		switch WebName {
		case "Novelbin", "NovelAll": //has anti scraping which makes it can't scrape too fast or got too many request
			time.Sleep(1000 * time.Millisecond)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	wp.Wg.Wait()
	close(wp.ch)
	close(wp.Res)
}

func (c *ListChapter) Reset() {
	c.Order = 0
	c.Title = ""
	c.Url = ""
	c.Content = ""
}

var DefaultPath string = "./ans-config"

var MapSearch = make(map[string]func(string, *sync.WaitGroup, chan<- []NovelData, chan<- error))
var MapToc = make(map[string]func(string, string) *NovelInfo)
var MapContent = make(map[string]func(*ListChapter, *WorkerPoolContent))
