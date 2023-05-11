package dto

import "fmt"

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
	Name    string
	Format  string
	Length  float64
	LengthH string
	Size    int64
	SizeH   string
}

func (f *File) String() string {
	return fmt.Sprintf("%T: %s", f, f.Name)
}
