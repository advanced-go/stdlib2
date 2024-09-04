package messaging

import (
	"fmt"
	"time"
)

func printAgentRun(uri string, ctrl, data <-chan *Message, state any) {
	fmt.Printf("test: AgentRun() -> [uri:%v] [ctrl:%v] [data:%v] [state:%v]\n", uri, ctrl != nil, data != nil, state != nil)
}

func testAgentRun(uri string, ctrl, data <-chan *Message, _ any) {
	for {
		select {
		case msg, open := <-ctrl:
			if !open {
				return
			}
			fmt.Printf("test: AgentRun() -> %v\n", msg)
			if msg.Event() == ShutdownEvent {
				return
			}
		default:
		}
		select {
		case msg, open := <-data:
			if !open {
				return
			}
			fmt.Printf("test: AgentRun() -> %v\n", msg)
		default:
		}
	}
}

func ExampleNewAgent_Error() {
	a, err := newAgent("", nil, nil, nil, nil, nil)
	fmt.Printf("test: newAgent() -> [agent:%v] [%v]\n", a, err)

	a, err = newAgent("urn:agent7", nil, nil, nil, nil, nil)
	fmt.Printf("test: newAgent() -> [agent:%v] [%v]\n", a, err)

	//Output:
	//test: newAgent() -> [agent:<nil>] [error: agent URI is empty]
	//test: newAgent() -> [agent:<nil>] [error: agent AgentFunc is nil]

}

func ExampleNewAgent() {
	uri := "urn:agent007"
	uri1 := "urn:agent009"

	a, _ := NewAgent(uri, printAgentRun, nil)
	a.Run()
	time.Sleep(time.Second)

	a, _ = NewAgentWithChannels(uri1, nil, nil, nil, printAgentRun, "data")
	a.Run()
	time.Sleep(time.Second)

	//Output:
	//test: AgentRun() -> [uri:urn:agent007] [ctrl:true] [data:true] [state:false]
	//test: AgentRun() -> [uri:urn:agent009] [ctrl:true] [data:false] [state:true]

}

func ExampleOnShutdown() {
	uri := "urn:agent007"

	a, _ := NewAgent(uri, printAgentRun, nil)
	if a1, ok := any(a).(*agent); ok {
		a1.running = true
	}
	a.Shutdown()

	a, _ = NewAgent(uri, printAgentRun, nil)
	if sd, ok := a.(OnShutdown); ok {
		sd.Add(func() { fmt.Printf("test: OnShutdown() -> func-1()\n") })
		sd.Add(func() { fmt.Printf("test: OnShutdown() -> func-2()\n") })
		sd.Add(func() { fmt.Printf("test: OnShutdown() -> func-3()\n") })
	}
	if a1, ok := any(a).(*agent); ok {
		a1.running = true
	}
	a.Shutdown()

	//Output:
	//test: OnShutdown() -> func-1()
	//test: OnShutdown() -> func-2()
	//test: OnShutdown() -> func-3()

}

func ExampleAgentRun() {
	uri := "urn:agent007"
	a, _ := NewAgent(uri, testAgentRun, nil)
	a.Run()
	a.Message(NewControlMessage(uri, "ExampleAgentRun()", StartupEvent))
	a.Message(NewDataMessage(uri, "ExampleAgentRun()", DataEvent))
	time.Sleep(time.Second)
	a.Shutdown()
	time.Sleep(time.Second)

	//Output:
	//test: AgentRun() -> [chan:CTRL] [from:ExampleAgentRun()] [to:urn:agent007] [event:startup]
	//test: AgentRun() -> [chan:DATA] [from:ExampleAgentRun()] [to:urn:agent007] [event:data]
	//test: AgentRun() -> [chan:CTRL] [from:urn:agent007] [to:urn:agent007] [event:shutdown]

}
