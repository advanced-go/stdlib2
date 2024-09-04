package controller

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"io"
	"net/http"
)

const (
	requestId = "123-request-id"
	relatesTo = "test-relates-to"
)

func Example_ResponseWriter() {
	requestId2 := "123-request-id"
	relatesTo2 := "test-relates-to"
	content := "this is response write content"
	w := newResponseWriter()

	w.Header().Add(core.XRequestId, requestId2)
	w.Header().Add(core.XRelatesTo, relatesTo2)
	w.WriteHeader(http.StatusAccepted)
	cnt, err := w.Write([]byte(content))
	fmt.Printf("test: responseWriter() -> [cnt:%v] [error:%v]\n", cnt, err)

	resp := w.Response()

	fmt.Printf("test: responseWriter() -> [write-requestId:%v] [response-requestId:%v]\n", requestId, resp.Header.Get(core.XRequestId))
	fmt.Printf("test: responseWriter() -> [write-relatesTo:%v] [response-relatesTo:%v]\n", relatesTo, resp.Header.Get(core.XRelatesTo))
	fmt.Printf("test: responseWriter() -> [write-statusCode:%v] [response-statusCode:%v]\n", http.StatusAccepted, resp.StatusCode)

	buf, _ := io.ReadAll(resp.Body)

	fmt.Printf("test: responseWriter() -> [write-content:%v] [response-content:%v]\n", content, string(buf))

	//Output:
	//test: responseWriter() -> [cnt:30] [error:<nil>]
	//test: responseWriter() -> [write-requestId:123-request-id] [response-requestId:123-request-id]
	//test: responseWriter() -> [write-relatesTo:test-relates-to] [response-relatesTo:test-relates-to]
	//test: responseWriter() -> [write-statusCode:202] [response-statusCode:202]
	//test: responseWriter() -> [write-content:this is response write content] [response-content:this is response write content]

}
