package http2

import (
	"fmt"
	"sync"
	"time"
)

type Lock4 struct{ m sync.Mutex }

type access struct {
	mu sync.Mutex
	//v  W
	mutex bool
}

func newAccess[W any]() *access {
	var w W
	a := new(access)
	if _, ok := any(w).(Lock4); ok {
		a.mutex = true
	}
	return a

}

func (a *access) get(goid string) {
	fmt.Printf("get-start() -> %v - %v\n", goid, time.Now().UTC())
	if a.mutex {
		a.mu.Lock()
		defer func() {
			a.mu.Unlock()
		}()

	}

	//}
	//defer fn()
	//a.mu.mu.Lock()

	time.Sleep(time.Second * 2)
	fmt.Printf("get-stop()  -> %v - %v\n", goid, time.Now().UTC())
	//a.mu.mu.Unlock()
	//fn()
}

func ExampleMutex_Lock() {
	//var t Lock
	a := newAccess[Lock4]()

	fmt.Printf("example-start()  ->     %v\n", time.Now().UTC())
	go func() {
		time.Sleep(time.Millisecond * 500)
		a.get("goid-2")
	}()
	a.get("goid-1")
	time.Sleep(time.Second * 3)
	fmt.Printf("example-stop()   ->     %v\n", time.Now().UTC())

	//Output:
	//fail

}
