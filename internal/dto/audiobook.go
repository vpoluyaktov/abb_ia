package dto

import (
	"encoding/json"
	"fmt"

	"abb_ia/internal/config"
)

type Audiobook struct {
	Author        string
	Title         string
	Description   string
	Genre         string
	Series        string
	SeriesNo      string
	Narrator      string
	Year          string
	CoverURL      string
	CoverFile     string
	IaURL         string
	Copyright     string
	LicenseUrl    string
	OutputDir     string
	Mp3Files      []Mp3File
	TotalDuration float64
	TotalSize     int64
	Parts         []Part
	IAItem        *IAItem
	Config        *config.Config
}

type Part struct {
	Number       int
	AACFile      string
	M4BFile      string
	FListFile    string
	MetadataFile string
	Format       string
	Size         int64
	Duration     float64
	Chapters     []Chapter
}

type Chapter struct {
	Number   int
	Name     string
	Size     int64
	Duration float64
	Start    float64
	End      float64
	Files    []Mp3File
}

type Mp3File struct {
	Number   int
	FileName string
	Size     int64
	Duration float64
}

func (ab *Audiobook) String() string {
	return fmt.Sprintf("%T: %s", ab, ab.Title)
}

func (ab *Audiobook) GetChapter(chapterNumber int) (*Chapter, error) {
	for _, part := range ab.Parts {
		for _, chapter := range part.Chapters {
			if chapter.Number == chapterNumber {
				return &chapter, nil
			}
		}
	}
	return nil, fmt.Errorf("no chapter found")
}

func (ab *Audiobook) SetChapter(chapterNumber int, ch Chapter) error {
	for partNo := range ab.Parts {
		part := &ab.Parts[partNo]
		for chapterNo := range part.Chapters {
			chapter := &part.Chapters[chapterNo]
			if chapter.Number == chapterNumber {
				ab.Parts[partNo].Chapters[chapterNo] = ch
				return nil
			}
		}
	}
	return fmt.Errorf("no chapter found")
}

func (ab *Audiobook) GetCopy() (*Audiobook, error) {
	// Convert the source struct to JSON
	jsonBytes, err := json.Marshal(ab)
	if err != nil {
		return nil, err
	}

	// Create a new destination struct
	destination := &Audiobook{}

	// Convert the JSON back to the destination struct
	err = json.Unmarshal(jsonBytes, destination)
	if err != nil {
		return nil, err
	}

	return destination, nil
}
