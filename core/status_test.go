package core

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func ExampleNewStatus_OK() {
	s := StatusOK()

	path := reflect.TypeOf(Status{}).PkgPath()
	path += "/" + reflect.TypeOf(Status{}).Name()
	fmt.Printf("test: NewStatus() -> [status:%v] [type:%v]\n", s, path)

	s = NewStatusError(http.StatusBadGateway, errors.New("this is an error message"))
	str := defaultFormatter(testTS, s.Code, HttpStatus(s.Code), "1234-56-789", []error{s.Err}, s.Trace())

	fmt.Printf("test: NewStatus() -> %v\n", str)

	//Output:
	//test: NewStatus() -> [status:OK] [type:github.com/advanced-go/stdlib/core/Status]
	//test: NewStatus() -> { "timestamp":"2024-03-01T18:23:50.205Z", "code":502, "status":"error: code not mapped: 502", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/stdlib/tree/main/core#ExampleNewStatus_OK" ] }

}

func ExampleNewStatus_Teapot() {
	s := NewStatus(http.StatusTeapot)
	fmt.Printf("test: NewStatus() -> [status:%v]\n", s)

	s = NewStatusError(http.StatusTeapot, errors.New("this is an error message"))
	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(testTS, s.Code, HttpStatus(s.Code), "1234-56-789", []error{s.Err}, s.Trace()))

	//Output:
	//test: NewStatus() -> [status:I'm A Teapot]
	//test: NewStatus() -> { "timestamp":"2024-03-01T18:23:50.205Z", "code":418, "status":"I'm A Teapot", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/stdlib/tree/main/core#ExampleNewStatus_Teapot" ] }

}

func ExampleNewStatus_Location() {
	s := errorFunc()
	s.AddLocation()

	str := formatter(testTS, s.Code, HttpStatus(s.Code), "1234-5678", []error{s.Err}, s.Trace())
	fmt.Printf("test: Location() -> [out:%v] [trace:%v]\n", str, s.Trace())

	//Output:
	//test: Location() -> [out:{ "timestamp":"2024-03-01T18:23:50.205Z", "code":400, "status":"Bad Request", "request-id":"1234-5678", "errors" : [ "test bad request error" ], "trace" : [ "https://github.com/advanced-go/stdlib/tree/main/core#ExampleNewStatus_Location","https://github.com/advanced-go/stdlib/tree/main/core#errorFunc" ] }
	//] [trace:[github/advanced-go/stdlib/core:errorFunc github/advanced-go/stdlib/core:ExampleNewStatus_Location]]

}

func errorFunc() *Status {
	return NewStatusError(http.StatusBadRequest, errors.New("test bad request error"))
}

func ExampleNewStatus_GenericLocation() {
	s := genericErrorFunc[Output]()
	s.AddLocation()

	str := formatter(testTS, s.Code, HttpStatus(s.Code), "1234-5678", []error{s.Err}, s.Trace())
	fmt.Printf("test: GenericLocation() -> [out:%v] [trace:%v]\n", str, s.Trace())

	//Output:
	//test: GenericLocation() -> [out:{ "timestamp":"2024-03-01T18:23:50.205Z", "code":400, "status":"Bad Request", "request-id":"1234-5678", "errors" : [ "test bad request error" ], "trace" : [ "https://github.com/advanced-go/stdlib/tree/main/core#ExampleNewStatus_GenericLocation","https://github.com/advanced-go/stdlib/tree/main/core#genericErrorFunc[...]" ] }
	//] [trace:[github/advanced-go/stdlib/core:genericErrorFunc[...] github/advanced-go/stdlib/core:ExampleNewStatus_GenericLocation]]

}

func genericErrorFunc[E ErrorHandler]() *Status {
	s := NewStatusError(http.StatusBadRequest, errors.New("test bad request error"))
	return s
}

/*
func ExampleNewStatus_TeapotHandled() {
	var e Output
	s := NewStatus(http.StatusTeapot)

	//fmt.Printf("test: NewStatus() -> [status:%v]\n", s)

	s.Error = errors.New("this is an error message")
	s.AddLocation("github/advanced-go/stdlib/core:AddLocation")
	s.AddLocation("github/advanced-go/stdlib/core:TopOfList")

	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(s.Code, []error{s.Error}, s.Trace(), "1234-56-789"))
    //e.Handle()
	//Output:
	//test: NewStatus() -> [status:I'm A Teapot]
	//test: NewStatus() -> { "code":418, "status":"I'm A Teapot", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/stdlib/tree/main/core#TopOfList","https://github.com/advanced-go/stdlib/tree/main/core#AddLocation" ] }

}


*/
