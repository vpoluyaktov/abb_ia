package controller

import (
	"time"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/event"
)

type processor interface {
	readMessages()
}

type Controller struct {
	dispatcher *event.Dispatcher
	processors []processor
}

func NewController(dispatcher *event.Dispatcher) *Controller {
	c := &Controller{}
	c.dispatcher = dispatcher
	c.processors = append(c.processors, NewSearchProcessor(c.dispatcher))
	return c
}

func (c *Controller) startEventListener() {
	for {
		for _, p := range c.processors {
			p.readMessages()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Controller) Run() {
	go c.startEventListener()
}
