package dto

import (
	"abb_ia/internal/config"
	"fmt"
)

type SearchCondition struct {
	Author    string
	Title     string
	SortBy    string
	SortOrder string
}

type SearchCommand struct {
	Condition SearchCondition
}

func (c *SearchCommand) String() string {
	return fmt.Sprintf("SearchCommand: Author: %s, Title: %s", c.Condition.Author, c.Condition.Title)
}

type UpdateSearchConfigCommand struct {
	Config config.Config
}

func (c *UpdateSearchConfigCommand) String() string {
	return fmt.Sprintf("UpdateSearchConfigCommand: Creator: %s, Title: %s", c.Config.DefaultAuthor, c.Config.DefaultTitle)
}

type GetNextPageCommand struct {
	Condition SearchCondition
}

func (c *GetNextPageCommand) String() string {
	return fmt.Sprintf("GetNextPageCommand: Author: %s, Title: %s", c.Condition.Author, c.Condition.Title)
}

type SearchComplete struct {
	Condition SearchCondition
}

func (c *SearchComplete) String() string {
	return fmt.Sprintf("SearchComplete: Author: %s, Title: %s", c.Condition.Author, c.Condition.Title)
}

type NothingFoundError struct {
	Condition SearchCondition
}

func (c *NothingFoundError) String() string {
	return fmt.Sprintf("NothingFoundError: Author: %s, Title: %s", c.Condition.Author, c.Condition.Title)
}

type LastPageMessage struct {
	Condition SearchCondition
}

func (c *LastPageMessage) String() string {
	return fmt.Sprintf("LastPageMessage: Author: %s, Title: %s", c.Condition.Author, c.Condition.Title)
}

type SearchProgress struct {
	ItemsTotal   int
	ItemsFetched int
}

func (p *SearchProgress) String() string {
	return fmt.Sprintf("SearchProgress: %d from %d", p.ItemsFetched, p.ItemsTotal)
}
