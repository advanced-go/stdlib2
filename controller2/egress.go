package controller2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"sync/atomic"
	"time"
)

func doEgress(duration time.Duration, do core.HttpExchange, req *http.Request) (resp *http.Response, status *core.Status) {
	var closed int64
	ch := make(chan struct{}, 1)
	tick := time.Tick(duration)
	go func() {
		defer func() {
			// Check for when a timeout is reached, the channel is closed, and there is a pending send
			if r := recover(); r != nil {
				fmt.Printf("panic: recovered in controller.doEgress() : %v\n", r)
			}
		}()
		resp, status = do(req)
		if atomic.LoadInt64(&closed) == 0 {
			ch <- struct{}{}
		}
	}()
	select {
	case <-tick:
		resp = &http.Response{StatusCode: http.StatusGatewayTimeout}
		status = core.NewStatusError(core.StatusDeadlineExceeded, errors.New("context deadline exceeded"))
	case <-ch:
	}
	atomic.StoreInt64(&closed, 1)
	close(ch)
	return
}
