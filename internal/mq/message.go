package mq

import (
	"fmt"
	"time"

	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
)

// MessagePriority defines the priority level of a message
type MessagePriority int

const (
	PriorityLow MessagePriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical

	// PullFrequency defines how often to check for new messages
	PullFrequency = 100 * time.Millisecond
)

// Message represents a message in the system
type Message struct {
	From       string
	To         string
	Dto        dto.Dto
	Async      bool
	Priority   MessagePriority `json:"priority"`
	CreatedAt  time.Time      `json:"created_at"`
	ExpiresAt  time.Time      `json:"expires_at"`
	RetryCount int            `json:"retry_count"`
	MaxRetries int            `json:"max_retries"`
	MessageID  string         `json:"message_id"`
	Metadata   map[string]any `json:"metadata"`
}

func (m *Message) String() string {
	return fmt.Sprintf("Message [From:%s, To:%s, ID:%s, Priority:%v] %s", m.From, m.To, m.MessageID, m.Priority, m.Dto.String())
}

func (m *Message) UnsupportedTypeError(reporter string) {
	logger.Error(fmt.Sprintf("%s: Unsupported message type: %T, sent From: %s, To: %s", reporter, m.Dto, m.From, m.To))
}

// CallBackFunc is a function type for message handlers
type CallBackFunc func(msg *Message)
