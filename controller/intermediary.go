package controller

import (
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

func NewControllerIntermediary(routeName string, c2 core.HttpExchange) core.HttpExchange {
	return func(r *http.Request) (resp *http.Response, status *core.Status) {
		if c2 == nil {
			return &http.Response{StatusCode: http.StatusBadRequest}, core.StatusBadRequest() //errors.New("error: AccessLog Intermediary HttpExchange is nil")
		}
		controllerCode := ""
		from := r.Header.Get(core.XFrom)

		var dur time.Duration
		if ct, ok := r.Context().Deadline(); ok {
			dur = time.Until(ct) * -1
		}
		start := time.Now().UTC()
		resp, status = c2(r)
		if status.Code == http.StatusGatewayTimeout {
			controllerCode = access.ControllerTimeout
		}
		access.Log(access.InternalTraffic, start, time.Since(start), r, resp, access.Routing{From: from, Route: routeName, Percent: -1}, access.Controller{Timeout: dur, Code: controllerCode})
		return
	}
}
