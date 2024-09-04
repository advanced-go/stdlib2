package core

import (
	"fmt"
	"net/http"
)

/*
func ExampleVersionContent() {
	s := "1.2.34"
	fmt.Printf("test: VersionContent() -> [%v]\n", VersionContent(s))

	//Output:
	//test: VersionContent() -> [{ "version": "1.2.34" }]

}

func ExampleHealthContent() {
	s := "jacked up!!"
	fmt.Printf("test: HealthContent() -> [%v]\n", HealthContent(s))

	//Output:
	//test: HealthContent() -> [{ "status": "jacked up!!" }]

}


*/

func ExampleHttpHandler() {
	ok := exchange(func(w http.ResponseWriter, r *http.Request) {})
	fmt.Printf("test: HttpHandler(anonymous-function) -> [ok:%v|\n", ok)

	ok = exchange(handler2)
	fmt.Printf("test: HttpHandler(function) -> [ok:%v|\n", ok)

	ok = exchange(handler3())
	fmt.Printf("test: HttpHandler(return-function) -> [ok:%v|\n", ok)

	//Output:
	//test: HttpHandler(anonymous-function) -> [ok:true|
	//test: HttpHandler(function) -> [ok:true|
	//test: HttpHandler(return-function) -> [ok:true|

}

func exchange(fn HttpHandler) bool {
	if fn == nil {
		return false
	}
	return true
}

func handler2(w http.ResponseWriter, r *http.Request) {
}

func handler3() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func authExchange(req *http.Request) (*http.Response, *Status) {
	if req.URL.Path == AuthorityRootPath {
		h := make(http.Header)
		h.Add(XAuthority, "github/advanced-go/stdlib")
		return &http.Response{StatusCode: http.StatusOK, Header: h}, StatusOK()
	}
	return &http.Response{StatusCode: http.StatusBadRequest}, NewStatus(http.StatusBadRequest)
}

func ExampleAuthority() {
	auth := Authority(authExchange)
	fmt.Printf("test: Authority() -> [auth:%v]\n", auth)

	//Output:
	//test: Authority() -> [auth:github/advanced-go/stdlib]

}
