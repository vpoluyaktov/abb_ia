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

type SearchReplaceChaptersCommand struct {
	Audiobook  *Audiobook
	SearchStr  string
	ReplaceStr string
}

func (c *SearchReplaceChaptersCommand) String() string {
	return fmt.Sprintf("SearchReplaceChaptersCommand: %s/%s", c.SearchStr, c.SearchStr)
}

type JoinChaptersCommand struct {
	Audiobook *Audiobook
}

func (c *JoinChaptersCommand) String() string {
	return fmt.Sprintf("JoinChaptersCommand: %s", c.Audiobook.String())
}

type RefreshChaptersCommand struct {
	Audiobook *Audiobook
}

func (c *RefreshChaptersCommand) String() string {
	return fmt.Sprintf("RefreshChaptersCommand: %s", c.Audiobook.String())
}