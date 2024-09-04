package controller

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	uri2 "github.com/advanced-go/stdlib/uri"
	"net/http"
	"sync"
)

var (
	ctrlMap        = NewControls()
	disableLogging = false
)

func DisableLogging(v bool) func() {
	disableLogging = v
	return func() {
		disableLogging = !v
	}
}

func updateHost(host, authority string, primary bool) (status *core.Status) {
	if host == "" || authority == "" {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New("invalid argument: host or authority is empty"))
	}
	var ctrl *Controller
	ctrl, status = LookupWithAuthority(authority)
	if !status.OK() {
		return
	}
	if primary {
		if ctrl.Router.Primary == nil {
			return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: primary resource is nil for authority: %v", authority)))
		}
		ctrl.Router.Primary.Host = host
	} else {
		if ctrl.Router.Secondary == nil {
			return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: secondary resource is nil for authority: %v", authority)))
		}
		ctrl.Router.Primary.Host = host
	}
	return core.StatusOK()
}

func UpdatePrimaryHost(host, authority string) (status *core.Status) {
	return updateHost(host, authority, true)
}

func UpdateSecondaryHost(host, authority string) (status *core.Status) {
	return updateHost(host, authority, false)
}

func updateExchange(list []core.HttpExchange, primary bool) (status *core.Status) {
	if list == nil {
		return core.NewStatus(core.StatusInvalidArgument)
	}

	var ctrl *Controller
	var authority = ""
	var rsc *Resource
	for _, ex := range list {
		authority = core.Authority(ex)
		ctrl, status = LookupWithAuthority(authority)
		if !status.OK() {
			continue
		}
		if primary {
			rsc = ctrl.Router.Primary
			if rsc == nil {
				return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: primary resource is nil for authority: %v", authority)))
			}
		} else {
			rsc = ctrl.Router.Secondary
			if rsc == nil {
				return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: secondary resource is nil for authority: %v", authority)))
			}
		}
		if rsc.Handler == nil {
			rsc.Handler = ex
		}
	}
	return status
}

func UpdatePrimaryExchange(list []core.HttpExchange) (status *core.Status) {
	return updateExchange(list, true)
}

func UpdateSecondaryExchange(list []core.HttpExchange) (status *core.Status) {
	return updateExchange(list, true)
}

func LookupWithConfig(cfg Config) (ctrl *Controller, status *core.Status) {
	if cfg.Authority != "" {
		return ctrlMap.lookup(cfg.Authority)
	}
	return ctrlMap.lookup(cfg.Host)
}

func LookupWithAuthority(authority string) (ctrl *Controller, status *core.Status) {
	return ctrlMap.lookup(authority)
}

func Lookup(req *http.Request) (ctrl *Controller, status *core.Status) {
	if req == nil || req.URL == nil {
		return nil, core.NewStatus(http.StatusNotFound)
	}

	// Try host first
	ctrl, status = ctrlMap.lookup(req.Host)
	if status.OK() {
		return
	}

	// Default to embedded authority
	p := uri2.Uproot(req.URL.Path)
	if p.Valid {
		ctrl, status = ctrlMap.lookup(p.Authority)
		if status.OK() {
			return
		}
	}
	return nil, core.NewStatus(http.StatusNotFound)
}

// controls - key value pairs of an authority -> *Controller
type controls struct {
	m *sync.Map
}

// NewControls - create a new Controls map
func NewControls() *controls {
	p := new(controls)
	p.m = new(sync.Map)
	return p
}

func (p *controls) register(ctrl *Controller) error {
	if ctrl == nil {
		return errors.New("invalid argument: Controller is nil")
	}
	_, ok1 := p.m.Load(ctrl.Router.Primary.Host)
	if ok1 {
		return errors.New(fmt.Sprintf("invalid argument: Controller already exists for authority: [%v]", ctrl.Router.Primary))
	}
	p.m.Store(ctrl.Router.Primary.Host, ctrl)
	return nil
}

func (p *controls) registerWithAuthority(ctrl *Controller) error {
	if ctrl == nil {
		return errors.New("invalid argument: Controller is nil")
	}
	//parsed := uri2.Uproot(ctrl.Router.Primary.Authority)
	//if !parsed.Valid {
	//	return errors.New(fmt.Sprintf("invalid argument: Controller primary authority is invalid: [%v] [%v]", ctrl.Router.Primary.Authority, parsed.Err))
	//}
	_, ok1 := p.m.Load(ctrl.Router.Primary.Authority)
	if ok1 {
		return errors.New(fmt.Sprintf("invalid argument: Controller already exists for authority : [%v]", ctrl.Router.Primary.Authority))
	}
	p.m.Store(ctrl.Router.Primary.Authority, ctrl)
	return nil
}

// Lookup - get a Controller using an authority
func (p *controls) lookup(authority string) (*Controller, *core.Status) {
	if authority == "" {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("invalid argument: authority is empty"))
	}
	v, ok := p.m.Load(authority)
	if !ok {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: Controller does not exist: [%v]", authority)))
	}
	if ctrl, ok1 := v.(*Controller); ok1 {
		return ctrl, core.StatusOK()
	}
	return nil, core.NewStatus(core.StatusInvalidContent)
}

func (p *controls) remove(key string) {
	p.m.Delete(key)
}
