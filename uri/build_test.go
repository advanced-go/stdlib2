package uri

import (
	"fmt"
	//"github.com/advanced-go/stdlib/core"
	"net/http"
	url2 "net/url"
)

func ExampleBuildQuery() {
	s := ""
	q := BuildQuery("")
	u, _ := url2.QueryUnescape(q)
	fmt.Printf("test: BuildQuery(\"%v\") -> [query:%v] [unesc:%v]\n", s, q, u)

	s = "region=*&zone=texas"
	q = BuildQuery(s)
	u, _ = url2.QueryUnescape(q)
	fmt.Printf("test: BuildQuery(\"%v\") -> [query:%v] [unesc:%v]\n", s, q, u)

	v := make(url2.Values)
	v.Add("region", "*")
	v.Add("zone", "texas")
	q = BuildQuery(v)
	u, _ = url2.QueryUnescape(q)
	fmt.Printf("test: BuildQuery(\"%v\") -> [query:%v] [unesc:%v]\n", s, q, u)

	//Output:
	//test: BuildQuery("") -> [query:] [unesc:]
	//test: BuildQuery("region=*&zone=texas") -> [query:region=%2A&zone=texas] [unesc:region=*&zone=texas]
	//test: BuildQuery("region=*&zone=texas") -> [query:region=%2A&zone=texas] [unesc:region=*&zone=texas]

}

func ExampleBuildValues() {
	q := ""
	values := BuildValues(q)
	fmt.Printf("test: BuildValues(\"%v\") -> [values:%v]\n", q, values)

	q = "regions=*&zone=&sub-zone=dallas"
	values = BuildValues(q)
	fmt.Printf("test: BuildValues(\"%v\") -> [values:%v]\n", q, values)

	q = "regions=*&zone=texas&sub-zone"
	values = BuildValues(q)
	fmt.Printf("test: BuildValues(\"%v\") -> [values:%v]\n", q, values)

	q = "regions=*&zone=texas&sub-zone=dallas"
	values = BuildValues(q)
	fmt.Printf("test: BuildValues(\"%v\") -> [values:%v]\n", q, values)

	//Output:
	//test: BuildValues("") -> [values:map[]]
	//test: BuildValues("regions=*&zone=&sub-zone=dallas") -> [values:map[regions:[*] sub-zone:[dallas] zone:[invalid]]]
	//test: BuildValues("regions=*&zone=texas&sub-zone") -> [values:map[regions:[*] sub-zone:[invalid] zone:[texas]]]
	//test: BuildValues("regions=*&zone=texas&sub-zone=dallas") -> [values:map[regions:[*] sub-zone:[dallas] zone:[texas]]]

}

func ExampleBuildURL() {
	host := ""
	version := ""
	path := "/search/yahoo"
	query := "q=golang&region=*"
	u := BuildURL(host, version, path, query)

	u1, err := url2.Parse(u)
	fmt.Printf("test: BuildURL(\"%v\",\"%v\",\"%v\",\"%v\") -> [uri:%v] [url:%v] [err:%v]\n", host, version, path, query, u, u1, err)

	host = "www.google.com"
	version = "v1"
	values := make(url2.Values)
	values.Add("q", "golang")
	values.Add("region", "*")
	u = BuildURL(host, version, path, values)
	u1, err = url2.Parse(u)
	fmt.Printf("test: BuildURL(\"%v\",\"%v\",\"%v\",\"%v\") -> [uri:%v] [url:%v] [err:%v]\n", host, version, path, values, u, u1, err)

	//Output:
	//test: BuildURL("","","/search/yahoo","q=golang&region=*") -> [uri:/search/yahoo?q=golang&region=%2A] [url:/search/yahoo?q=golang&region=%2A] [err:<nil>]
	//test: BuildURL("www.google.com","v1","/search/yahoo","map[q:[golang] region:[*]]") -> [uri:https://www.google.com/v1/search/yahoo?q=golang&region=%2A] [url:https://www.google.com/v1/search/yahoo?q=golang&region=%2A] [err:<nil>]

}

func ExampleBuildURL_WithAuthority() {
	host := ""
	version := ""
	authority := ""
	path := "/search/yahoo"
	query := "q=golang&region=*"
	u := BuildURLWithAuthority2(host, authority, version, path, query)

	//u1, err := url2.Parse(u)
	fmt.Printf("test: BuildURLWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\",\"%v\") -> [uri:%v]\n", host, version, authority, path, query, u)

	host = "www.google.com"
	version = "v1"
	authority = "github/advanced-go/stdlib"
	values := BuildValues("q=golang&region=*")
	u = BuildURLWithAuthority2(host, authority, version, path, values)
	//u1, err = url2.Parse(u)
	fmt.Printf("test: BuildURLWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\",\"%v\") -> [uri:%v]\n", host, version, authority, path, values, u)

	//Output:
	//test: BuildURLWithAuthority("","","","/search/yahoo","q=golang&region=*") -> [uri::search/yahoo?q=golang&region=%2A]
	//test: BuildURLWithAuthority("www.google.com","v1","github/advanced-go/stdlib","/search/yahoo","map[q:[golang] region:[*]]") -> [uri:https://www.google.com/github/advanced-go/stdlib:v1/search/yahoo?q=golang&region=%2A]

}

func ExampleTransformURL() {
	host := "www.google.com"
	authority := "github/advanced-go/search"
	uri := "http://localhost:8081/github/advanced-go/search:google?" + BuildQuery("q=golang&region=*")

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	url := TransformURL(host, req.URL)
	fmt.Printf("test: TransformURL(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, host, authority, url)

	//Output:
	//test: TransformURL("http://localhost:8081/github/advanced-go/search:google?q=golang&region=%2A") [host:www.google.com] [auth:github/advanced-go/search] [url:https://www.google.com/github/advanced-go/search:google?q=golang&region=%2A]

}

func ExampleTransformURL_Host() {
	uri := "/search?" + BuildQuery("q=golang&region=*")
	host := ""

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	url := TransformURL(host, req.URL)
	fmt.Printf("test: TransformURL(\"%v\") [host:%v] [url:%v]\n", uri, host, url)

	host = "localhost:8080"
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = TransformURL(host, req.URL)
	fmt.Printf("test: TransformURL(\"%v\") [host:%v] [url:%v]\n", uri, host, url)

	uri = "/update"
	host = "www.google.com"
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = TransformURL(host, req.URL)
	fmt.Printf("test: TransformURL(\"%v\") [host:%v] [url:%v]\n", uri, host, url)

	//Output:
	//test: TransformURL("/search?q=golang&region=%2A") [host:] [url:http://localhost/search?q=golang&region=%2A]
	//test: TransformURL("/search?q=golang&region=%2A") [host:localhost:8080] [url:http://localhost:8080/search?q=golang&region=%2A]
	//test: TransformURL("/update") [host:www.google.com] [url:https://www.google.com/update]

}

func _ExampleTransformURL_Authority() {
	uri := "/google?q=golang"
	host := ""
	authority := "github/advanced-go/search"

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	url := TransformURL(host, req.URL)
	fmt.Printf("test: BuildUri(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, host, authority, url)

	host = "www.google.com"
	req, _ = http.NewRequest(http.MethodGet, uri, nil)
	url = TransformURL(host, req.URL)
	fmt.Printf("test: BuildUri(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, host, authority, url)

	/*
		// localhost
		rsc = NewPrimaryResource("localhost:8080", "", 0, "/health/liveness", httpCall)
		req, _ = http.NewRequest(http.MethodGet, uri, nil)
		url = rsc.BuildUri(req.URL)
		fmt.Printf("test: BuildUri(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url)

		// non-localhost
		uri = "/update"
		rsc = NewPrimaryResource("www.google.com", "", 0, "/health/liveness", httpCall)
		req, _ = http.NewRequest(http.MethodGet, uri, nil)
		url = rsc.BuildUri(req.URL)
		fmt.Printf("test: BuildUri(\"%v\") [host:%v] [auth:%v] [url:%v]\n", uri, rsc.Host, rsc.Authority, url

	*/

	//Output:
	//test: BuildUri("/google?q=golang") [host:] [auth:github/advanced-go/search] [url:http://localhost/github/advanced-go/search:google?q=golang]
	//test: BuildUri("/google?q=golang") [host:www.google.com] [auth:github/advanced-go/search] [url:https://www.google.com/github/advanced-go/search:google?q=golang]

}
