package epub

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		imagePath, _ := RetrieveImage(content.Image)
		imagePathInternal, err := epub.AddImage(imagePath, "cover.jpg")
		if err != nil {
			fmt.Println(err.Error(), "<<<<<< error image path")
		}

		errCover := epub.SetCover(imagePathInternal, "")
		if errCover != nil {
			fmt.Println(err.Error(), "<<<< error set cover")
		}
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

func RetrieveImage(source string) (string, error) {
	fmt.Println(source, "<<<<<<< source")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", source, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	response, e := client.Do(req)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf("%#v", response.Body)
	fmt.Println(response.StatusCode, "<<<<<< status code")
	defer response.Body.Close()

	//open a file for writing
	location := "./cmd/epub/assets/cover.jpg"
	file, err := os.Create(location)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	d, err := io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(d, "<<<<<<<")

	return location, nil
}
