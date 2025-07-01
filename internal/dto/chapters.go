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

type SearchReplaceDescriptionCommand struct {
	Audiobook  *Audiobook
	SearchStr  string
	ReplaceStr string
}

func (c *SearchReplaceDescriptionCommand) String() string {
	return fmt.Sprintf("SearchReplaceDescriptionCommand: %s/%s", c.SearchStr, c.SearchStr)
}

type RefreshDescriptionCommand struct {
	Audiobook *Audiobook
}

func (c *RefreshDescriptionCommand) String() string {
	return fmt.Sprintf("RefreshDescriptionCommand: %s", c.Audiobook.String())
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


type RecalculatePartsCommand struct {
	Audiobook *Audiobook
}

func (c *RecalculatePartsCommand) String() string {
	return fmt.Sprintf("RecalculatePartsCommand: %s", c.Audiobook.String())
}

type UseMP3NamesCommand struct {
	Audiobook *Audiobook
}

func (c *UseMP3NamesCommand) String() string {
	return fmt.Sprintf("UseMP3NamesCommand: %s", c.Audiobook.String())
}

// SortChaptersCommand represents a command to sort chapters by name
type SortChaptersCommand struct {
	Audiobook *Audiobook
}

func (c *SortChaptersCommand) String() string {
	return fmt.Sprintf("SortChaptersCommand: %s", c.Audiobook.String())
}
