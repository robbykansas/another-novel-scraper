package content

import (
	"robbykansas/another-novel-scraper/cmd/epub"
	"robbykansas/another-novel-scraper/cmd/sources"
	"strings"
)

var AllContent = map[string]func(string, string) *sources.NovelInfo{
	"Novelhall":    sources.NovelhallContent,
	"1stKissNovel": sources.FirstKissNovelContent,
}

func GetContent(content string, folder string, title string) {
	dataContent := strings.Split(content, ",")
	WebName := dataContent[0]
	Url := dataContent[1]

	epubData := AllContent[WebName](Url, title)

	epub.SetEpub(folder, epubData)
}
