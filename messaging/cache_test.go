package messaging

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
)

const (
	StatusNotProvided = int(95)
)

func ExampleCache_Add() {
	resp := NewCache()

	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-0", StartupEvent, core.NewStatus(StatusNotProvided)))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-1", StartupEvent, core.StatusOK()))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-2", PingEvent, core.NewStatus(StatusNotProvided)))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-3", PingEvent, core.NewStatus(StatusNotProvided)))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-4", PingEvent, core.StatusOK()))

	fmt.Printf("test: count() -> : %v\n", resp.Count())

	m, ok := resp.Get("invalid")
	fmt.Printf("test: Get(%v) -> : [ok:%v] [msg-nil:%v]\n", "invalid", ok, m == nil)

	m, ok = resp.Get("from-uri-3")
	fmt.Printf("test: Get(%v) -> : [ok:%v] [msg-to:%v]\n", "from-uri-3", ok, len(m.To()) > 0)

	fmt.Printf("test: include(%v,%v) -> : %v\n", ShutdownEvent, StatusNotProvided, resp.Include(ShutdownEvent, StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", ShutdownEvent, StatusNotProvided, resp.Exclude(ShutdownEvent, StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", StartupEvent, StatusNotProvided, resp.Include(StartupEvent, StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", StartupEvent, StatusNotProvided, resp.Exclude(StartupEvent, StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", PingEvent, StatusNotProvided, resp.Include(PingEvent, StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", PingEvent, StatusNotProvided, resp.Exclude(PingEvent, StatusNotProvided))

	//Output:
	//test: count() -> : 5
	//test: Get(invalid) -> : [ok:false] [msg-nil:true]
	//test: Get(from-uri-3) -> : [ok:true] [msg-to:true]
	//test: include(event:shutdown,95) -> : []
	//test: exclude(event:shutdown,95) -> : [from-uri-0 from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:startup,95) -> : [from-uri-0]
	//test: exclude(event:startup,95) -> : [from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:ping,95) -> : [from-uri-2 from-uri-3]
	//test: exclude(event:ping,95) -> : [from-uri-0 from-uri-1 from-uri-4]

}

func ExampleCache_Uri() {
	resp := NewCache()

	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-0", StartupEvent, core.NewStatus(StatusNotProvided)))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-1", StartupEvent, core.StatusOK()))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-2", PingEvent, core.NewStatus(StatusNotProvided)))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-3", PingEvent, core.NewStatus(StatusNotProvided)))
	resp.Add(NewMessageWithStatus(ChannelNone, "to-uri", "from-uri-4", PingEvent, core.StatusOK()))

	fmt.Printf("test: count() -> : %v\n", resp.Count())

	m, ok := resp.Get("invalid")
	fmt.Printf("test: Get(%v) -> : [ok:%v] [msg-nil:%v]\n", "invalid", ok, m == nil)

	m, ok = resp.Get("from-uri-3")
	fmt.Printf("test: Get(%v) -> : [ok:%v] [msg-to:%v]\n", "from-uri-3", ok, m.To())

	fmt.Printf("test: include(%v,%v) -> : %v\n", ShutdownEvent, StatusNotProvided, resp.Include(ShutdownEvent, StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", ShutdownEvent, StatusNotProvided, resp.Exclude(ShutdownEvent, StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", StartupEvent, StatusNotProvided, resp.Include(StartupEvent, StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", StartupEvent, StatusNotProvided, resp.Exclude(StartupEvent, StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", PingEvent, StatusNotProvided, resp.Include(PingEvent, StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", PingEvent, StatusNotProvided, resp.Exclude(PingEvent, StatusNotProvided))

	//Output:
	//test: count() -> : 5
	//test: Get(invalid) -> : [ok:false] [msg-nil:true]
	//test: Get(from-uri-3) -> : [ok:true] [msg-to:to-uri]
	//test: include(event:shutdown,95) -> : []
	//test: exclude(event:shutdown,95) -> : [from-uri-0 from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:startup,95) -> : [from-uri-0]
	//test: exclude(event:startup,95) -> : [from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:ping,95) -> : [from-uri-2 from-uri-3]
	//test: exclude(event:ping,95) -> : [from-uri-0 from-uri-1 from-uri-4]

}
