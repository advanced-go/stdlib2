package http2

import "sync"

type Mutex interface {
	Acquire() func()
}

type Lock struct {
	mu sync.Mutex
}

func (l Lock) Acquire() func() {
	l.mu.Lock()
	return func() { l.mu.Unlock() }
}

type Lock2 struct {
	mu sync.Mutex
}

func (l Lock2) Acquire() func() {
	l.mu.Lock()
	return func() { l.mu.Unlock() }
}

type Lock3 sync.Mutex

func (l Lock3) Lock() {
	//l.Lock()
}

//func (l Lock2) Acquire() func() {
//	l.mu.Lock()
//	return func() { l.mu.Unlock() }
//}

type NoLock struct{}

func (l NoLock) Acquire() func() {
	return func() {}
}
