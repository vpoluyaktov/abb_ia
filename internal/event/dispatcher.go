package event

type Dispatcher struct {
	eventQueue []*Message
	listeners  map[string]Fn
}

type Message struct {
	Sender    string
	Recipient string
	Body      interface{}
	Priority  int
	Async     bool
}

type Fn func(*Message)

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{}
	return d
}

func (d *Dispatcher) SendMessage(message *Message) {
	if !message.Async && hasKey(d.listeners, message.Recipient) {
		d.listeners[message.Recipient](message)
	} else {
		d.eventQueue = append(d.eventQueue, message)
	}
}

func (d *Dispatcher) GetMessage(recipient string) *Message {
	for i, m := range d.eventQueue {
		if m.Recipient == recipient {
			d.removeMessage(i)
			return m
		}
	}
	return nil
}

func (d *Dispatcher) RegisterListener(recipient string, callBackFunc Fn) {
	if d.listeners == nil {
		d.listeners = make(map[string]Fn)
	}
	d.listeners[recipient] = callBackFunc
}

// Check if a map contains a given key
func hasKey(m map[string]Fn, key string) bool {
	_, ok := m[key]
	return ok
}

func (d *Dispatcher) removeMessage(index int)  {
	if index <= len(d.eventQueue) {
		q := make([]*Message, 0)
		q = append(q, d.eventQueue[:index]...)
		q = append(q, d.eventQueue[index+1:]...)
		d.eventQueue = q
	} 
}
