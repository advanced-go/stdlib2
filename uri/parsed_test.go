package uri

import (
	"fmt"
	"net/url"
)

func ExampleURLParse_Raw() {
	u := "http://localhost:8080/github/advanced-go/stdlib/uri.Uproot?q=golang"
	uri, err := url.Parse(u)
	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [host:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Host, uri.Path, uri.Fragment, uri.RawQuery, err)

	u = "http://localhost:8080/github/advanced-go/stdlib/uri:Uproot?q=golang"
	uri, err = url.Parse(u)
	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Path, uri.Fragment, uri.RawQuery, err)

	u = "http://localhost:8080/github/advanced-go/stdlib/uri?q=golang#Uproot"
	uri, err = url.Parse(u)
	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Path, uri.Fragment, uri.RawQuery, err)

	u = uri.Path
	uri, err = url.Parse(u)
	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Path, uri.Fragment, uri.RawQuery, err)

	//Output:
	//test: ParseRaw("http://localhost:8080/github/advanced-go/stdlib/uri.Uproot?q=golang") -> [scheme:http] [host:localhost:8080] [path:/github/advanced-go/stdlib/uri.Uproot] [frag:] [query:q=golang] [err:<nil>]
	//test: ParseRaw("http://localhost:8080/github/advanced-go/stdlib/uri:Uproot?q=golang") -> [scheme:http] [path:/github/advanced-go/stdlib/uri:Uproot] [frag:] [query:q=golang] [err:<nil>]
	//test: ParseRaw("http://localhost:8080/github/advanced-go/stdlib/uri?q=golang#Uproot") -> [scheme:http] [path:/github/advanced-go/stdlib/uri] [frag:Uproot] [query:q=golang] [err:<nil>]
	//test: ParseRaw("/github/advanced-go/stdlib/uri") -> [scheme:] [path:/github/advanced-go/stdlib/uri] [frag:] [query:] [err:<nil>]

}

func ExampleParsed() {
	uri := "http://localhost:8081/github/advanced-go/guidance:v1/resiliency?" + BuildQuery("region=region1")
	p := Uproot(uri)

	fmt.Printf("test: Uproot() -> [auth:%v] [vers:%v] [rsc:%v] [path:%v] [query:%v]\n", p.Authority, p.Version, p.Resource, p.Path, p.Query)

	//Output:
	//test: Uproot() -> [auth:github/advanced-go/guidance] [vers:v1] [rsc:resiliency] [path:resiliency] [query:region=region1]

}

func ExampleParsed_Version() {
	p := new(Parsed)

	prev := p.Path
	parseVersion(p)
	fmt.Printf("test: parseVersion(\"%v\") -> [path:%v] [vers:%v]\n", prev, p.Path, p.Version)

	p.Path = "search"
	prev = p.Path
	parseVersion(p)
	fmt.Printf("test: parseVersion(\"%v\") -> [path:%v] [vers:%v]\n", prev, p.Path, p.Version)

	p.Path = "v1/search"
	prev = p.Path
	parseVersion(p)
	fmt.Printf("test: parseVersion(\"%v\") -> [path:%v] [vers:%v]\n", prev, p.Path, p.Version)

	//Output:
	//test: parseVersion("") -> [path:] [vers:]
	//test: parseVersion("search") -> [path:search] [vers:]
	//test: parseVersion("v1/search") -> [path:search] [vers:v1]

}

func ExampleParsed_PathURL() {
	url := "https://www.google.com/github/advanced-go/search:google"
	p := Uproot(url)
	u := p.PathQuery()
	fmt.Printf("test: Parsed(\"%v\") -> [pathURL:%v] [query:%v]\n", url, u, u.Query().Encode())

	url = "https://www.google.com/github/advanced-go/search:v2/google"
	p = Uproot(url)
	u = p.PathQuery()
	fmt.Printf("test: Parsed(\"%v\") -> [pathURL:%v] [query:%v]\n", url, u, u.Query().Encode())

	url = "https://www.google.com/github/advanced-go/search:v2/google?" + BuildQuery("q=golang")
	p = Uproot(url)
	u = p.PathQuery()
	fmt.Printf("test: Parsed(\"%v\") -> [pathURL:%v] [query:%v]\n", url, u, u.Query().Encode())

	//Output:
	//test: Parsed("https://www.google.com/github/advanced-go/search:google") -> [pathURL:google] [query:]
	//test: Parsed("https://www.google.com/github/advanced-go/search:v2/google") -> [pathURL:google] [query:]
	//test: Parsed("https://www.google.com/github/advanced-go/search:v2/google?q=golang") -> [pathURL:google?q=golang] [query:q=golang]

}

func ExampleParseURL() {
	values := make(url.Values)
	values.Add("q", "*")
	uri := "http://localhost:8081/github/advanced-go/search:yahoo?" + BuildQuery(values)
	u, _ := url.Parse(uri)

	url1, parsed := ParseURL("", u)
	fmt.Printf("test: ParseURL(\"%v\") -> [url:%v] [host:%v] [path:%v] [query:%v]\n", uri, url1, parsed.Host, parsed.Path, parsed.Query)

	uri = "http://www.google.com/search/+all/usa?" + BuildQuery(values)
	u, _ = url.Parse(uri)

	url1, parsed = ParseURL("", u)
	fmt.Printf("test: ParseURL(\"%v\") -> [url:%v] [host:%v] [path:%v] [query:%v]\n", uri, url1, parsed.Host, parsed.Path, parsed.Query)

	//Output:
	//test: ParseURL("http://localhost:8081/github/advanced-go/search:yahoo?q=%2A") -> [url:http://localhost:8081/github/advanced-go/search:yahoo?q=*] [host:localhost:8081] [path:yahoo] [query:q=*]
	//test: ParseURL("http://www.google.com/search/+all/usa?q=%2A") -> [url:http://www.google.com/search/+all/usa?q=*] [host:www.google.com] [path:/search/+all/usa] [query:q=*]

}
