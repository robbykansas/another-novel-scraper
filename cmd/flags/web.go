package flags

import (
	"fmt"
	"robbykansas/another-novel-scraper/cmd/models"
)

type Web models.Web

var AllowedWeb []string

var errorScraping = "Error Scraping"

func (n Web) String() string {
	return string(n)
}

func (n *Web) Type() string {
	return "Available Web"
}

func (n *Web) Set(value string) error {
	for _, web := range AllowedWeb {
		if string(web) == value {
			*n = Web(value)
			return nil
		}
	}

	// return fmt.Errorf(AvailableWeb)
	return fmt.Errorf(errorScraping)
}
