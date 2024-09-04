package controller

import (
	"errors"
	"fmt"
)

// RegisterController - add a controller for an egress route
func RegisterController(ctrl *Controller) error {
	if ctrl == nil {
		return errors.New(fmt.Sprintf("invalid argument: Controller is nil"))
	}
	if ctrl.Router == nil {
		return errors.New(fmt.Sprintf("invalid argument: Controller router is nil"))
	}
	if ctrl.Router.Primary == nil {
		return errors.New(fmt.Sprintf("invalid argument: Controller router primary resource is nil"))
	}
	if len(ctrl.Router.Primary.Authority) == 0 {
		if ctrl.Router.Primary.Host == "" {
			return errors.New(fmt.Sprintf("invalid argument: Controller router primary resource host is empty"))
		}
		return ctrlMap.register(ctrl)
	}
	return ctrlMap.registerWithAuthority(ctrl)
}

func RegisterControllerList(ctrl []*Controller) error {
	var err error

	for _, c := range ctrl {
		err = RegisterController(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterControllerWithDefer(ctrl *Controller, def func()) func() {
	err := RegisterController(ctrl)
	// !panic
	fn := def
	return func() {
		if fn != nil {
			fn()
		}
		if err != nil {
			fmt.Printf("register controller error : %v\n", err)
			return
		}
		UnregisterController(ctrl)
	}
}

func RegisterControllerListWithDefer(ctrl []*Controller) func() {
	var fn func()

	for _, c := range ctrl {
		fn = RegisterControllerWithDefer(c, fn)
	}
	return fn
}

func UnregisterController(ctrl *Controller) {
	if ctrl != nil {
		key := ctrl.Router.Primary.Authority
		if key == "" {
			key = ctrl.Router.Primary.Host
		}
		ctrlMap.remove(key)
	}
}

func UnregisterControllerList(ctrl []*Controller) {
	for _, c := range ctrl {
		UnregisterController(c)
	}
}
