package controller

import (
	"context"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

func doInternal(duration time.Duration, handler core.HttpExchange, req *http.Request) (resp *http.Response, status *core.Status) {
	if duration <= 0 || handler == nil {
		return &http.Response{StatusCode: http.StatusBadRequest}, core.NewStatus(http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(req.Context(), duration)
	defer cancel()
	r2 := req.Clone(ctx)
	resp, status = handler(r2)
	return
}
