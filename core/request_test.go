package core

import (
	"fmt"
	"net/http"
)

func ExampleAddRequestId() {
	h := AddRequestId(nil)
	fmt.Printf("test: AddRequestId(nil) -> [empty:%v]\n", len(h.Get(XRequestId)) == 0)

	head := make(http.Header)
	h = AddRequestId(head)
	fmt.Printf("test: AddRequestId(head) -> [empty:%v]\n", len(h.Get(XRequestId)) == 0)

	head = make(http.Header)
	head.Add(XRequestId, "123-45-head")
	h = AddRequestId(head)
	fmt.Printf("test: AddRequestId(head) -> %v\n", h.Get(XRequestId))

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	h = AddRequestId(req)
	fmt.Printf("test: RequestId(request) -> [empty:%v]\n", len(h.Get(XRequestId)) == 0)

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	req.Header.Set(XRequestId, "123-456-request")
	h = AddRequestId(req)
	fmt.Printf("test: RequestId(request) -> %v\n", h.Get(XRequestId))

	//Output:
	//test: AddRequestId(nil) -> [empty:false]
	//test: AddRequestId(head) -> [empty:false]
	//test: AddRequestId(head) -> 123-45-head
	//test: RequestId(request) -> [empty:false]
	//test: RequestId(request) -> 123-456-request

}

func ExampleRequestId() {
	id := RequestId("123-456-string")
	fmt.Printf("test: RequestId(string) -> %v\n", id)

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	req.Header.Set(XRequestId, "123-456-request")
	id = RequestId(req)
	fmt.Printf("test: RequestId(request) -> %v\n", id)

	h := make(http.Header)
	h.Set(XRequestId, "123-456-header")
	id = RequestId(h)
	fmt.Printf("test: RequestId(header) -> %v\n", id)

	//Output:
	//test: RequestId(string) -> 123-456-string
	//test: RequestId(request) -> 123-456-request
	//test: RequestId(header) -> 123-456-header

}
