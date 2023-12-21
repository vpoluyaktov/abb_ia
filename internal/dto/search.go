package dto

import "fmt"

type SearchCommand struct {
	SearchCondition string
}

func (c *SearchCommand) String() string {
	return fmt.Sprintf("SearchCommand: %s", c.SearchCondition)
}

type GetNextPageCommand struct {
	SearchCondition string
}

func (c *GetNextPageCommand) String() string {
	return fmt.Sprintf("GetNextPageCommand: %s", c.SearchCondition)
}

type SearchComplete struct {
	SearchCondition string
}

func (c *SearchComplete) String() string {
	return fmt.Sprintf("SearchComplete: %s", c.SearchCondition)
}

type NothingFoundError struct {
	SearchCondition string
}

func (c *NothingFoundError) String() string {
	return fmt.Sprintf("NothingFoundError: %s", c.SearchCondition)
}

type LastPageMessage struct {
	SearchCondition string
}

func (c *LastPageMessage) String() string {
	return fmt.Sprintf("LastPageMessage: %s", c.SearchCondition)
}

type SearchProgress struct {
	ItemsTotal   int
	ItemsFetched int
}

func (p *SearchProgress) String() string {
	return fmt.Sprintf("SearchProgress: %d from %d", p.ItemsFetched, p.ItemsTotal)
}
