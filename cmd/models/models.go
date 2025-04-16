package models

import "sync"

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

func (c *ListChapter) Reset() {
	c.Order = 0
	c.Title = ""
	c.Url = ""
	c.Content = ""
}

var DefaultPath string = "./ans-config"

var MapSearch = make(map[string]func(string, *sync.WaitGroup, chan<- []NovelData, chan<- error))
var MapToc = make(map[string]func(string, string) *NovelInfo)
var MapContent = make(map[string]func(*ListChapter, *sync.WaitGroup, chan<- *ListChapter, *sync.Pool))

type NovelData struct {
	WebName          string
	Title            string
	Url              string
	AvailableChapter string
}
