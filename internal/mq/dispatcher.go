package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
)

// Dispatcher provides a message dispatcher with advanced features
type Dispatcher struct {
	mu              sync.RWMutex
	messageChannels map[string]chan *Message
	listeners       map[string]CallBackFunc
	metrics         *DispatcherMetrics
	ctx            context.Context
	cancel         context.CancelFunc
}

// RegisterHandler registers a callback function for a specific recipient
func (d *Dispatcher) RegisterHandler(recipient string, handler CallBackFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.listeners[recipient] = handler

	// Create a message channel if it doesn't exist
	if _, exists := d.messageChannels[recipient]; !exists {
		ch := make(chan *Message, 100) // Buffered channel
		d.messageChannels[recipient] = ch
		// Start a goroutine to handle messages for this recipient
		go d.processMessages(recipient, ch)
	}

	logger.Debug(fmt.Sprintf("Registered message handler for recipient: %s", recipient))
}

// RegisterListener registers a listener for a specific recipient
func (d *Dispatcher) RegisterListener(recipient string, handler CallBackFunc) {
	d.RegisterHandler(recipient, handler)
}

// DispatcherMetrics tracks performance and usage metrics
type DispatcherMetrics struct {
	mu                sync.RWMutex
	messagesSent      int64
	messagesProcessed int64
	errors            int64
	startTime         time.Time
}

// NewDispatcher creates a new message dispatcher
func NewDispatcher() *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Dispatcher{
		messageChannels: make(map[string]chan *Message),
		listeners:       make(map[string]CallBackFunc),
		metrics: &DispatcherMetrics{
			startTime: time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// Start the message cleanup goroutine
	go d.cleanupExpiredMessages()
	go d.collectMetrics()

	return d
}

// SendMessage sends a message with advanced features
func (d *Dispatcher) SendMessage(from string, to string, dto dto.Dto, priority MessagePriority) error {
	msg := &Message{
		From:       from,
		To:         to,
		Dto:        dto,
		Async:      true,
		Priority:   priority,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(5 * time.Minute), // Default 5-minute expiration
		MessageID:  generateMessageID(),
		MaxRetries: 3,
		Metadata:   make(map[string]any),
	}

	d.mu.Lock()
	ch, exists := d.messageChannels[to]
	if !exists {
		ch = make(chan *Message, 100) // Buffered channel
		d.messageChannels[to] = ch
		// Start a goroutine to handle messages for this recipient
		go d.processMessages(to, ch)
	}
	d.mu.Unlock()

	select {
	case ch <- msg:
		d.incrementMetric(&d.metrics.messagesSent)
		logger.Debug(fmt.Sprintf("Message sent - ID: %s, From: %s, To: %s, Priority: %v",
			msg.MessageID, from, to, priority))
		return nil
	default:
		return fmt.Errorf("failed to send message to %s: channel full", to)
	}
}

// GetMessage retrieves a message from a recipient's channel
func (d *Dispatcher) GetMessage(recipient string) (*Message, error) {
	d.mu.RLock()
	ch, exists := d.messageChannels[recipient]
	d.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no channel for recipient: %s", recipient)
	}

	select {
	case msg := <-ch:
		return msg, nil
	default:
		// This is normal behavior when polling and no messages are available
		// logger.Debug(fmt.Sprintf("No message available for recipient: %s", recipient))
		return nil, nil
	}
}

func (d *Dispatcher) processMessages(recipient string, ch chan *Message) {
	for {
		select {
		case <-d.ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			if time.Now().After(msg.ExpiresAt) {
				logger.Warn(fmt.Sprintf("Message expired, skipping processing - ID: %s", msg.MessageID))
				continue
			}

			if err := d.handleMessage(msg); err != nil {
				if msg.RetryCount < msg.MaxRetries {
					msg.RetryCount++
					// Exponential backoff for retries
					time.Sleep(time.Duration(msg.RetryCount*msg.RetryCount) * time.Second)
					ch <- msg
				} else {
					logger.Error(fmt.Sprintf("Message processing failed after max retries - ID: %s, Error: %v", msg.MessageID, err))
				}
			} else {
				d.incrementMetric(&d.metrics.messagesProcessed)
			}
		}
	}
}

func (d *Dispatcher) handleMessage(msg *Message) error {
	d.mu.RLock()
	listener, exists := d.listeners[msg.To]
	d.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no listener registered for recipient: %s", msg.To)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Panic recovered in message handler - ID: %s, Panic: %v", msg.MessageID, r))
		}
	}()

	listener(msg)
	return nil
}

func (d *Dispatcher) cleanupExpiredMessages() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.mu.RLock()
			for recipient, ch := range d.messageChannels {
				// Non-blocking cleanup of expired messages
				go func(recipient string, ch chan *Message) {
					for {
						select {
						case msg := <-ch:
							if msg != nil && time.Now().Before(msg.ExpiresAt) {
								// Put non-expired message back
								ch <- msg
							}
						default:
							return
						}
					}
				}(recipient, ch)
			}
			d.mu.RUnlock()
		}
	}
}

func (d *Dispatcher) collectMetrics() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			metrics := d.GetMetrics()
			logger.Debug(fmt.Sprintf("Dispatcher metrics - Sent: %d, Processed: %d, Errors: %d",
				metrics.messagesSent, metrics.messagesProcessed, metrics.errors))
		}
	}
}

func (d *Dispatcher) incrementMetric(metric *int64) {
	d.metrics.mu.Lock()
	*metric++
	d.metrics.mu.Unlock()
}

func (d *Dispatcher) GetMetrics() DispatcherMetrics {
	d.metrics.mu.RLock()
	defer d.metrics.mu.RUnlock()
	return *d.metrics
}

// Shutdown gracefully shuts down the dispatcher
func (d *Dispatcher) Shutdown(timeout time.Duration) error {
	d.cancel()
	
	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan bool)
	go func() {
		d.mu.Lock()
		defer d.mu.Unlock()
		
		// Close all channels
		for _, ch := range d.messageChannels {
			close(ch)
		}
		done <- true
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out")
	case <-done:
		return nil
	}
}

func generateMessageID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}
