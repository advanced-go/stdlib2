package messaging

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"net/http"
)

func handler(msg *Message) {
	fmt.Printf(msg.Event())
}

func ExampleReplyTo() {
	msg := NewMessageWithReply(ChannelNone, "test", "", "startup", handler)
	SendReply(msg, core.StatusOK())

	msg = NewMessage(ChannelNone, "test", "", "startup")
	SendReply(msg, core.StatusOK())

	//Output:
	//startup

}

func ExampleStatusContent() {
	status := core.NewStatus(http.StatusTeapot)
	m := NewMessage(ChannelNone, "to", "from", StartupEvent)

	err := m.SetContent("", status)
	fmt.Printf("test: SetContent(\"\",status) -> [%v]\n", err)

	err = m.SetContent(ContentTypeStatus, nil)
	fmt.Printf("test: SetContent(\"%v\",nil) -> [%v]\n", ContentTypeStatus, err)

	err = m.SetContent(ContentTypeStatus, status)
	fmt.Printf("test: SetContent(\"%v\",status) -> [%v]\n", ContentTypeStatus, err)

	m = NewMessage(ChannelNone, "to", "from", StartupEvent)
	ct, body, ok := m.Content()
	fmt.Printf("test: Content() -> [ct:%v] [body:%v] [ok:%v]\n", ct, body, ok)

	m.SetContent(ContentTypeStatus, status)
	ct, body, ok = m.Content()
	fmt.Printf("test: Content() -> [ct:%v] [body:%v] [ok:%v]\n", ct, body, ok)

	s := m.Status()
	fmt.Printf("test: Status() -> [body:%v]\n", s)

	m = NewMessageWithStatus(ChannelNone, "to", "from", StartupEvent, status)
	s = m.Status()
	fmt.Printf("test: NewMessageWithStatus() -> [body:%v]\n", s)

	//Output:
	//test: SetContent("",status) -> [error: content type is empty]
	//test: SetContent("application/status",nil) -> [error: content is nil]
	//test: SetContent("application/status",status) -> [<nil>]
	//test: Content() -> [ct:] [body:<nil>] [ok:false]
	//test: Content() -> [ct:application/status] [body:I'm A Teapot] [ok:true]
	//test: Status() -> [body:I'm A Teapot]
	//test: NewMessageWithStatus() -> [body:I'm A Teapot]

}

func ExampleConfigContent() {
	cfg := make(map[string]string)
	cfg["uri"] = "http://www.google/com"
	m := NewMessage(ChannelNone, "to", "from", StartupEvent)

	ct, body, ok := m.Content()
	fmt.Printf("test: Content() -> [ct:%v] [body:%v] [ok:%v]\n", ct, body, ok)

	m.SetContent(ContentTypeConfig, cfg)
	ct, body, ok = m.Content()
	fmt.Printf("test: Content() -> [ct:%v] [body:%v] [ok:%v]\n", ct, body, ok)

	s := m.Config()
	fmt.Printf("test: Status() -> [body:%v]\n", s)

	//Output:
	//test: Content() -> [ct:] [body:<nil>] [ok:false]
	//test: Content() -> [ct:application/config] [body:map[uri:http://www.google/com]] [ok:true]
	//test: Status() -> [body:map[uri:http://www.google/com]]

}
