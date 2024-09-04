package controller

import (
	"bytes"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	uri2 "github.com/advanced-go/stdlib/uri"
	"io"
	"net/http"
	"time"
)

/*
func testDo(req *http.Request) (*http.Response, *core.Status) {
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		return resp, core.NewStatus(resp.StatusCode)
	}
	resp = &http.Response{StatusCode: core.StatusDeadlineExceeded}
	return resp, core.NewStatusError(core.StatusDeadlineExceeded, err)
}
*/

func testDo(r *http.Request) (*http.Response, *core.Status) {
	req, _ := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if resp == nil {
			resp = &http.Response{StatusCode: http.StatusGatewayTimeout, Body: io.NopCloser(bytes.NewReader([]byte("Timeout [Get \"https://www.google.com/search?q=golang\": context deadline exceeded]")))}
			return resp, core.NewStatus(core.StatusDeadlineExceeded)
		}
		resp.Body = io.NopCloser(bytes.NewReader([]byte(err.Error())))
		return resp, core.NewStatus(http.StatusInternalServerError)
	}
	resp.Body = io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("%v OK", resp.StatusCode))))
	return resp, core.NewStatus(resp.StatusCode)
}

func ExampleDo_Error() {
	ctrl := NewController("google-search", NewPrimaryResource("www.google.com", "", 0, "/health/liveness", httpCall), nil)

	_, status := ctrl.Do(testDo, nil)
	fmt.Printf("test: Do(testDo,nil) -> [status:%v]\n", status)

	//Output:
	//test: Do(testDo,nil) -> [status:Invalid Argument [invalid argument : request is nil]]

}

func ExampleDo_Exchange() {
	//defer DisableLogging(true)()
	//auth := "github/advanced-go/search"
	ctrl := NewExchangeController("google-search", testDo)
	uri := "http://localhost:8081/github/advanced-go/search:yahoo?" + uri2.BuildQuery("q=golang")
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set(core.XFrom, "github/advanced-go/from")

	resp, status := ctrl.Do(nil, req)
	var buf []byte
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_0s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: Do_0s() -> [status-code:200] [status:OK] [buf:true]

}

func ExampleDo_Internal() {
	defer DisableLogging(true)()
	ctrl := NewController("google-search", NewPrimaryResource("www.google.com", "", 0, "/health/liveness", httpCall), nil)
	uri := "/search?" + uri2.BuildQuery("q=golang")
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := ctrl.Do(nil, req)
	var buf []byte
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_0s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("www.google.com", "", time.Millisecond*5, "/health/liveness", httpCall), nil)
	resp, status = ctrl.Do(nil, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_5ms() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>]
	//test: Do_0s() -> [status-code:200] [status:OK] [buf:false]
	//test: httpCall() -> [content:false] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>]
	//test: Do_5ms() -> [status-code:504] [status:Timeout] [buf:false]

}

func ExampleDo_Egress() {
	var buf []byte
	ctrl := NewController("google-search", NewPrimaryResource("www.google.com", "", 0, "/health/liveness", nil), nil)
	uri := "/search?" + uri2.BuildQuery("q=golang")
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := ctrl.Do(testDo, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_0s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("www.google.com", "", time.Millisecond*5, "/health/liveness", nil), nil)
	resp, status = ctrl.Do(testDo, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_5ms() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: Do_0s() -> [status-code:200] [status:OK] [buf:true]
	//test: Do_5ms() -> [status-code:504] [status:Deadline Exceeded [context deadline exceeded]] [buf:true]

}
