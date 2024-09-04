package controller

import (
	"errors"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

type Controller struct {
	RouteName string
	Router    *Router
}

func NewController(routeName string, primary, secondary *Resource) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Router = NewRouter(primary, secondary)
	return c
}

func NewExchangeController(routeName string, ex core.HttpExchange) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	authority := core.Authority(ex)
	c.Router = NewRouter(NewPrimaryResource("", authority, 0, "", ex), nil)
	return c
}

//func (c *Controller) Name() string {
//	return c.RouteName
//}

func (c *Controller) Do(do core.HttpExchange, req *http.Request) (resp *http.Response, status *core.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusBadRequest}, core.NewStatusError(core.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	traffic := access.EgressTraffic
	from := req.Header.Get(core.XFrom)
	rsc := c.Router.RouteTo()
	if rsc.Handler != nil {
		traffic = access.InternalTraffic
		do = rsc.Handler
	} else {
		if do == nil {
			return &http.Response{StatusCode: http.StatusBadRequest}, core.NewStatusError(core.StatusInvalidArgument, errors.New("invalid argument : core.HttpExchange is nil"))
		}
	}
	inDuration, outDuration := durations(rsc, req)
	duration := time.Duration(0)
	controllerCode := ""
	newURL := rsc.BuildURL(req.URL)
	req.URL = newURL
	if req.URL != nil {
		req.Host = req.URL.Host
	}
	start := time.Now().UTC()

	// if no timeout or an existing deadline and existing deadline is <= timeout, then use the existing request
	if outDuration == 0 || (inDuration > 0 && inDuration <= outDuration) {
		duration = inDuration * -1
		resp, status = do(req)
	} else {
		duration = outDuration
		// Internal call
		if rsc.Handler != nil {
			//ctx, cancel := context.WithTimeout(req.Context(), outDuration)
			//defer cancel()
			//r2 := req.Clone(ctx)
			resp, status = doInternal(outDuration, do, req)
		} else {
			resp, status = doEgress(outDuration, do, req)
		}
	}
	elapsed := time.Since(start)
	if resp != nil {
		c.Router.UpdateStats(resp.StatusCode, rsc)
		if resp.StatusCode == http.StatusGatewayTimeout {
			controllerCode = access.ControllerTimeout
		}
	} else {
		resp = &http.Response{StatusCode: status.HttpCode()}
	}
	if !disableLogging {
		access.Log(traffic, start, elapsed, req, resp, access.Routing{From: from, Route: c.RouteName, To: rsc.Name, Percent: -1}, access.Controller{Timeout: duration, RateLimit: -1, RateBurst: 0, Code: controllerCode})
	}
	return
}

func durations(rsc *Resource, req *http.Request) (in time.Duration, out time.Duration) {
	deadline, ok := req.Context().Deadline()
	if ok {
		in = time.Until(deadline) // * -1
	}
	if rsc.Duration > 0 {
		out = rsc.Duration
	}
	return
}
