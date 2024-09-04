package messaging

import (
	"errors"
	"github.com/advanced-go/stdlib/core"
)

const (
	ChannelSize = 16
)

// OnShutdown - add functions to be run on shutdown
type OnShutdown interface {
	Add(func())
}

// Agent - intelligent agent
// TODO : Track agent assignment as part of the URI or separate identifier??
// //Uri() string
//
//	//Message(m *Message)
//	Track agent NID or class/type?
type Agent interface {
	Mailbox
	Run()
	Shutdown()
}

type OpsAgent interface {
	Agent
	core.ErrorHandler
	// AddActivity - context + content
	AddActivity(agentId string, content any)
}

// AgentFunc - agent function
type AgentFunc func(uri string, ctrl, data <-chan *Message, state any)

type agent struct {
	running  bool
	uri      string
	state    any
	ctrl     chan *Message
	data     chan *Message
	status   chan *Message
	run      AgentFunc
	shutdown func()
}

func NewAgent(uri string, run AgentFunc, state any) (Agent, error) {
	return newAgent(uri, make(chan *Message, ChannelSize), make(chan *Message, ChannelSize), make(chan *Message, ChannelSize), run, state)
}

func NewAgentWithChannels(uri string, ctrl, data, status chan *Message, run AgentFunc, state any) (Agent, error) {
	return newAgent(uri, ctrl, data, status, run, state)
}

func newAgent(uri string, ctrl, data, status chan *Message, run AgentFunc, state any) (Agent, error) {
	if len(uri) == 0 {
		return nil, errors.New("error: agent URI is empty")
	}
	if run == nil {
		return nil, errors.New("error: agent AgentFunc is nil")
	}
	if ctrl == nil {
		ctrl = make(chan *Message, ChannelSize)
	}
	a := new(agent)
	a.uri = uri
	a.state = state
	a.ctrl = ctrl
	a.data = data
	a.status = status
	a.run = run
	return a, nil
}

// Uri - identity
func (a *agent) Uri() string {
	return a.uri
}

// String - identity
func (a *agent) String() string {
	return a.uri
}

// Message - message an agent
func (a *agent) Message(msg *Message) {
	Mux(msg, a.ctrl, a.data, a.status)
}

// Run - run the agent
func (a *agent) Run() {
	if a.running {
		return
	}
	a.running = true
	go a.run(a.uri, a.ctrl, a.data, a.state)
}

// Shutdown - shutdown the agent
func (a *agent) Shutdown() {
	if !a.running {
		return
	}
	a.running = false
	if a.shutdown != nil {
		a.shutdown()
	}
	a.Message(NewControlMessage(a.uri, a.uri, ShutdownEvent))
}

// Add - add a shutdown function
func (a *agent) Add(f func()) {
	a.shutdown = AddShutdown(a.shutdown, f)
}

// Mux - multiplex a message over channels
func Mux(msg *Message, ctrl, data, status chan *Message) {
	if msg == nil {
		return
	}
	switch msg.Channel() {
	case ChannelControl:
		if ctrl != nil {
			ctrl <- msg
		}
	case ChannelData:
		if data != nil {
			data <- msg
		}
	case ChannelStatus:
		if status != nil {
			status <- msg
		}
	default:
	}
}

func AddShutdown(curr, next func()) func() {
	if next == nil {
		return nil
	}
	if curr == nil {
		curr = next
	} else {
		// !panic
		prev := curr
		curr = func() {
			prev()
			next()
		}
	}
	return curr
}

//	return
//}
/*
	if msg1, ok := t.(*core.Status); ok {
		if status != nil {
			status <- msg1
		}
	}

*/
//}

/*
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in agent.Shutdown() : %v\n", r)
		}
	}()

*/
