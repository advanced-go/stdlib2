package controller2

import (
	"fmt"
	uri2 "github.com/advanced-go/stdlib/uri"
	"io"
	"net/http"
	"time"
)

func ExampleExchange_Error() {
	//ctrl := NewController("google-search", NewPrimaryResource("www.google.com", "", 0,  httpCall), nil)
	//RegisterController(ctrl)
	_, status := Exchange(nil, nil, nil)
	fmt.Printf("test: Exchange(nil) -> [status:%v]\n", status)

	//Output:
	//test: Exchange(nil) -> [status:Invalid Argument [invalid argument : request is nil]]

}

func ExampleExchange_Internal() {
	//defer DisableLogging(true)()
	//authority := "github/advanced-go/search"
	ctrl := NewController("google-search", NewPrimaryResource("www.google.com", "", 0, httpCall), nil)
	uri := "https://www.google.com/search?" + uri2.BuildQuery("q=golang")
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := Exchange(req, testDo, ctrl)
	var buf []byte
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Exchange_0s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("www.google.com", "", time.Millisecond*5, httpCall), nil)
	resp, status = Exchange(req, testDo, ctrl)
	if status.OK() && resp.Body != nil {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Exchange_5ms() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>]
	//test: Exchange_0s() -> [status-code:200] [status:OK] [buf:false]
	//test: httpCall() -> [content:false] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>]
	//test: Exchange_5ms() -> [status-code:504] [status:Timeout] [buf:false]

}

func ExampleExchange_Egress() {
	var buf []byte
	ctrl := NewController("google-search", NewPrimaryResource("www.google.com", "", 0, nil), nil)
	uri := "https://www.google.com/search?" + uri2.BuildQuery("q=golang")
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := Exchange(req, testDo, ctrl)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Exchange_0s() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("www.google.com", "", time.Millisecond*5, nil), nil)
	resp, status = Exchange(req, testDo, ctrl)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Exchange_5ms() -> [status-code:%v] [status:%v] [buf:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: Exchange_0s() -> [status-code:200] [status:OK] [buf:true]
	//test: Exchange_5ms() -> [status-code:504] [status:Deadline Exceeded [context deadline exceeded]] [buf:true]

}
