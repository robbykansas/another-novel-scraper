package flags

import (
	"fmt"
)

type ChosenTitle string

var AllowedTitle []string

var errorChosenTitle = "Title not available"

func (n ChosenTitle) String() string {
	return string(n)
}

func (n *ChosenTitle) Type() string {
	return "ChosenTitle"
}

func (n *ChosenTitle) Set(value string) error {
	for _, title := range AllowedTitle {
		if title == value {
			*n = ChosenTitle(value)
			return nil
		}
	}

	return fmt.Errorf(errorChosenTitle)
}
