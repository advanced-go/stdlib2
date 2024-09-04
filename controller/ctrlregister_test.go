package controller

import "fmt"

const (
	auth1 = "localhost:8087"
	auth2 = "github/advanced-go/test-controller"
	auth3 = "www.google3.com"
)

var testCtrls = []*Controller{
	NewController("test-route-1", NewPrimaryResource(auth1, "", 0, "", nil), nil),
	NewController("test-route-2", NewPrimaryResource("", auth2, 0, "", nil), nil),
	NewController("test-route-2", NewPrimaryResource("", auth3, 0, "", nil), nil),
}

func ExampleRegisterController_Defer() {
	authority1 := "localhost:9999"
	authority2 := "github/advanced-go/test-guidance"
	ctrl := NewController("test-route-1", NewPrimaryResource(authority1, "", 0, "", nil), nil)
	ctrl2 := NewController("test-route-2", NewPrimaryResource("", authority2, 0, "", nil), nil)

	fn := RegisterControllerWithDefer(ctrl, nil)
	fn = RegisterControllerWithDefer(ctrl2, fn)

	_, status := LookupWithAuthority(authority1)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", authority1, status)

	_, status = LookupWithAuthority(authority2)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", authority2, status)

	fn()
	_, status = LookupWithAuthority(authority1)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", authority1, status)

	_, status = LookupWithAuthority(authority2)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", authority2, status)

	//Output:
	//test: LookupWithAuthority("localhost:9999") -> [status:OK]
	//test: LookupWithAuthority("github/advanced-go/test-guidance") -> [status:OK]
	//test: LookupWithAuthority("localhost:9999") -> [status:Invalid Argument [invalid argument: Controller does not exist: [localhost:9999]]]
	//test: LookupWithAuthority("github/advanced-go/test-guidance") -> [status:Invalid Argument [invalid argument: Controller does not exist: [github/advanced-go/test-guidance]]]

}

func ExampleRegisterControllerList_Defer() {
	fn := RegisterControllerListWithDefer(testCtrls)
	fmt.Printf("test: RegisterControllerListWithDefer() -> [fn:%v]\n", fn != nil)

	_, status := LookupWithAuthority(auth1)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth1, status)
	_, status = LookupWithAuthority(auth2)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth2, status)
	_, status = LookupWithAuthority(auth3)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth3, status)

	fn() //UnregisterControllerList(testCtrls)

	_, status = LookupWithAuthority(auth1)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth1, status)
	_, status = LookupWithAuthority(auth2)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth2, status)
	_, status = LookupWithAuthority(auth3)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth3, status)

	//Output:
	//test: RegisterControllerListWithDefer() -> [fn:true]
	//test: LookupWithAuthority("localhost:8087") -> [status:OK]
	//test: LookupWithAuthority("github/advanced-go/test-controller") -> [status:OK]
	//test: LookupWithAuthority("www.google3.com") -> [status:OK]
	//test: LookupWithAuthority("localhost:8087") -> [status:Invalid Argument [invalid argument: Controller does not exist: [localhost:8087]]]
	//test: LookupWithAuthority("github/advanced-go/test-controller") -> [status:Invalid Argument [invalid argument: Controller does not exist: [github/advanced-go/test-controller]]]
	//test: LookupWithAuthority("www.google3.com") -> [status:Invalid Argument [invalid argument: Controller does not exist: [www.google3.com]]]

}

func ExampleRegisterControllerList() {
	err := RegisterControllerList(testCtrls)
	fmt.Printf("test: RegisterControllerList() -> [err:%v]\n", err)

	_, status := LookupWithAuthority(auth1)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth1, status)
	_, status = LookupWithAuthority(auth2)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth2, status)
	_, status = LookupWithAuthority(auth3)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth3, status)

	UnregisterControllerList(testCtrls)

	_, status = LookupWithAuthority(auth1)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth1, status)
	_, status = LookupWithAuthority(auth2)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth2, status)
	_, status = LookupWithAuthority(auth3)
	fmt.Printf("test: LookupWithAuthority(\"%v\") -> [status:%v]\n", auth3, status)

	//Output:
	//test: RegisterControllerList() -> [err:<nil>]
	//test: LookupWithAuthority("localhost:8087") -> [status:OK]
	//test: LookupWithAuthority("github/advanced-go/test-controller") -> [status:OK]
	//test: LookupWithAuthority("www.google3.com") -> [status:OK]
	//test: LookupWithAuthority("localhost:8087") -> [status:Invalid Argument [invalid argument: Controller does not exist: [localhost:8087]]]
	//test: LookupWithAuthority("github/advanced-go/test-controller") -> [status:Invalid Argument [invalid argument: Controller does not exist: [github/advanced-go/test-controller]]]
	//test: LookupWithAuthority("www.google3.com") -> [status:Invalid Argument [invalid argument: Controller does not exist: [www.google3.com]]]

}

func ExampleRegisterController_Error() {
	err := RegisterController(nil)
	fmt.Printf("test: RegisterController(nil) -> [err:%v]\n", err)

	ctrl := NewController("test-route", nil, nil)
	ctrl.Router = nil
	err = RegisterController(ctrl)
	fmt.Printf("test: RegisterController(ctrl) -> [err:%v]\n", err)

	ctrl = NewController("test-route", nil, nil)
	err = RegisterController(ctrl)
	fmt.Printf("test: RegisterController(ctrl) -> [err:%v]\n", err)

	ctrl = NewController("test-route", NewPrimaryResource("", "", 0, "", nil), nil)
	err = RegisterController(ctrl)
	fmt.Printf("test: RegisterController(ctrl) -> [err:%v]\n", err)

	//Output:
	//test: RegisterController(nil) -> [err:invalid argument: Controller is nil]
	//test: RegisterController(ctrl) -> [err:invalid argument: Controller router is nil]
	//test: RegisterController(ctrl) -> [err:invalid argument: Controller router primary resource is nil]
	//test: RegisterController(ctrl) -> [err:invalid argument: Controller router primary resource host is empty]

}
