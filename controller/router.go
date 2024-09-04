package controller

import (
	"sync/atomic"
)

const (
	PrimaryName   = "primary"
	SecondaryName = "secondary"
	primary       = 0
	secondary     = 1
)

type Router struct {
	Primary   *Resource
	Secondary *Resource
	active    atomic.Int64
}

func NewRouter(primary, secondary *Resource) *Router {
	r := new(Router)
	r.Primary = primary
	if primary != nil {
		r.Primary.Name = PrimaryName
	}
	r.Secondary = secondary
	if r.Secondary != nil {
		r.Secondary.Name = SecondaryName
	}
	return r
}

func (r *Router) RouteTo() *Resource {
	if r.active.Load() == primary {
		return r.Primary
	}
	return r.Secondary
}

func (r *Router) UpdateStats(statusCode int, rsc *Resource) {

}

func (r *Router) Swap() (swapped bool) {
	old := r.active.Load()
	if old == primary {
		swapped = r.active.CompareAndSwap(old, secondary)
	} else {
		swapped = r.active.CompareAndSwap(old, primary)
	}
	return
}
