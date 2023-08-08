package dto

import "fmt"

type Audiobook struct {
	Title    string
	Author   string
	Chapters []Chapter
	IAItem   *IAItem
}

type Chapter struct {
	Number int
	Start float64
	End   float64
	Duration float64
	Name  string
}

func (ab *Audiobook) String() string {
	return fmt.Sprintf("%T: %s", ab, ab.Title)
}
