package messaging

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

type Mailbox interface {
	Uri() string
	Message(m *Message)
}

// Exchange - controller2 directory
type Exchange struct {
	m *sync.Map
}

// NewExchange - create a new controller2
func NewExchange() *Exchange {
	e := new(Exchange)
	e.m = new(sync.Map)
	return e
}

// Count - number of agents
func (d *Exchange) Count() int {
	count := 0
	d.m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// List - a list of agent uri's
func (d *Exchange) List() []string {
	var uri []string
	d.m.Range(func(key, value any) bool {
		if str, ok := key.(string); ok {
			uri = append(uri, str)
		}
		return true
	})
	sort.Strings(uri)
	return uri
}

// Send - send a message
func (d *Exchange) Send(msg *Message) error {
	// TO DO : authenticate shutdown control message
	//if msg != nil && msg.Event() == ShutdownEvent {
	//	return nil
	//}
	if msg == nil {
		return errors.New(fmt.Sprintf("error: controller2.Send() failed as message is nil"))
	}
	a := d.Get(msg.To())
	if a == nil {
		return errors.New(fmt.Sprintf("error: controller2.Send() failed as the message To is empty or invalid : [%v]", msg.To()))
	}
	a.Message(msg)
	return nil
}

// Broadcast - broadcast a message to all entries, deleting the entry if the message event is Shutdown
func (d *Exchange) Broadcast(msg *Message) {
	if msg == nil {
		return //errors.New(fmt.Sprintf("error: exchange.Broadcast() failed as message is nil"))
	}
	// TODO : Need to disallow shutdown message??
	if msg.Event() == ShutdownEvent {
		return
	}
	go func() {
		for _, uri := range d.List() {
			a := d.Get(uri)
			if a == nil {
				continue
			}
			a.Message(msg)
			//if msg.Event() == ShutdownEvent {
			//	d.m.Delete(uri)
			//}
		}
	}()
}

// RegisterMailbox - register a mailbox
func (d *Exchange) RegisterMailbox(m Mailbox) error {
	if m == nil {
		return errors.New("error: controller2.Register() agent is nil")
	}
	_, ok := d.m.Load(m.Uri())
	if ok {
		return errors.New(fmt.Sprintf("error: controller2.Register() agent already exists: [%v]", m.Uri()))
	}
	d.m.Store(m.Uri(), m)
	if sd, ok1 := m.(OnShutdown); ok1 {
		sd.Add(func() {
			d.m.Delete(m.Uri())
		})
	}
	return nil
}

// GetMailbox - find a mailbox
func (d *Exchange) GetMailbox(uri string) Mailbox {
	if len(uri) == 0 {
		return nil
	}
	v, ok1 := d.m.Load(uri)
	if !ok1 {
		return nil
	}
	if a, ok2 := v.(Mailbox); ok2 {
		return a
	}
	return nil
}

// Register - register an agent
func (d *Exchange) Register(agent Agent) error {
	if agent == nil {
		return errors.New("error: exchange.Register() agent is nil")
	}
	_, ok := d.m.Load(agent.Uri())
	if ok {
		return errors.New(fmt.Sprintf("error: exchange.Register() agent already exists: [%v]", agent.Uri()))
	}
	d.m.Store(agent.Uri(), agent)
	if sd, ok1 := agent.(OnShutdown); ok1 {
		sd.Add(func() {
			d.m.Delete(agent.Uri())
		})
	}
	return nil
}

// Get - find an agent
func (d *Exchange) Get(uri string) Agent {
	if len(uri) == 0 {
		return nil
	}
	v, ok1 := d.m.Load(uri)
	if !ok1 {
		return nil
	}
	if a, ok2 := v.(Agent); ok2 {
		return a
	}
	return nil
}

// Shutdown - shutdown all agents
func (d *Exchange) Shutdown() {
	go func() {
		for _, uri := range d.List() {
			a := d.Get(uri)
			if a == nil {
				continue
			}
			a.Shutdown()
		}
	}()
}

/*
func (d *controller2) shutdown(uri string) runtime.Status {
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	//for _, e := range d.m {
	//	if e.ctrl != nil {
	//		e.ctrl <- Message{To: e.uri, Event: core.ShutdownEvent}
	//	}
	//}
	m, status := d.get(uri)
	if !status.OK() {
		return status
	}
	if m.data != nil {
		close(m.data)
	}
	if m.ctrl != nil {
		close(m.ctrl)
	}
	d.m.Delete(uri)
	return runtime.StatusOK()
}
*/
