package http2

import "fmt"

func Example_ReadHttp() {
	base := "file://[cwd]/resource/"
	in := "put-req.txt"
	out := "put-resp.txt"
	args, req, resp := ReadHttp(base, in, out)

	fmt.Printf("test: ReadHttp(%v,%v) [err:%v] [req:%v], [resp:%v]\n", in, out, args, req != nil, resp != nil)

	//Output:
	//test: ReadHttp(put-req.txt,put-resp.txt) [err:[]] [req:true], [resp:true]

}
