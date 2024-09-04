package http2

import (
	"fmt"
	json2 "github.com/advanced-go/stdlib/json"
	"github.com/advanced-go/stdlib2/core"
	"net/http"
)

func finalize2(resp *http.Response) {
	if resp.Header == nil {
		resp.Header = make(http.Header)
	}
	resp.Header.Add(core.XAuthority, "github/advanced-go/stdlib")
	if resp.Request != nil {
		resp.Header.Add("x-method", resp.Request.Method)
	}
}

var (
	content3 = NewListContent[core.Origin, Patch, postContent2](false, match2, patch2, post2)
	rsc3     = NewResource[core.Origin, Patch, postContent2]("rsc-origin", content3, finalize2)
)

func getList3() []core.Origin {
	var list []core.Origin

	if l, ok := any(content3).(*ListContent[core.Origin, Patch, postContent2]); ok {
		list = l.List
	}
	return list
}

func ExampleResource_Put() {
	reader, _, status := json2.NewReadCloser(origin2)
	if !status.OK() {
		fmt.Printf("test: Put() -> [read-closer-status:%v]\n", status)
	} else {
		req, _ := http.NewRequest(http.MethodPut, "https://localhost:8081/github/advanced-go/documents:resiliency", reader)
		resp, status1 := rsc3.Do(req)
		fmt.Printf("test: Put() -> [status:%v] [status-code:%v] [header:%v] [%v]\n", status1, resp.StatusCode, resp.Header, getList3())
	}

	//Output:
	//test: Put() -> [status:OK] [status-code:200] [header:map[X-Authority:[github/advanced-go/stdlib] X-Method:[PUT]]] [[{region1 Zone1  www.host1.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]]

}

func ExampleResource_Get() {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=zone1", nil)
	resp, status := rsc3.Do(req)
	if !status.OK() {
		fmt.Printf("test: Do() -> [status:%v]\n", status)
	} else {
		items, status1 := json2.New[[]core.Origin](resp.Body, nil)
		fmt.Printf("test: Get() -> [status:%v] [status-code:%v] [header:%v] [%v]\n", status1, resp.StatusCode, resp.Header, items)
	}

	//Output:
	//test: Get() -> [status:OK] [status-code:200] [header:map[Content-Type:[application/json] X-Authority:[github/advanced-go/stdlib] X-Method:[GET]]] [[{region1 Zone1  www.host1.com } {region2 Zone1  www.google.com }]]

}

func ExampleResource_Delete() {
	// Empty
	req, _ := http.NewRequest(http.MethodDelete, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", nil)
	resp, status := rsc3.Do(req)
	fmt.Printf("test: Do() -> [status:%v] [status-code:%v]\n", status, resp.StatusCode)

	// Delete 1
	prevCount := content3.Count()
	req, _ = http.NewRequest(http.MethodDelete, "https://localhost:8081/github/advanced-go/documents:resiliency?host=www.host2.com", nil)
	resp, status = rsc3.Do(req)
	fmt.Printf("test: Do() -> [status:%v] [prev:%v] [curr:%v] %v\n", status, prevCount, content3.Count(), getList3())

	// Delete remaining 2
	prevCount = content3.Count()
	req, _ = http.NewRequest(http.MethodDelete, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=zone1", nil)
	resp, status = rsc3.Do(req)
	fmt.Printf("test: Do() -> [status:%v] [prev:%v] [curr:%v] %v\n", status, prevCount, content3.Count(), getList3())

	// Re-initialize
	reader, _, status1 := json2.NewReadCloser(origin2)
	if !status1.OK() {
		fmt.Printf("test: Put() -> [read-closer-status:%v]\n", status1)
	} else {
		req, _ = http.NewRequest(http.MethodPut, "https://localhost:8081/github/advanced-go/documents:resiliency", reader)
		resp, status1 = rsc3.Do(req)
		fmt.Printf("test: Put() -> [status:%v] [status-code:%v] [header:%v] [%v]\n", status1, resp.StatusCode, resp.Header, getList3())
	}

	//Output:
	//test: Do() -> [status:Not Found] [status-code:404]
	//test: Do() -> [status:OK] [prev:3] [curr:2] [{region1 Zone1  www.host1.com } {region2 Zone1  www.google.com }]
	//test: Do() -> [status:OK] [prev:2] [curr:0] []
	//test: Put() -> [status:OK] [status-code:200] [header:map[X-Authority:[github/advanced-go/stdlib] X-Method:[PUT]]] [[{region1 Zone1  www.host1.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]]

}

func ExampleResource_Patch() {
	p := Patch{Updates: []Operation{
		{Op: OpReplace, Path: core.HostKey, Value: "www.search.yahoo.com"},
	}}
	fmt.Printf("test: Patch-before() -> %v\n", getList3())

	rc, _, status := json2.NewReadCloser(p)
	if !status.OK() {
		fmt.Printf("test: NewReaderCloser() -> [status:%v]\n", status)
	} else {
		req, _ := http.NewRequest(http.MethodPatch, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", rc)
		resp, status1 := rsc3.Do(req)
		fmt.Printf("test: Patch-after() -> [status:%v] [status-code:%v] %v\n", status1, resp.StatusCode, getList3())
	}

	//Output:
	//test: Patch-before() -> [{region1 Zone1  www.host1.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]
	//test: Patch-after() -> [status:OK] [status-code:200] [{region1 Zone1  www.search.yahoo.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]

}

func ExampleResource_Post() {
	p := postContent2{}
	fmt.Printf("test: Post-before() -> %v\n", getList3())

	rc, _, status := json2.NewReadCloser(p)
	if !status.OK() {
		fmt.Printf("test: NewReaderCloser() -> [status:%v]\n", status)
	} else {
		req, _ := http.NewRequest(http.MethodPost, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", rc)
		resp, status1 := rsc3.Do(req)
		fmt.Printf("test: Post-after() -> [status:%v] [status-code:%v] %v\n", status1, resp.StatusCode, getList3())
	}

	//Output:
	//test: Post-before() -> [{region1 Zone1  www.search.yahoo.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]
	//test: Post-after() -> [status:OK] [status-code:200] [{region1 Zone1  www.facebook.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]

}
