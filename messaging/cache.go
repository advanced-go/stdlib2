package messaging

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Cache - message cache by uri
type Cache struct {
	m    map[string]*Message
	errs []error
	mu   sync.RWMutex
}

// NewCache - create a message cache
func NewCache() *Cache {
	c := new(Cache)
	c.m = make(map[string]*Message)
	return c
}

// Count - return the count of items
func (r *Cache) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, _ = range r.m {
		count++
	}
	return count
}

// Filter - apply a filter against a traversal of all items
func (r *Cache) Filter(event string, code int, include bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var uri []string
	for u, resp := range r.m {
		s := resp.Status()
		if s == nil {
			fmt.Printf("no status available : %v\n", u)
			continue
		}
		if include {
			if s.Code == code && resp.Event() == event {
				uri = append(uri, u)
			}
		} else {
			if s.Code != code || resp.Event() != event {
				uri = append(uri, u)
			}
		}
	}
	sort.Strings(uri)
	return uri
}

// Include - filter for items that include a specific event
func (r *Cache) Include(event string, status int) []string {
	return r.Filter(event, status, true)
}

// Exclude - filter for items that do not include a specific event
func (r *Cache) Exclude(event string, status int) []string {
	return r.Filter(event, status, false)
}

// Add - add a message
func (r *Cache) Add(msg *Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if msg.From() == "" {
		err := errors.New("invalid argument: message from is empty")
		r.errs = append(r.errs, err)
		return err
	}
	if _, ok := r.m[msg.From()]; !ok {
		r.m[msg.From()] = msg
		return nil
	}
	err0 := errors.New(fmt.Sprintf("invalid argument: message found [%v]", msg.From()))
	r.errs = append(r.errs, err0)
	return err0
}

// Get - get a message based on a URI
func (r *Cache) Get(uri string) (*Message, bool) {
	if uri == "" {
		return nil, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[uri]; ok {
		return r.m[uri], true
	}
	return nil, false //errors.New(fmt.Sprintf("invalid argument: uri not found [%v]", uri))
}

// Uri - list the URI's in the cache
func (r *Cache) Uri() []string {
	var uri []string
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, _ := range r.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}

// ErrorList - list of errors
func (r *Cache) ErrorList() []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.errs
}

// NewCacheHandler - handler to receive messages into a cache.
func NewCacheHandler(cache *Cache) Handler {
	return func(msg *Message) {
		err := cache.Add(msg)
		if err != nil {
			fmt.Printf("error: messaging cache handler cache.Add() %v\n", err)
		}
	}
}
