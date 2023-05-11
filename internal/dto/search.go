package dto

import "fmt"

// format list ranged by priority
var FormatList = []string{"16Kbps MP3", "24Kbps MP3", "32Kbps MP3", "40Kbps MP3", "48Kbps MP3", "56Kbps MP3", "64Kbps MP3", "80Kbps MP3", "96Kbps MP3", "112Kbps MP3", "128Kbps MP3", "144Kbps MP3", "160Kbps MP3", "224Kbps MP3", "256Kbps MP3", "320Kbps MP3", "VBR MP3"}

type SearchCommand struct {
	SearchCondition string
}

func (c *SearchCommand) String() string {
	return fmt.Sprintf("SearchCommand: %s", c.SearchCondition)
}

type SearchProgress struct {
	ItemsTotal   int
	ItemsFetched int
}

func (p *SearchProgress) String() string {
	return fmt.Sprintf("SearchProgress: %d from %d", p.ItemsFetched, p.ItemsTotal)
}
