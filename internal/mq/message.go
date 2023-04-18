package mq

import (
	"fmt"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/dto"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)

type Message struct {
	From  string
	To    string
	Type  string
	Dto   dto.Dto
	Async bool
}

func (m *Message) String() string {
	return fmt.Sprintf("Message - From:" + m.From + ", To:" + m.To + " " + m.Dto.String())
}

func (m *Message) DtoCastError() {
	logger.Error("MQ DTO cast error. From: " + m.From + ", To: " + m.To + ", DTO Type: " + m.Type + ", object type: " + fmt.Sprintf("%T", m.Dto))
}

func (m *Message) UnsupportedTypeError() {
	logger.Error("Unsupported message type: " + m.Type + ", sent From: " + m.From + ", To: " + m.To)
}
