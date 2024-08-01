package flags

import (
	"fmt"
	"strings"
)

type Web string

type NovelData struct {
	WebName          string
	AvailableChapter string
}

const (
	Syosetu   Web = "Syosetu"
	NovelFull Web = "NovelFull"
)

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
