package controller2

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ExampleResource_BuildURL() {
	uri := "/search?q=golang&region=*"

	// No host
	rsc := NewPrimaryResource("", "", 0, nil)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	url := rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [url:%v]\n", uri, rsc.Host, url)

	// localhost
	uri = "search?q=golang&region=*"
	rsc = NewPrimaryResource("localhost:8080", "", 0, nil)
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [url:%v]\n", uri, rsc.Host, url)

	// non-localhost
	uri = "/update/resource"
	rsc = NewPrimaryResource("www.google.com", "", 0, nil)
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [url:%v]\n", uri, rsc.Host, url)

	//Output:
	//test: BuildURL("/search?q=golang&region=*") [host:] [url:http://internalhost/search?q=golang&region=%2A]
	//test: BuildURL("search?q=golang&region=*") [host:localhost:8080] [url:http://localhost:8080/search?q=golang&region=%2A]
	//test: BuildURL("/update/resource") [host:www.google.com] [url:https://www.google.com/update/resource]

}

func _ExampleResource_BuildURL() {
	uri := "/search?q=golang&region=*"

	// No host, no authority
	rsc := NewPrimaryResource("", "", 0, nil)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	url := rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url)

	// No host, with authority
	uri = "/yahoo?q=golang&region=*"
	rsc = NewPrimaryResource("", "github/advanced-go/search", 0, nil)
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url)

	// localhost
	uri = "/search?q=golang&region=*"
	rsc = NewPrimaryResource("localhost:8080", "", 0, nil)
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url)

	// non-localhost
	uri = "/update/resource"
	rsc = NewPrimaryResource("www.google.com", "", 0, nil)
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url)

	// authority
	uri = "/update/storage?q=golang&region=*"
	rsc = NewPrimaryResource("www.google.com", "github/advanced-go/search", 0, nil)
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = rsc.BuildURL(req.URL)
	fmt.Printf("test: BuildURL(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url)

	//Output:
	//test: BuildURL("/search?q=golang&region=*") [host:] [auth:] [url:http://internalhost/search?q=golang&region=%2A]
	//test: BuildURL("/yahoo?q=golang&region=*") [host:] [auth:github/advanced-go/search] [url:http://internalhost/github/advanced-go/search:yahoo?q=golang&region=%2A]
	//test: BuildURL("/search?q=golang&region=*") [host:localhost:8080] [auth:] [url:http://localhost:8080/search?q=golang&region=%2A]
	//test: BuildURL("/update/resource") [host:www.google.com] [auth:] [url:https://www.google.com/update/resource]
	//test: BuildURL("/update/storage?q=golang&region=*") [host:www.google.com] [auth:github/advanced-go/search] [url:https://www.google.com/github/advanced-go/search:update/storage?q=golang&region=%2A]

}

func ExampleTimeout() {
	var dIn time.Duration = -1
	uri := "/search?q=golang"
	rsc := NewPrimaryResource("localhost:8080", "", dIn, httpCall)

	dOut := rsc.timeout(nil)
	fmt.Printf("test: timeout(nil) -> [timeout:%v]\n", dOut)

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [timeout:%v]\n", dIn, dOut)

	dIn = time.Millisecond * 100
	rsc = NewPrimaryResource("localhost:8080", "", dIn, httpCall)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [timeout:%v]\n", dIn, dOut)

	//Output:
	//test: timeout(nil) -> [timeout:0s]
	//test: timeout(req) -> [duration:-1ns] [timeout:0s]
	//test: timeout(req) -> [duration:100ms] [timeout:100ms]

}

func ExampleTimeout_Deadline() {
	dIn := time.Millisecond * 200
	deadline := time.Millisecond * 100
	uri := "/search?q=golang"
	rsc := NewPrimaryResource("localhost:8080", "", dIn, httpCall)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	dOut := rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [deadline:%v] [timeout:%v]\n", dIn, deadline, dOut)

	dIn = time.Millisecond * 100
	rsc = NewPrimaryResource("localhost:8080", "", dIn, httpCall)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [deadline:%v] [timeout:%v]\n", dIn, deadline, dOut)

	dIn = time.Millisecond * 50
	rsc = NewPrimaryResource("localhost:8080", "", dIn, httpCall)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [deadline:%v] [timeout:%v]\n", dIn, deadline, dOut)

	//Output:
	//test: timeout(req) -> [duration:200ms] [deadline:100ms] [timeout:-100ms]
	//test: timeout(req) -> [duration:100ms] [deadline:100ms] [timeout:-100ms]
	//test: timeout(req) -> [duration:50ms] [deadline:100ms] [timeout:50ms]

}
