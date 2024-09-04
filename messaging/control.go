package messaging

import (
	"errors"
	"github.com/advanced-go/stdlib/core"
)

// NewControlAgent - create an agent that only listens on a control channel, and has a default AgentRun func
func NewControlAgent(uri string, ctrlHandler Handler) (Agent, error) {
	if ctrlHandler == nil {
		return nil, errors.New("error: control agent message handler is nil")
	}
	return NewAgentWithChannels(uri, nil, nil, nil, controlAgentRun, ctrlHandler)
}

// controlAgentRun - a simple run function that only handles control messages, and dispatches via a message handler
func controlAgentRun(_ string, ctrl, _ <-chan *Message, state any) {
	if state == nil || ctrl == nil {
		return
	}
	var ctrlHandler Handler
	if h, ok := state.(Handler); ok {
		ctrlHandler = h
	} else {
		return
	}
	for {
		select {
		case msg, open := <-ctrl:
			if !open {
				return
			}
			switch msg.Event() {
			case ShutdownEvent:
				ctrlHandler(NewMessageWithStatus(ChannelControl, "", "", msg.Event(), core.StatusOK()))
				//m.close()
				return
			default:
				ctrlHandler(msg)
			}
		default:
		}
	}
}
