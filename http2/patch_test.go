package http2

import "fmt"

func ExamplePatch() {
	patch := Patch{Updates: []Operation{
		{Op: OpAdd},
	},
	}

	fmt.Printf("test: Patch() -> [patch:%v]\n", patch)

	//Output:
	//test: Patch() -> [patch:{[{add  <nil>}]}]

}
