package messaging

import (
	"fmt"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"time"
)

func _ExampleTicker() {
	t := NewTicker(time.Second * 2)
	ctrl := make(chan *Message)

	go tickerRun(ctrl, t)
	time.Sleep(time.Second * 20)

	ctrl <- NewControlMessage("to", "from", ShutdownEvent)
	time.Sleep(time.Second * 2)

	//Output:
	//test: Ticker() -> 2024-07-11T14:39:57.164Z
	//test: Ticker() -> 2024-07-11T14:39:59.164Z
	//test: Ticker() -> 2024-07-11T14:40:04.182Z
	//test: Ticker() -> 2024-07-11T14:40:09.180Z
	//test: Ticker() -> 2024-07-11T14:40:11.193Z
	//test: Ticker() -> 2024-07-11T14:40:13.184Z

}

func tickerRun(ctrl <-chan *Message, t *Ticker) {
	count := 0
	t.Start(0)
	for {
		select {
		case <-t.C():
			fmt.Printf("test: Ticker() -> %v\n", fmt2.FmtRFC3339Millis(time.Now().UTC()))
			count++
			if count == 2 {
				t.Start(time.Second * 5)
			}
			if count == 4 {
				t.Reset()
			}
		case msg := <-ctrl:
			switch msg.Event() {
			case ShutdownEvent:
				return
			default:
			}
		default:
		}
	}
}
