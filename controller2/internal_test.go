package controller2

import (
	"context"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"io"
	"net/http"
	"time"
)

func httpCall(r *http.Request) (resp *http.Response, status *core.Status) {
	cnt := 0
	var err0 error
	var err1 error

	if r.URL.Path == core.AuthorityRootPath {
		h := make(http.Header)
		h.Add(core.XAuthority, "authority")
		return &http.Response{StatusCode: http.StatusOK, Header: h}, core.StatusOK()
	}
	resp, err0 = http.DefaultClient.Do(r)
	if err0 != nil {
		resp = new(http.Response)
		if r.Context().Err() == context.DeadlineExceeded {
			status = core.NewStatus(http.StatusGatewayTimeout)
		} else {
			status = core.NewStatus(http.StatusInternalServerError)
		}
		resp.StatusCode = status.Code
	} else {
		var buf []byte
		buf, err1 = io.ReadAll(resp.Body)
		if err1 != nil {
			if err1 == context.DeadlineExceeded {
				status = core.NewStatus(http.StatusGatewayTimeout)
			} else {
				status = core.NewStatus(http.StatusInternalServerError)
			}
		} else {
			resp.ContentLength = int64(len(buf))
			cnt = len(buf)
			status = core.StatusOK()
		}
	}
	fmt.Printf("test: httpCall() -> [content:%v] [do-err:%v] [read-err:%v]\n", cnt > 0, err0, err1)
	return
}

func ExampleDoInternal() {
	var buf []byte
	uri := "https://www.google.com/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := doInternal(0, httpCall, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: DoInternal_0s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	resp, status = doInternal(time.Second*5, httpCall, req)
	buf = nil
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: DoInternal_5s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	resp, status = doInternal(time.Millisecond*5, httpCall, req)
	buf = nil
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: DoInternal_5ms() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: DoInternal_0s() -> [status-code:400] [status:Bad Request] [buf:false]
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>]
	//test: DoInternal_5s() -> [status-code:200] [status:OK] [buf:false]
	//test: httpCall() -> [content:false] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>]
	//test: DoInternal_5ms() -> [status-code:504] [status:Timeout] [buf:false]

}
