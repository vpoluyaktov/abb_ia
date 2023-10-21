package dto

import "fmt"

type IAItem struct {
	ID           string
	Title        string
	Creator      string
	Description  string
	Cover        string
	Server       string
	Dir          string
	FilesCount   int
	TotalLength  float64
	TotalSize    int64
	Files        []File
}

func (i *IAItem) String() string {
	return fmt.Sprintf("%T: %s", i, i.Title)
}

type File struct {
	Name    string
	Format  string
	Length  float64
	Size    int64
}

func (f *File) String() string {
	return fmt.Sprintf("%T: %s", f, f.Name)
}
