package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"net/http"
)

// TODO : add support for control messages or restart, apply-changes, rollback-changes

const (
	StartupEvent           = "event:startup"
	ShutdownEvent          = "event:shutdown"
	RestartEvent           = "event:restart"
	ProcessEvent           = "event:process"
	HostStartupEvent       = "event:host-startup"
	PingEvent              = "event:ping"
	ReconfigureEvent       = "event:reconfigure"
	ChangesetApplyEvent    = "event:changeset-apply"
	ChangesetRollbackEvent = "event:changeset-rollback"
	DataEvent              = "event:data"
	StatusEvent            = "event:status"
	DataChangeEvent        = "event:data-change"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive

	ContentType       = "Content-Type"
	XRelatesTo        = "x-relates-to"
	XMessageId        = "x-message-id"
	XTo               = "x-to"
	XFrom             = "x-from"
	XEvent            = "x-event"
	XChannel          = "x-channel"
	XAgentId          = "x-agent-id"
	XForwardTo        = "x-forward-to"
	ContentTypeStatus = "application/status"
	ContentTypeConfig = "application/config"
	ChannelData       = "DATA"
	ChannelControl    = "CTRL"
	ChannelStatus     = "STATUS"
	ChannelNone       = "NONE"
)

// SendFunc - uniform interface for messaging
type SendFunc func(m *Message)

// Handler - uniform interface for message handling
type Handler func(msg *Message)

// Map - map of messages
type Map map[string]*Message

// Message - message
type Message struct {
	Header  http.Header
	Body    any
	ReplyTo Handler
}

func NewMessage(channel, to, from, event string) *Message {
	m := new(Message)
	if len(channel) == 0 {
		channel = ChannelNone
	}
	m.Header = make(http.Header)
	m.Header.Add(XChannel, channel)
	m.Header.Add(XTo, to)
	m.Header.Add(XFrom, from)
	m.Header.Add(XEvent, event)
	return m
}

func NewControlMessage(to, from, event string) *Message {
	return NewMessage(ChannelControl, to, from, event)
}

func NewDataMessage(to, from, event string) *Message {
	return NewMessage(ChannelData, to, from, event)
}

func NewStatusMessage(to, from string, status *core.Status) *Message {
	m := NewMessage(ChannelStatus, to, from, StatusEvent)
	m.SetContent(ContentTypeStatus, status)
	m.Body = status
	return m
}

func NewMessageWithReply(channel, to, from, event string, replyTo Handler) *Message {
	m := NewMessage(channel, to, from, event)
	m.ReplyTo = replyTo
	return m
}

func NewMessageWithStatus(channel, to, from, event string, status *core.Status) *Message {
	m := NewMessage(channel, to, from, event)
	m.SetContent(ContentTypeStatus, status)
	m.Body = status
	return m
}

func (m *Message) String() string {
	status := m.Status()
	if status != nil {
		return fmt.Sprintf("[chan:%v] [from:%v] [to:%v] [%v] [status:%v]", m.Channel(), m.From(), m.To(), m.Event(), status)
	}
	return fmt.Sprintf("[chan:%v] [from:%v] [to:%v] [%v]", m.Channel(), m.From(), m.To(), m.Event())
}

func (m *Message) To() string {
	return m.Header.Get(XTo)
}

func (m *Message) SetTo(uri string) {
	m.Header.Set(XTo, uri)
}

func (m *Message) From() string {
	return m.Header.Get(XFrom)
}

func (m *Message) SetFrom(uri string) {
	m.Header.Set(XFrom, uri)
}

func (m *Message) Event() string {
	return m.Header.Get(XEvent)
}

func (m *Message) RelatesTo() string {
	return m.Header.Get(XRelatesTo)
}

func (m *Message) ForwardTo() string {
	return m.Header.Get(XForwardTo)
}

func (m *Message) Channel() string {
	return m.Header.Get(XChannel)
}

func (m *Message) IsContentType(ct string) bool {
	return m.Header.Get(ContentType) == ct
}

func (m *Message) Status() *core.Status {
	ct := m.Header.Get(ContentType)
	if ct != ContentTypeStatus || m.Body == nil {
		return nil
	}
	if s, ok := m.Body.(*core.Status); ok {
		return s
	}
	return nil
}

func (m *Message) Config() map[string]string {
	ct := m.Header.Get(ContentType)
	if ct != ContentTypeConfig || m.Body == nil {
		return nil
	}
	if m, ok := m.Body.(map[string]string); ok {
		return m
	}
	return nil
}

func (m *Message) Content() (string, any, bool) {
	if m.Body == nil {
		return "", nil, false
	}
	ct := m.Header.Get(ContentType)
	if len(ct) == 0 {
		return "", nil, false
	}
	return ct, m.Body, true
}

func (m *Message) SetContent(contentType string, content any) error {
	if len(contentType) == 0 {
		return errors.New("error: content type is empty")
	}
	if content == nil {
		return errors.New("error: content is nil")
	}
	m.Body = content
	m.Header.Add(ContentType, contentType)
	return nil
}

func (m *Message) SetContentType(contentType string) {
	if len(contentType) == 0 {
		return //errors.New("error: content type is empty")
	}
	m.Header.Add(ContentType, contentType)
}

func (m *Message) ContentType() string {
	return m.Header.Get(ContentType)
}

// SendReply - function used by message recipient to reply with a Status
func SendReply(msg *Message, status *core.Status) {
	if msg == nil || msg.ReplyTo == nil {
		return
	}
	m := NewMessageWithStatus(ChannelNone, msg.From(), msg.To(), msg.Event(), status)
	m.Header.Add(XRelatesTo, msg.RelatesTo())
	msg.ReplyTo(m)
}
