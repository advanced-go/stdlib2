package uri

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	testRespName = "file://[cwd]/timeseries1test/get-all-resp-v1.txt"
	defaultKey   = "default"
	googleKey    = "google"
	yahooKey     = "yahoo"
	bingKey      = "bing"
)

func ExampleBuildHostWithScheme() {
	host := ""
	o := BuildHostWithScheme(host)
	fmt.Printf("test: BuildHostWithScheme(\"%v\") -> [origin:%v]\n", host, o)

	host = "www.google.com"
	o = BuildHostWithScheme(host)
	fmt.Printf("test: BuildHostWithScheme(\"%v\") -> [origin:%v]\n", host, o)

	host = "localhost:8080"
	o = BuildHostWithScheme(host)
	fmt.Printf("test: BuildHostWithScheme(\"%v\") -> [origin:%v]\n", host, o)

	host = "internalhost"
	o = BuildHostWithScheme(host)
	fmt.Printf("test: BuildHostWithScheme(\"%v\") -> [origin:%v]\n", host, o)

	//Output:
	//test: BuildHostWithScheme("") -> [origin:]
	//test: BuildHostWithScheme("www.google.com") -> [origin:https://www.google.com]
	//test: BuildHostWithScheme("localhost:8080") -> [origin:http://localhost:8080]
	//test: BuildHostWithScheme("internalhost") -> [origin:http://internalhost]

}

func ExampleBuildPath() {
	auth := "github/advanced-go/timeseries"
	rsc := "access"
	ver := "v2"
	values := make(url.Values)

	p := BuildPath(rsc, values)
	fmt.Printf("test: BuildPath(\"%v\") -> [%v]\n", rsc, p)

	p = BuildPathWithAuthority(auth, "", rsc, values)
	fmt.Printf("test: BuildPathWithAuthority(\"%v\",\"%v\",\"%v\") -> [%v]\n", auth, ver, rsc, p)

	values.Add("region", "*")
	p = BuildPath(rsc, values)
	fmt.Printf("test: BuildPath(\"%v\") -> [%v]\n", rsc, p)

	p = BuildPathWithAuthority(auth, ver, rsc, values)
	fmt.Printf("test: BuildPathWithAuthority(\"%v\",\"%v\",\"%v\") -> [%v]\n", auth, ver, rsc, p)

	//Output:
	//test: BuildPath("access") -> [access]
	//test: BuildPathWithAuthority("github/advanced-go/timeseries","v2","access") -> [github/advanced-go/timeseries:access]
	//test: BuildPath("access") -> [access?region=*]
	//test: BuildPathWithAuthority("github/advanced-go/timeseries","v2","access") -> [github/advanced-go/timeseries:v2/access?region=*]

}

func ExampleResolve_Url() {
	errType := 123
	host := ""
	path := "/search"
	values := make(url.Values)
	r := NewResolver("localhost:8081")

	url1 := r.Url(host, path, errType, nil)
	fmt.Printf("test: Url(\"%v\",\"%v\") -> [%v]\n", host, path, url1)

	url1 = r.Url(host, path, values, nil)
	fmt.Printf("test: Url(\"%v\",\"%v\") -> [%v]\n", host, path, url1)

	values.Add("q", "golang")
	url1 = r.Url(host, path, values, nil)
	fmt.Printf("test: Url(\"%v\",\"%v\") -> [%v]\n", host, path, url1)

	url1 = r.Url(host, path, "q=golang", nil)
	fmt.Printf("test: Url_String(\"%v\",\"%v\") -> [%v]\n", host, path, url1)

	host = "www.google.com"
	url1 = r.Url(host, path, values, nil)
	fmt.Printf("test: Url_String(\"%v\",\"%v\") -> [%v]\n", host, path, url1)

	h := make(http.Header)
	h.Add(BuildPath(path, values), "https://www.search.yahoo.com?q=golang")
	host = "www.google.com"
	url1 = r.Url(host, path, values, h)
	fmt.Printf("test: Url_Override(\"%v\",\"%v\") -> [%v]\n", host, path, url1)

	//Output:
	//test: Url("","/search") -> [http://localhost:8081/searcherror: query type is invalid int]
	//test: Url("","/search") -> [http://localhost:8081/search]
	//test: Url("","/search") -> [http://localhost:8081/search?q=golang]
	//test: Url_String("","/search") -> [http://localhost:8081/search?q=golang]
	//test: Url_String("www.google.com","/search") -> [https://www.google.com/search?q=golang]
	//test: Url_Override("www.google.com","/search") -> [https://www.search.yahoo.com?q=golang]

}

func ExampleResolve_UrlWithAuthority() {
	host := ""
	auth := "github/advanced-go/timeseries"
	rsc := "access"
	ver := ""
	values := make(url.Values)
	r := NewResolver("")

	url1 := r.UrlWithAuthority(host, auth, "", rsc, values, nil)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	values.Add("region", "*")
	url1 = r.UrlWithAuthority(host, auth, "", rsc, values, nil)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	url1 = r.UrlWithAuthority(host, auth, "", rsc, "region=*", nil)
	fmt.Printf("test: UrlWithAuthority_String(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	host = "www.google.com"
	url1 = r.UrlWithAuthority(host, auth, "", rsc, values, nil)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	host = "localhost:8080"
	ver = "v2"
	//rsc = "v2/" + rsc
	url1 = r.UrlWithAuthority(host, auth, ver, rsc, values, nil)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	h := make(http.Header)
	url1 = r.UrlWithAuthority(host, auth, ver, rsc, values, h)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	h.Add(BuildPathWithAuthority(auth, ver, rsc, values), testRespName)
	url1 = r.UrlWithAuthority(host, auth, ver, rsc, values, h)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	host = "www.google.com"
	rsc = "search"
	values.Del("region")
	values.Add("q", "golang")
	auth = ""
	url1 = r.UrlWithAuthority(host, auth, ver, rsc, values, nil)
	fmt.Printf("test: UrlWithAuthority(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, ver, rsc, url1)

	//Output:
	//test: UrlWithAuthority("","github/advanced-go/timeseries","","access") -> [/github/advanced-go/timeseries:access]
	//test: UrlWithAuthority("","github/advanced-go/timeseries","","access") -> [/github/advanced-go/timeseries:access?region=*]
	//test: UrlWithAuthority_String("","github/advanced-go/timeseries","","access") -> [/github/advanced-go/timeseries:access?region=*]
	//test: UrlWithAuthority("www.google.com","github/advanced-go/timeseries","","access") -> [https://www.google.com/github/advanced-go/timeseries:access?region=*]
	//test: UrlWithAuthority("localhost:8080","github/advanced-go/timeseries","v2","access") -> [http://localhost:8080/github/advanced-go/timeseries:v2/access?region=*]
	//test: UrlWithAuthority("localhost:8080","github/advanced-go/timeseries","v2","access") -> [http://localhost:8080/github/advanced-go/timeseries:v2/access?region=*]
	//test: UrlWithAuthority("localhost:8080","github/advanced-go/timeseries","v2","access") -> [file://[cwd]/timeseries1test/get-all-resp-v1.txt]
	//test: UrlWithAuthority("www.google.com","","v2","search") -> [https://www.google.com/search?q=golang]

}
