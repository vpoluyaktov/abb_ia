package dto

type Dto interface {
}

const CommandType = "dto.Command"
type Command struct {
	Command string
}