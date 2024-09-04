package controller2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"time"
)

type Config struct {
	RouteName string `json:"route"`
	Host      string `json:"host"`
	Authority string `json:"authority"`
	Duration  time.Duration
}

func New(cfg *Config, handler core.HttpExchange) *Controller {
	var prime *Resource
	var second *Resource
	if handler == nil {
		prime = NewPrimaryResource(cfg.Host, cfg.Authority, cfg.Duration, nil)
	} else {
		prime = NewPrimaryResource("", cfg.Authority, cfg.Duration, handler)
		second = NewSecondaryResource(cfg.Host, cfg.Authority, cfg.Duration, nil)
	}
	return NewController(cfg.RouteName, prime, second)
}

type Controller struct {
	RouteName string
	Primary   *Resource
	Secondary *Resource
}

func NewController(routeName string, primary, secondary *Resource) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Primary = primary
	c.Secondary = secondary
	return c
}

func RegisterControllerFromConfig(config *Config, ex core.HttpExchange) *core.Status {
	ctrl := New(config, ex)
	err := RegisterController(ctrl)
	if err != nil {
		return core.NewStatusError(core.StatusInvalidArgument, err)
	}
	return core.StatusOK()
}

// RegisterController - add a controller for an egress route
func RegisterController(ctrl *Controller) error {
	if ctrl == nil {
		return errors.New(fmt.Sprintf("invalid argument: Controller is nil"))
	}
	//if ctrl.Router == nil {
	//	return errors.New(fmt.Sprintf("invalid argument: Controller router is nil"))
	//}
	if ctrl.Primary == nil {
		return errors.New(fmt.Sprintf("invalid argument: Controller rimary resource is nil"))
	}
	if len(ctrl.Primary.Authority) == 0 {
		if ctrl.Primary.Host == "" {
			return errors.New(fmt.Sprintf("invalid argument: Controller primary resource host is empty"))
		}
		return ctrlMap.register(ctrl)
	}
	return ctrlMap.registerWithAuthority(ctrl)
}
