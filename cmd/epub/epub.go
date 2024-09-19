package epub

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"robbykansas/another-novel-scraper/cmd/sources"
	"strconv"

	"github.com/go-shiori/go-epub"
)

func SetEpub(folder string, content *sources.NovelInfo) {
	epub, err := epub.NewEpub(content.Title)
	if err != nil {
		log.Fatalf("Error creating epub: %v", err)
	}

	if content.Author != "0" {
		epub.SetAuthor(content.Author)
	}

	if content.Synopsis != "" {
		epub.SetDescription(content.Synopsis)
	}

	if content.Image != "" {
		imagePath, _ := RetrieveImage(content.Image)
		imagePathInternal, err := epub.AddImage(imagePath, "cover.jpg")
		if err != nil {
			log.Fatal(err)
		}

		errCover := epub.SetCover(imagePathInternal, "")
		if errCover != nil {
			log.Fatal(errCover)
		}
	}

	var vol = 1
	var counter = 0

	for _, item := range content.Data {
		var sectionBody string
		var titleSection string
		if counter == 0 {
			var tempEndVol string
			stringVol := strconv.Itoa(vol)
			if len(content.Data) >= 100 {
				tempEndVol = content.Data[98].Title
			} else {
				tempEndVol = content.Data[len(content.Data)-1].Title
			}
			startVol := content.Data[counter].Title
			titleSection = fmt.Sprintf("Vol %s: %s => %s", stringVol, startVol, tempEndVol)
			sectionBody = fmt.Sprintf("<h1> %s </h1>", titleSection)
		}

		counter += 1

		if counter%100 == 0 {
			vol += 1
			stringVol := strconv.Itoa(vol)
			if vol*100 < len(content.Data) {
				startVol := content.Data[counter-1].Title
				endVol := content.Data[(vol*100)-2].Title
				titleSection = fmt.Sprintf("Vol %s: %s => %s", stringVol, startVol, endVol)
				sectionBody = fmt.Sprintf("<h1> %s </h1>", titleSection)
			} else {
				startVol := content.Data[counter-1].Title
				endVol := content.Data[len(content.Data)-1].Title
				titleSection = fmt.Sprintf("Vol %s: %s => %s", stringVol, startVol, endVol)
				sectionBody = fmt.Sprintf("<h1> %s </h1>", titleSection)
			}
		}

		sectionPath, err := epub.AddSection(sectionBody, titleSection, "", "")
		if err != nil {
			log.Fatal(err)
		}

		subSectionBody := fmt.Sprintf(`<h1> %s </h1>
		<p> %s </p>`, item.Title, item.Content)

		epub.AddSubSection(sectionPath, subSectionBody, item.Title, "", "")
	}

	errEpub := epub.Write(fmt.Sprintf("%s/%s.epub", folder, content.Title))
	if errEpub != nil {
		log.Fatal(errEpub)
	}

	location := "./cmd/epub/assets/cover.jpg"
	os.Remove(location)
}

func RetrieveImage(source string) (string, error) {
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

	defer response.Body.Close()

	//open a file for writing
	location := "./cmd/epub/assets/cover.jpg"
	file, err := os.Create(location)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	io.Copy(file, response.Body)

	return location, nil
}
