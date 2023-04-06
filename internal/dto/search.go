package dto

const SearchCommandType = "dto.SearchCommand"

type SearchCommand struct {
	SearchCondition string
}

const SearchResultType = "dto.SearchResult"

type SearchResult struct {
	ItemName string
}
