package flags

import (
	"fmt"
)

type Web string

type NovelData struct {
	WebName          string
	Title            string
	Url              string
	AvailableChapter string
}

type NovelInfo struct {
	WebName   Web
	SearchUrl string
}

var WebInfo = []NovelInfo{
	{
		WebName:   "Novelhall",
		SearchUrl: "https://www.novelhall.com/index.php?s=so&module=book&keyword=%s",
	},
	{
		WebName:   "1stKissNovel",
		SearchUrl: "https://1stkissnovel.org/?s=%s&post_type=wp-manga",
	},
}

// var AvailableWeb = fmt.Sprintf("available Web for scraping: %s", strings.Join(ListWeb, ", "))
var errorScraping = "Error Scraping"

func (w Web) String() string {
	return string(w)
}

func (w *Web) Type() string {
	return "Available Web"
}

func (w *Web) Set(value string) error {
	for _, web := range WebInfo {
		if string(web.WebName) == value {
			*w = Web(value)
			return nil
		}
	}

	// return fmt.Errorf(AvailableWeb)
	return fmt.Errorf(errorScraping)
}
