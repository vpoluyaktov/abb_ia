package mq

import (
	"fmt"

	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
)

type Message struct {
	From  string
	To    string
	Dto   dto.Dto
	Async bool
}

func (m *Message) String() string {
	return fmt.Sprintf("Message [From:" + m.From + ", To:" + m.To + " " + m.Dto.String() + "]")
}

func (m *Message) UnsupportedTypeError(reporter string) {
	logger.Error(reporter + ": Unsupported message type: " + fmt.Sprintf("%T", m.Dto) + ", sent From: " + m.From + ", To: " + m.To)
}
