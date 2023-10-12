package dto

import "fmt"

type ChaptersCreate struct {
	Audiobook *Audiobook
}

func (c *ChaptersCreate) String() string {
	return fmt.Sprintf("ChaptersCreate %s", c.Audiobook.String())
}

type AddPartCommand struct {
	Part *Part
}

func (c *AddPartCommand) String() string {
	return fmt.Sprintf("AddPartCommand: %d", c.Part.Number)
}

type AddChapterCommand struct {
	Chapter *Chapter
}

func (c *AddChapterCommand) String() string {
	return fmt.Sprintf("AddChapterCommand: %s", c.Chapter.Name)
}

type ChaptersReady struct {
	Audiobook *Audiobook
}

func (c *ChaptersReady) String() string {
	return fmt.Sprintf("ChaptersReady: %s", c.Audiobook.String())
}
