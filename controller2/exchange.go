package controller2

import (
	"errors"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

func Exchange(req *http.Request, do core.HttpExchange, ctrl *Controller) (resp *http.Response, status *core.Status) {
	if req == nil || do == nil || ctrl == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, core.NewStatusError(core.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	//var ctrl *Controller
	//ctrl, status = lookup(req)
	//if !status.OK() {
	//	return do(req)
	//}
	localDo := do
	traffic := access.EgressTraffic
	rsc := ctrl.Primary
	if rsc.Handler != nil {
		localDo = rsc.Handler
		traffic = access.InternalTraffic
	}
	inDuration, outDuration := durations(rsc, req)
	duration := time.Duration(0)
	reasonCode := ""
	newURL := rsc.BuildURL(req.URL)
	req.URL = newURL
	if req.URL != nil {
		req.Host = req.URL.Host
	}
	start := time.Now().UTC()
	from := req.Header.Get(core.XFrom)

	// if no timeout or an existing deadline and existing deadline is <= timeout, then use the existing request
	if outDuration == 0 || (inDuration > 0 && inDuration <= outDuration) {
		duration = inDuration * -1
		resp, status = localDo(req)
	} else {
		duration = outDuration
		if rsc.Handler != nil {
			resp, status = doInternal(outDuration, localDo, req)
		} else {
			resp, status = doEgress(outDuration, localDo, req)
		}
	}
	if resp != nil {
		if resp.StatusCode == http.StatusGatewayTimeout {
			reasonCode = access.ControllerTimeout
		}
	} else {
		resp = &http.Response{StatusCode: status.HttpCode()}
	}
	access.Log(traffic, start, time.Since(start), req, resp, access.Routing{From: from, Route: ctrl.RouteName, To: rsc.Name, Percent: -1}, access.Controller{Timeout: duration, Code: reasonCode})
	return
}

func durations(rsc *Resource, req *http.Request) (in time.Duration, out time.Duration) {
	if req != nil && req.Context() != nil {
		deadline, ok := req.Context().Deadline()
		if ok {
			in = time.Until(deadline) // * -1
		}
	}
	if rsc.Duration > 0 {
		out = rsc.Duration
	}
	return
}
