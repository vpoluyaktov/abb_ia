package dto

import "fmt"

type Audiobook struct {
	Title string
	Author string
	IAItem *IAItem
}

func (ab *Audiobook) String() string {
	return fmt.Sprintf("%T: %s", ab, ab.Title)
}

