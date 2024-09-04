package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

const (
	maxWait = timeout + time.Millisecond*100
)

//var pingStart = time.Now()

func ExamplePing_Good() {
	uri1 := "urn:ping:good"
	ex := NewExchange()

	c := make(chan *Message, 16)
	a, _ := NewAgentWithChannels(uri1, c, nil, nil, emptyRun, nil)
	ex.Register(a)
	go pingGood(c)
	status := Ping(ex, uri1)
	fmt.Printf("test: Ping(good) -> [%v] [timeout:%v] [duration<3:%v]\n", status, timeout, status.Duration < time.Second*3)

	//Output:
	//test: Ping(good) -> [OK] [timeout:3s] [duration<3:true]

}

func ExamplePing_Timeout() {
	uri2 := "urn:ping:timeout"
	c := make(chan *Message, 16)

	ex := NewExchange()
	a, _ := NewAgentWithChannels(uri2, c, nil, nil, emptyRun, nil)
	ex.Register(a)
	go pingTimeout(c)
	status := Ping(ex, uri2)
	fmt.Printf("test: Ping(timeout) -> [%v] [timeout:%v] [duration>3:%v]\n", status, timeout, status.Duration > time.Second*3)

	//Output:
	//test: Ping(timeout) -> [Timeout] [timeout:3s] [duration>3:true]

}

func ExamplePing_Error() {
	uri3 := "urn:ping:error"
	ex := NewExchange()

	c := make(chan *Message, 16)
	a, _ := NewAgentWithChannels(uri3, c, nil, nil, emptyRun, nil)
	ex.Register(a)
	go pingError(c, errors.New("ping response error"))
	status := Ping(ex, uri3)
	fmt.Printf("test: Ping(error) -> [status:%v] [timeout:%v] [duration<3:%v]\n", status, timeout, status.Duration < time.Second*3)

	//Output:
	//recovered in messaging.NewReceiverReplyTo() : send on closed channel
	//test: Ping(error) -> [status:I'm A Teapot [ping response error]] [timeout:3s] [duration<3:true]

}

func ExamplePing_Delay() {
	uri4 := "urn:ping:delay"
	ex := NewExchange()

	c := make(chan *Message, 16)
	a, _ := NewAgentWithChannels(uri4, c, nil, nil, emptyRun, nil)
	ex.Register(a)
	go pingDelay(c)
	status := Ping(ex, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [timeout:%v] [duration>timeout/2:%v]\n", status, timeout, status.Duration > timeout/2)

	//Output:
	//test: Ping(delay) -> [OK] [timeout:3s] [duration>timeout/2:true]

}

func pingGood(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			SendReply(msg, core.StatusOK())
		default:
		}
	}
}

func pingTimeout(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(maxWait)
			SendReply(msg, core.StatusOK())
		default:
		}
	}
}

func pingError(c chan *Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				SendReply(msg, core.NewStatusError(http.StatusTeapot, errors.New("ping response error")))
			}
		default:
		}
	}
}

func pingDelay(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(timeout / 2)
			SendReply(msg, core.StatusOK())
		default:
		}
	}
}
