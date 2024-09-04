package http2

import (
	"fmt"
	"github.com/advanced-go/stdlib2/core"
	"net/http"
	"time"
)

type postContent2 struct {
	Item string
}

var (
	origin2 = []core.Origin{
		{Region: "region1", Zone: "Zone1", Host: "www.host1.com"},
		{Region: "region1", Zone: "Zone2", Host: "www.host2.com"},
		{Region: "region2", Zone: "Zone1", Host: "www.google.com"},
	}
)

func match2(req *http.Request, item *core.Origin) bool {
	filter := core.NewOrigin(req.URL.Query())
	if core.OriginMatch(*item, filter) {
		return true
	}
	return false
}

func patch2(_ *http.Request, list *[]core.Origin, patch *Patch) *core.Status {
	count := 0
	for _, op := range patch.Updates {
		switch op.Op {
		case OpReplace:
			if op.Path == core.HostKey {
				count++
				if s, ok1 := op.Value.(string); ok1 {
					(*list)[0].Host = s
				}
			}
		default:
		}
	}
	if count == 0 {
		return core.StatusNotFound()
	}
	return core.StatusOK()
}

func post2(_ *http.Request, list *[]core.Origin, _ *postContent2) *core.Status {
	(*list)[0].Host = "www.facebook.com"
	return core.StatusOK()
}

var (
	content = NewListContent[core.Origin, Patch, postContent2](true, match2, patch2, post2)
)

func _ExampleListContent_Count2() {
	fmt.Printf("example-start()  ->       %v\n", time.Now().UTC())
	go func() {
		time.Sleep(time.Millisecond * 500)
		//content.count2("goid-2")
	}()
	//content.count2("goid-1")
	time.Sleep(time.Second * 3)
	fmt.Printf("example-stop()   ->       %v\n", time.Now().UTC())

	//Output:
	//fail

}

func ExampleListContent_Put() {
	req, _ := http.NewRequest(http.MethodPut, "https://ww.google.com/search?q=golang", nil)
	status := content.Put(req, origin2)
	fmt.Printf("test: Put() -> [status:%v]\n", status)

	//Output:
	//test: Put() -> [status:OK]

}

func ExampleListContent_Get() {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=zOne1", nil)
	items, status := content.Get(req)
	fmt.Printf("test: Get() -> [status:%v] [count:%v] [items:%v]\n", status, len(items), items)

	//Output:
	//test: Get() -> [status:OK] [count:2] [items:[{region1 Zone1  www.host1.com } {region2 Zone1  www.google.com }]]

}

func ExampleListContent_GetEmpty() {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", nil)
	items, status := content.Get(req)
	fmt.Printf("test: Get() -> [status:%v] [count:%v] [items:%v]\n", status, len(items), items)

	//Output:
	//test: Get() -> [status:Not Found] [count:0] [items:[]]

}

func getList() []core.Origin {
	var list []core.Origin

	if l, ok := any(content).(*ListContent[core.Origin, Patch, postContent2]); ok {
		list = l.List
	}
	return list
}

func ExampleListContent_Delete() {
	// Empty
	req, _ := http.NewRequest(http.MethodDelete, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", nil)
	status := content.Delete(req)
	fmt.Printf("test: Delete-0() -> [status:%v]\n", status)

	// Delete 1
	prevCount := content.Count()
	req, _ = http.NewRequest(http.MethodDelete, "https://localhost:8081/github/advanced-go/documents:resiliency?host=www.host2.com", nil)
	status = content.Delete(req)
	fmt.Printf("test: Delete-1() -> [status:%v] [prev:%v] [curr:%v] %v\n", status, prevCount, content.Count(), getList())

	// Delete remaining 2
	prevCount = content.Count()
	req, _ = http.NewRequest(http.MethodDelete, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=zone1", nil)
	status = content.Delete(req)
	fmt.Printf("test: Delete-2() -> [status:%v] [prev:%v] [curr:%v] %v\n", status, prevCount, content.Count(), getList())

	req, _ = http.NewRequest(http.MethodPut, "https://ww.google.com/search?q=golang", nil)
	status = content.Put(req, origin2)
	fmt.Printf("test: Put() -> [status:%v] [count:%v]\n", status, content.Count())

	//Output:
	//test: Delete-0() -> [status:Not Found]
	//test: Delete-1() -> [status:OK] [prev:3] [curr:2] [{region1 Zone1  www.host1.com } {region2 Zone1  www.google.com }]
	//test: Delete-2() -> [status:OK] [prev:2] [curr:0] []
	//test: Put() -> [status:OK] [count:3]

}

func ExampleListContent_Patch() {
	p := Patch{Updates: []Operation{
		{Op: OpReplace, Path: core.HostKey, Value: "www.search.yahoo.com"},
	}}
	fmt.Printf("test: Patch-before() -> %v\n", getList())

	req, _ := http.NewRequest(http.MethodPatch, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", nil)
	status := content.Patch(req, &p)
	fmt.Printf("test: Patch-after() -> [status:%v] %v\n", status, getList())

	//Output:
	//test: Patch-before() -> [{region1 Zone1  www.host1.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]
	//test: Patch-after() -> [status:OK] [{region1 Zone1  www.search.yahoo.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]

}

func ExampleListContent_Post() {
	p := postContent2{}

	fmt.Printf("test: Post-before() -> %v\n", getList())

	req, _ := http.NewRequest(http.MethodPost, "https://localhost:8081/github/advanced-go/documents:resiliency?zone=invalid", nil)
	status := content.Post(req, &p)
	fmt.Printf("test: Post-after() -> [status:%v] %v\n", status, getList())

	//Output:
	//test: Post-before() -> [{region1 Zone1  www.search.yahoo.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]
	//test: Post-after() -> [status:OK] [{region1 Zone1  www.facebook.com } {region1 Zone2  www.host2.com } {region2 Zone1  www.google.com }]

}
