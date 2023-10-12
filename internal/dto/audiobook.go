package dto

import "fmt"

type Audiobook struct {
	Author      string
	Title       string
	Description string
	Cover       string
	Series      string
	SeriesNo    int
	Copyright   string
	Parts       []Part
	IAItem      *IAItem
}

type Part struct {
	Number   int
	FileName string
	Size     int64
	Duration float64
	Chapters []Chapter
	Files    []Mp3File
}

type Chapter struct {
	Number   int
	Name     string
	Size     int64
	Duration float64
	Start    float64
	End      float64
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

func (ab *Audiobook) SetChapter(chapterNumber int, ch *Chapter) error {
	for _, part := range ab.Parts {
		for _, chapter := range part.Chapters {
			if chapter.Number == chapterNumber {
				chapter = *ch
				return nil
			}
		}
	}
	return fmt.Errorf("no chapter found")
}