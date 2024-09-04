package uri

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	testRespName2 = "file://[cwd]/uritest/get-all-resp-v1.txt"
)

func ExampleBuildRsc() {
	ver := ""
	rsc := "access"
	r := BuildRsc(ver, rsc)

	fmt.Printf("test: BuildRsc(\"%v\",\"%v\") -> [%v]\n", ver, rsc, r)

	ver = "v1"
	r = BuildRsc(ver, rsc)
	fmt.Printf("test: BuildRsc(\"%v\",\"%v\") -> [%v]\n", ver, rsc, r)

	//Output:
	//test: BuildRsc("","access") -> [access]
	//test: BuildRsc("v1","access") -> [v1/access]

}

/*
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


*/

func ExampleBuildPath2() {
	auth := "github/advanced-go/timeseries"
	vers := "v2"
	rsc := "access"
	values := make(url.Values)
	p := BuildPath2(auth, vers, rsc, values)

	fmt.Printf("test: BuildPath2(\"%v\",\"%v\",\"%v\") -> [%v]\n", auth, vers, rsc, p)

	values.Add("region", "*")
	p = BuildPath2(auth, vers, rsc, values)
	fmt.Printf("test: BuildPath2(\"%v\",\"%v\",\"%v\") -> [%v]\n", auth, vers, rsc, p)

	//Output:
	//test: BuildPath2("github/advanced-go/timeseries","v2","access") -> [github/advanced-go/timeseries:v2/access]
	//test: BuildPath2("github/advanced-go/timeseries","v2","access") -> [github/advanced-go/timeseries:v2/access?region=*]

}

func ExampleResolve2() {
	host := ""
	auth := "github/advanced-go/timeseries"
	vers := "v2"
	rsc := "access"
	values := make(url.Values)

	url := Resolve(host, auth, vers, rsc, values, nil)
	fmt.Printf("test: Resolve(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, vers, rsc, url)

	values.Add("region", "*")
	url = Resolve(host, auth, vers, rsc, values, nil)
	fmt.Printf("test: Resolve(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, vers, rsc, url)

	host = "www.google.com"
	url = Resolve(host, auth, vers, rsc, values, nil)
	fmt.Printf("test: Resolve(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, vers, rsc, url)

	host = "localhost:8080"
	url = Resolve(host, auth, vers, rsc, values, nil)
	fmt.Printf("test: Resolve(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, vers, rsc, url)

	h := make(http.Header)
	url = Resolve(host, auth, vers, rsc, values, h)
	fmt.Printf("test: Resolve(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, vers, rsc, url)

	h.Add(BuildPath2(auth, vers, rsc, values), testRespName2)
	url = Resolve(host, auth, vers, rsc, values, h)
	fmt.Printf("test: Resolve(\"%v\",\"%v\",\"%v\",\"%v\") -> [%v]\n", host, auth, vers, rsc, url)

	//Output:
	//test: Resolve("","github/advanced-go/timeseries","v2","access") -> [github/advanced-go/timeseries:v2/access]
	//test: Resolve("","github/advanced-go/timeseries","v2","access") -> [github/advanced-go/timeseries:v2/access?region=*]
	//test: Resolve("www.google.com","github/advanced-go/timeseries","v2","access") -> [https://www.google.com/github/advanced-go/timeseries:v2/access?region=*]
	//test: Resolve("localhost:8080","github/advanced-go/timeseries","v2","access") -> [http://localhost:8080/github/advanced-go/timeseries:v2/access?region=*]
	//test: Resolve("localhost:8080","github/advanced-go/timeseries","v2","access") -> [http://localhost:8080/github/advanced-go/timeseries:v2/access?region=*]
	//test: Resolve("localhost:8080","github/advanced-go/timeseries","v2","access") -> [file://[cwd]/uritest/get-all-resp-v1.txt]

}
