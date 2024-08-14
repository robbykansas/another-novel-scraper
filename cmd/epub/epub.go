package epub

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/sources"

	"github.com/go-shiori/go-epub"
)

func SetEpub(folder string, content *sources.NovelInfo) {
	epub, err := epub.NewEpub(content.Title)
	if err != nil {
		fmt.Println(err.Error(), "<<<< error epub")
	}

	if content.Author != "0" {
		epub.SetAuthor(content.Author)
	}

	if content.Image != "" {
		imagePath, _ := epub.AddImage(content.Image, "cover.png")
		epub.SetCover(imagePath, "")
	}

	for _, item := range content.Data {
		sectionBody := `<h1>` + item.Title + `</h1>
		<p>` + item.Content + `</p>`
		_, err := epub.AddSection(sectionBody, item.Title, "", "")
		if err != nil {
			fmt.Println(err.Error(), "<<<<<<<<< Error set content epub")
		}
	}

	errEpub := epub.Write(fmt.Sprintf("%s/%s.epub", folder, content.Title))
	if errEpub != nil {
		fmt.Println(errEpub.Error(), "<<<<< error epub write")
	}
}
