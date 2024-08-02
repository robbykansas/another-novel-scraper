package flags

import (
	"fmt"
	"strings"
)

type Web string
type SearchUrl string

type NovelData struct {
	WebName          string
	Title            string
	Url              string `json:"-"`
	AvailableChapter string
}

type ListWeb struct {
	WebName string
}

const (
	Syosetu   Web = "Syosetu"
	NovelFull Web = "NovelFull"
)

const (
	Novelhall SearchUrl = "https://www.novelhall.com/index.php?s=so&module=book&keyword=%s"
)

var listSearchUrl = []string{
	string(Novelhall),
}

var listWeb = []string{string(Syosetu), string(NovelFull)}

var AvailableWeb = fmt.Sprintf("available Web for scraping: %s", strings.Join(listWeb, ", "))

func (w Web) String() string {
	return string(w)
}

func (w *Web) Type() string {
	return "Available Web"
}

func (w *Web) Set(value string) error {
	for _, web := range listWeb {
		if web == value {
			*w = Web(value)
			return nil
		}
	}

	return fmt.Errorf(AvailableWeb)
}
