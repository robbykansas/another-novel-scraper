package flags

import (
	"fmt"
)

type Web string

var AllowedWeb []string

// var AvailableWeb = fmt.Sprintf("available Web for scraping: %s", strings.Join(ListWeb, ", "))
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
