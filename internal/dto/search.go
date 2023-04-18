package dto

import "fmt"

// format list ranged by priority
var FormatList = []string{"16Kbps MP3", "24Kbps MP3", "32Kbps MP3", "40Kbps MP3", "48Kbps MP3", "56Kbps MP3", "64Kbps MP3", "80Kbps MP3", "96Kbps MP3", "112Kbps MP3", "128Kbps MP3", "144Kbps MP3", "160Kbps MP3", "224Kbps MP3", "256Kbps MP3", "320Kbps MP3", "VBR MP3"}

const SearchCommandType = "dto.SearchCommand"
type SearchCommand struct {
	SearchCondition string
}
func (c *SearchCommand) String() string {
	return fmt.Sprintf("SearchCommand: %s", c.SearchCondition)
}

const IAItemType = "dto.IAItem"
type IAItem struct {
	ID           string
	Title        string
	Creator      string
	Description  string
	Server       string
	Dir          string
	FilesCount   int
	TotalLength  float64
	TotalLengthH string
	TotalSize    int64
	TotalSizeH   string
	Files        []File
}
func (i *IAItem) String() string {
	return fmt.Sprintf("%T: %s", i, i.Title)
}

type File struct {
	Name   string
	Format string
	Length float64
	LengthH string
	Size   int64
	SizeH  string
}
func (f *File) String() string {
	return fmt.Sprintf("%T: %s", f, f.Name)
}