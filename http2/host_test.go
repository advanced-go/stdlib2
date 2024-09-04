package http2

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/httpx"
	"net/http"
)

const (
	rsc1Name = "rsc1-name"
	rsc2Name = "rsc-2-name"
	url1     = "http://localhost:8080/github/advanced-go/documents:resiliency"
	url2     = "http://localhost:8080/github/advanced-go/guidance:resiliency"
)

func rsc1(r *http.Request) (*http.Response, *core.Status) {
	if r.Method == http.MethodGet && r.URL.Path == core.AuthorityRootPath {
		return httpx.NewAuthorityResponse(rsc1Name), core.StatusOK()
	}
	return &http.Response{StatusCode: http.StatusBadRequest}, core.StatusOK()
}

func rsc2(r *http.Request) (*http.Response, *core.Status) {
	if r.Method == http.MethodGet && r.URL.Path == core.AuthorityRootPath {
		return httpx.NewAuthorityResponse(rsc2Name), core.StatusOK()
	}
	return &http.Response{StatusCode: http.StatusGatewayTimeout}, core.StatusOK()
}

func ExampleNewHost_Error() {
	_, err := NewHost("", nil, nil)
	fmt.Printf("test: NewHost() -> [err:%v]\n", err)

	_, err = NewHost("github/advanced-go/stdlib", nil, nil)
	fmt.Printf("test: NewHost() -> [err:%v]\n", err)

	_, err = NewHost("github/advanced-go/stdlib", mapper, nil)
	fmt.Printf("test: NewHost() -> [err:%v]\n", err)

	_, err = NewHost("github/advanced-go/stdlib", mapper, rsc1, rsc1)
	fmt.Printf("test: NewHost() -> [err:%v]\n", err)

	//Output:
	//test: NewHost() -> [err:error: authority is empty]
	//test: NewHost() -> [err:resource map function is nil]
	//test: NewHost() -> [err:error: invalid resource map, resource name is empty]
	//test: NewHost() -> [err:error: invalid resource name, Exchange already exists for: rsc1-name]

}

/*
func ExampleNewHost_Error() {
	h := NewHost("github/advanced-go/stdlib", nil, nil)
	req, _ := http.NewRequest(http.MethodPut, "https://www.google.com/search?Q=golang", nil)

	resp := h.Do(nil)
	buf, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: NewHost() -> [status-code:%v] [content:%v]\n", resp.StatusCode, string(buf))

	resp = h.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: NewHost() -> [status-code:%v] [content:%v]\n", resp.StatusCode, string(buf))

	h = NewHost("github/advanced-go/stdlib", func(req *http.Request) string { return "invalid" }, nil)
	resp = h.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	fmt.Printf("test: NewHost() -> [status-code:%v] [content:%v]\n", resp.StatusCode, string(buf))

	//Output:
	//test: NewHost() -> [status-code:400] [content:bad request: http.Request is nil]
	//test: NewHost() -> [status-code:400] [content:invalid resource map, resource name is empty for: [https://www.google.com/search?Q=golang]]
	//test: NewHost() -> [status-code:400] [content:invalid resource map, HttpExchange not found for: [https://www.google.com/search?Q=golang]]

}

*/

func mapper(r *http.Request) string {
	if r.URL.String() == url1 {
		return rsc1Name
	}
	return rsc2Name
}

func ExampleNewHost_Do() {
	h, _ := NewHost("github/advanced-go/stdlib", mapper, rsc1, rsc2)

	req, _ := http.NewRequest(http.MethodGet, url1, nil)
	resp, _ := h.Do(req)
	fmt.Printf("test: Do(\"%v\") -> [map:%v] [status-code:%v]\n", url1, mapper(req), resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, url2, nil)
	resp, _ = h.Do(req)
	fmt.Printf("test: Do(\"%v\") -> [map:%v] [status-code:%v]\n", url2, mapper(req), resp.StatusCode)

	//Output:
	//test: Do("http://localhost:8080/github/advanced-go/documents:resiliency") -> [map:rsc1-name] [status-code:400]
	//test: Do("http://localhost:8080/github/advanced-go/guidance:resiliency") -> [map:rsc-2-name] [status-code:504]

}
