package controller

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"io"
	"net/http"
	"time"
)

func ExampleDoEgress() {
	uri := "https://www.google.com/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := doEgress(time.Second*5, testDo, req)
	buf, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ExampleDoEgress_OK -> [status-code:%v] [status:%v] [content:%v]\n", resp.StatusCode, status, len(buf) > 0)

	resp, status = doEgress(time.Millisecond*500, func(r *http.Request) (*http.Response, *core.Status) {
		time.Sleep(time.Second * 2)
		return testDo(r)
	}, req)
	time.Sleep(time.Second * 3)
	fmt.Printf("test: ExampleDoEgress_Recover -> [status-code:%v] [status:%v] [content:%v]\n", resp.StatusCode, status, false)

	//Output:
	//test: ExampleDoEgress_OK -> [status-code:200] [status:OK] [content:true]
	//test: ExampleDoEgress_Recover -> [status-code:504] [status:Deadline Exceeded [context deadline exceeded]] [content:false]

}
