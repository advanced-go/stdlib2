package controller

import (
	"fmt"
	"time"
)

func ExampleControlsLookup() {
	p := NewControls()
	//path := "http://localhost:8080/github/advanced-go/search:google?q=golang"
	auth := "github/advanced-go/search"
	ctrl := NewController("test-route", NewPrimaryResource("", auth, time.Second*2, "", httpCall), nil)

	_, status := p.lookup("")
	fmt.Printf("test: Lookup(\"\") -> [status:%v]\n", status)

	_, status = p.lookup(auth)
	fmt.Printf("test: Lookup(%v) -> [status:%v]\n", auth, status)

	err := p.registerWithAuthority(ctrl)
	fmt.Printf("test: Register() -> [err:%v]\n", err)

	handler, status1 := p.lookup(auth)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", auth, status1, handler != nil)

	host := "www.google.com"
	ctrl = NewController("test-route", NewPrimaryResource(host, "", time.Second*2, "", httpCall), nil)

	err = p.register(ctrl)
	fmt.Printf("test: Register() -> [err:%v]\n", err)
	handler, status1 = p.lookup(host)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", host, status1, handler != nil)

	//Output:
	//test: Lookup("") -> [status:Invalid Argument [invalid argument: authority is empty]]
	//test: Lookup(github/advanced-go/search) -> [status:Invalid Argument [invalid argument: Controller does not exist: [github/advanced-go/search]]]
	//test: Register() -> [err:<nil>]
	//test: Lookup(github/advanced-go/search) -> [status:OK] [handler:true]
	//test: Register() -> [err:<nil>]
	//test: Lookup(www.google.com) -> [status:OK] [handler:true]

}

func _ExampleControlsLookup_Config() {
	p := NewControls()
	//path := "http://localhost:8080/github/advanced-go/search:google?q=golang"
	auth := "github/advanced-go/search"
	ctrl := NewController("test-route", NewPrimaryResource("", auth, time.Second*2, "", httpCall), nil)

	_, status := p.lookup("")
	fmt.Printf("test: Lookup(\"\") -> [status:%v]\n", status)

	_, status = p.lookup(auth)
	fmt.Printf("test: Lookup(%v) -> [status:%v]\n", auth, status)

	err := p.registerWithAuthority(ctrl)
	fmt.Printf("test: Register() -> [err:%v]\n", err)

	handler, status1 := p.lookup(auth)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", auth, status1, handler != nil)

	host := "www.google.com"
	ctrl = NewController("test-route", NewPrimaryResource(host, "", time.Second*2, "", httpCall), nil)

	err = p.register(ctrl)
	fmt.Printf("test: Register() -> [err:%v]\n", err)
	handler, status1 = p.lookup(host)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", host, status1, handler != nil)

	//Output:
	//test: Lookup("") -> [status:Invalid Argument [invalid argument: authority is empty]]
	//test: Lookup(github/advanced-go/search) -> [status:Invalid Argument [invalid argument: Controller does not exist: [github/advanced-go/search]]]
	//test: Register() -> [err:<nil>]
	//test: Lookup(github/advanced-go/search) -> [status:OK] [handler:true]
	//test: Register() -> [err:<nil>]
	//test: Lookup(www.google.com) -> [status:OK] [handler:true]

}

func ExampleUpdatePrimaryExchange() {

}
