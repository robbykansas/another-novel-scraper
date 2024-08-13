package search

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/flags"
	"robbykansas/another-novel-scraper/cmd/sources"

	"sync"
)

var AllWebInfo = []flags.WebInfo{
	sources.NovelhallInfo,
	sources.FirstKissNovelInfo,
}

var AllSources = map[string]func(string, flags.WebInfo, *sync.WaitGroup, chan<- []flags.NovelData, chan<- error){
	"Novelhall":    sources.NovelhallSearch,
	"1stKissNovel": sources.FirstKissNovelSearch,
}

func SearchTitle(title string) (map[string][]flags.NovelData, error) {
	var wg sync.WaitGroup
	var channelRes = make(chan []flags.NovelData, 2)
	var channelErr = make(chan error)
	groupedTitle := make(map[string][]flags.NovelData)

	for _, search := range AllWebInfo {
		wg.Add(1)
		go AllSources[search.WebName.String()](title, search, &wg, channelRes, channelErr)
	}

	go func() {
		wg.Wait()
		close(channelErr)
		close(channelRes)
	}()

	for err := range channelErr {
		if err != nil {
			fmt.Println(err, "<<<<<<<<<< error")
		}
	}

	for res := range channelRes {
		if len(res) > 0 {
			for _, g := range res {
				groupedTitle[g.Title] = append(groupedTitle[g.Title], g)
			}
		}
	}

	if len(groupedTitle) > 0 {
		return groupedTitle, nil
	} else {
		return nil, fmt.Errorf("Error")
	}
}
