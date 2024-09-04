package controller

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"time"
)

type Config struct {
	RouteName    string `json:"route"`
	Host         string `json:"host"`
	Authority    string `json:"authority"`
	LivenessPath string `json:"liveness"`
	Duration     time.Duration
}

func New(cfg *Config, handler core.HttpExchange) *Controller {
	var prime *Resource
	var second *Resource
	if handler == nil {
		prime = NewPrimaryResource(cfg.Host, cfg.Authority, cfg.Duration, cfg.LivenessPath, nil)
	} else {
		prime = NewPrimaryResource("", cfg.Authority, cfg.Duration, cfg.LivenessPath, handler)
		second = NewSecondaryResource(cfg.Host, cfg.Authority, cfg.Duration, cfg.LivenessPath, nil)
	}
	return NewController(cfg.RouteName, prime, second)
}

func GetRoute(name string, config []Config) (Config, bool) {
	for _, cfg := range config {
		if cfg.RouteName == name {
			return cfg, true
		}
	}
	return Config{}, false
}

func RegisterControllerFromConfig(config *Config, ex core.HttpExchange) *core.Status {
	ctrl := New(config, ex)
	err := RegisterController(ctrl)
	if err != nil {
		return core.NewStatusError(core.StatusInvalidArgument, err)
	}
	return core.StatusOK()
}

func RegisterControllerFromRoutes(routeName string, config []Config, ex core.HttpExchange) *core.Status {
	cfg, ok := GetRoute(routeName, config)
	if !ok {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("error: route name not found: %v\n", routeName)))
	}
	return RegisterControllerFromConfig(&cfg, ex)
}
