package controller

import (
	"time"

	"github.com/vpoluyaktov/audiobook_creator_IA/internal/mq"
)

type controller interface {
	checkMQ()
}

type Conductor struct {
	dispatcher  *mq.Dispatcher
	controllers []controller
}

func NewConductor(dispatcher *mq.Dispatcher) *Conductor {
	c := &Conductor{}
	c.dispatcher = dispatcher
	c.controllers = append(c.controllers, NewSearchController(c.dispatcher))
	return c
}

func (c *Conductor) startEventListener() {
	for {
		for _, p := range c.controllers {
			p.checkMQ()
		}
		time.Sleep(mq.PullFrequency)
	}
}

func (c *Conductor) Run() {
	go c.startEventListener()
}
