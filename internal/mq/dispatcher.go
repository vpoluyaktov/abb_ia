package mq

import (
	"container/list"
	"reflect"
	"sync"
	"time"

	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
)

/**
 * Dispatcher is a struct that provides a mechanism for dispatching messages
 * to multiple recipients.
 *
 * @param mu - A mutex to ensure thread safety.
 * @param recipients - A map of strings to message queues.
 * @param listeners - A map of strings to callback functions.
 *
 * @returns Dispatcher - A struct that provides a mechanism for dispatching
 * messages to multiple recipients.
 *
 * This code is useful for dispatching messages to multiple recipients in a
 * thread-safe manner. The mutex ensures that only one thread can access the
 * data at a time, while the maps provide a way to store and access the
 * recipients and listeners. The callback functions allow for custom
 * functionality to be executed when a message is received.
 */
type Dispatcher struct {
	mu         sync.Mutex
	recipients map[string]messageQueue
	listeners  map[string]CallBackFunc
}

type messageQueue struct {
	messages *list.List
}

const PullFrequency = 10 * time.Millisecond

type CallBackFunc func(*Message)

/**
 * Creates a new Dispatcher instance.
 *
 * @returns A new Dispatcher instance.
 *
 * This function is useful for creating a new Dispatcher instance which is used to
 * manage message queues and callbacks. The instance is initialized with empty
 * maps for recipients and listeners.
 */
func NewDispatcher() *Dispatcher {
	d := &Dispatcher{}
	d.recipients = make(map[string]messageQueue)
	d.listeners = make(map[string]CallBackFunc)
	return d
}

/**
 * Sends a message to a recipient.
 * @param from The sender of the message.
 * @param to The recipient of the message.
 * @param dto The data transfer object (DTO) to be sent.
 * @param async Whether the message should be sent asynchronously.
 *
 * This function is useful for sending messages between different components of an application. If the message is sent asynchronously, it is pushed to a queue and the recipient will receive it when they are ready. If the message is sent synchronously, the recipient's method is called in blocking mode.
 */
func (d *Dispatcher) SendMessage(from string, to string, dto dto.Dto, async bool) {

	m := &Message{}
	m.From = from
	m.To = to
	m.Dto = dto
	m.Async = async

	if async {
		// push message to queue
		if _, ok := d.recipients[m.To]; !ok {
			d.recipients[m.To] = messageQueue{list.New()}
		}
		// check if such message is already in queue
		d.mu.Lock()
		if !d.messageExists(m) {
			d.recipients[m.To].messages.PushBack(m)
			logger.Debug("MQ <-- async " + m.String())
		}
		d.mu.Unlock()
	} else if _, ok := d.listeners[m.To]; ok {
		// call recepient method in blocking mode
		logger.Debug("MQ <-- sync  " + m.String())
		d.listeners[m.To](m)
		logger.Debug("MQ  sync --> " + m.String())
	}
}

/**
 * Retrieves a message from the Dispatcher for the given recipient.
 *
 * @param recipient  The recipient of the message.
 * @returns The message, or `undefined` if no message is available.
 *
 * This function is useful for retrieving messages from the Dispatcher for a given recipient. It uses a lock to ensure thread safety when accessing the message queue. It also logs the message for debugging purposes.
 */
func (d *Dispatcher) GetMessage(recipient string) *Message {
	var m *Message
	d.mu.Lock()
	if _, ok := d.recipients[recipient]; ok {
		e := d.recipients[recipient].messages.Front()
		if e != nil {
			d.recipients[recipient].messages.Remove(e)
			m = e.Value.(*Message)
			logger.Debug("MQ async --> " + m.String())
		}
	}
	d.mu.Unlock()
	return m
}

/**
 * Checks if a given message exists in the list of recipients.
 * @param m - The message to check for.
 * @returns A boolean indicating if the message exists.
 *
 * This function is useful for checking if a given message exists in the list of recipients. It iterates through the list of recipients and checks if the given message is equal to any of the messages in the list. If it is, it returns true, otherwise it returns false.
 */
func (d *Dispatcher) messageExists(m *Message) bool {
	found := false
	for e := d.recipients[m.To].messages.Front(); e != nil; e = e.Next() {
		mm := e.Value.(*Message)
		if reflect.DeepEqual(m, mm) {
			found = true
			break
		}
	}
	return found
}

func (d *Dispatcher) RegisterListener(recipient string, callBackFunc CallBackFunc) {
	d.listeners[recipient] = callBackFunc
}
