package search

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/flags"
	"robbykansas/another-novel-scraper/cmd/sources"
)

func SearchTitle(title string) (map[string][]flags.NovelData, error) {
	groupedTitle := make(map[string][]flags.NovelData)
	AllSources := map[string]func(string, *flags.NovelInfo) ([]flags.NovelData, error){
		"Novelhall":    sources.NovelhallSearch,
		"1stKissNovel": sources.FirstKissNovelSearch,
	}

	for _, search := range flags.WebInfo {
		result, err := AllSources[search.WebName.String()](title, &search)
		if err != nil {
			fmt.Println(err, "<<<<<<< error this")
		}

		for _, g := range result {
			groupedTitle[g.WebName] = append(groupedTitle[g.WebName], g)
		}
	}

	// fmt.Printf("%+v\n", groupedTitle)

	return groupedTitle, nil
}
