package messagingtest

import (
	"fmt"
	"github.com/advanced-go/stdlib/messaging"
)

func ExampleNewAgent() {
	a := NewAgent()
	if _, ok := any(a).(messaging.OpsAgent); ok {
		fmt.Printf("test: OpsAgent() -> ok\n")
	} else {
		fmt.Printf("test: OpsAgent() -> fail\n")
	}

	//Output:
	//test: OpsAgent() -> ok

}

func ExampleOld() {
	fmt.Printf("test")

	//Output:
	//test

}
