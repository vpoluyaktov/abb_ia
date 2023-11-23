package dto

import "fmt"

type IAItem struct {
	ID          string
	Title       string
	Creator     string
	Description string
	CoverUrl    string
	IaURL       string
	LicenseUrl  string
	Server      string
	Dir         string
	TotalLength float64
	TotalSize   int64
	AudioFiles  []AudioFile
	ImageFiles  []ImageFile
}

func (i *IAItem) String() string {
	return fmt.Sprintf("%T: %s", i, i.Title)
}

type AudioFile struct {
	Name   string
	Title  string
	Format string
	Length float64
	Size   int64
}

func (f *AudioFile) String() string {
	return fmt.Sprintf("%T: %s", f, f.Name)
}

type ImageFile struct {
	Name   string
	Format string
	Size   int64
}

func (f *ImageFile) String() string {
	return fmt.Sprintf("%T: %s", f, f.Name)
}
