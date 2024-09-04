package http2

import (
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"sync"
)

type MatchFunc[T any] func(r *http.Request, item *T) bool
type PatchProcessFunc[T any, U any] func(r *http.Request, list *[]T, content *U) *core.Status
type PostProcessFunc[T any, V any] func(r *http.Request, list *[]T, content *V) *core.Status

type Content2[T any, U any, V any] interface {
	//Count2(string) int
	Count() int
	Empty()
	Get(r *http.Request) ([]T, *core.Status)
	Put(r *http.Request, items []T) *core.Status
	Delete(r *http.Request) *core.Status
	Patch(r *http.Request, post *U) *core.Status
	Post(r *http.Request, post *V) *core.Status
}

type ListContent[T any, U any, V any] struct {
	List  []T
	mu    sync.Mutex
	mutex bool
	match MatchFunc[T]
	patch PatchProcessFunc[T, U]
	post  PostProcessFunc[T, V]
}

func NewListContent[T any, U any, V any](mutex bool, match MatchFunc[T], patch PatchProcessFunc[T, U], post PostProcessFunc[T, V]) Content2[T, U, V] {
	c := new(ListContent[T, U, V])
	c.mutex = mutex
	c.match = match
	c.patch = patch
	c.post = post
	if c.match == nil {
		if c.match == nil {
			c.match = func(r *http.Request, item *T) bool { return false }
		}
	}
	return c
}

func (c *ListContent[T, U, V]) acquire() func() {
	if c.mutex {
		c.mu.Lock()
		return func() {
			c.mu.Unlock()
		}
	} else {
		return func() {}
	}
}

/*
func (c *ListContent[T, U, V]) count2(goid string) int {
	fmt.Printf("Count-start() -> %v - %v\n", goid, time.Now().UTC())
	defer c.acquire()()
	time.Sleep(time.Second * 2)
	fmt.Printf("Count-stop()  -> %v - %v\n", goid, time.Now().UTC())
	return len(c.List)
}

*/

func (c *ListContent[T, U, V]) Count() int {
	defer c.acquire()()
	return len(c.List)
}

func (c *ListContent[T, U, V]) Empty() {
	defer c.acquire()()
	c.List = nil
}

func (c *ListContent[T, U, V]) Get(r *http.Request) ([]T, *core.Status) {
	if r == nil {
		return nil, core.NewStatus(core.StatusInvalidArgument)
	}
	defer c.acquire()()
	var items []T
	for _, target := range c.List {
		if c.match(r, &target) {
			items = append(items, target)
		}
	}
	if len(items) == 0 {
		return nil, core.StatusNotFound()
	}
	return items, core.StatusOK()
}

func (c *ListContent[T, U, V]) Put(r *http.Request, items []T) *core.Status {
	if r == nil {
		return core.NewStatus(core.StatusInvalidArgument)
	}
	defer c.acquire()()
	if len(items) != 0 {
		c.List = append(c.List, items...)
	}
	return core.StatusOK()
}

func (c *ListContent[T, U, V]) Delete(r *http.Request) *core.Status {
	if r == nil {
		return core.NewStatus(core.StatusInvalidArgument)
	}
	defer c.acquire()()
	count := 0
	deleted := true
	for deleted {
		deleted = false
		for i, target := range c.List {
			if c.match(r, &target) {
				c.List = append(c.List[:i], c.List[i+1:]...)
				deleted = true
				count++
				break
			}
		}
	}
	if count == 0 {
		return core.StatusNotFound()
	}
	return core.StatusOK()
}

func (c *ListContent[T, U, V]) Patch(r *http.Request, patch *U) *core.Status {
	if r == nil || patch == nil {
		return core.NewStatus(core.StatusInvalidArgument)
	}
	if c.patch == nil {
		return core.NewStatus(http.StatusBadRequest)
	}
	defer c.acquire()()
	return c.patch(r, &c.List, patch)
}

func (c *ListContent[T, U, V]) Post(r *http.Request, post *V) *core.Status {
	if r == nil || post == nil {
		return core.NewStatus(core.StatusInvalidArgument)
	}
	if c.post == nil {
		return core.NewStatus(http.StatusBadRequest)
	}
	defer c.acquire()()
	return c.post(r, &c.List, post)
}
