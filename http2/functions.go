package http2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/httpx/httpxtest"
	"github.com/advanced-go/stdlib/io"
	"github.com/advanced-go/stdlib2/core"
	"net/http"
	"testing"
)

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"
)

type Args struct {
	Item string
	Got  string
	Want string
	Err  *core.Status
}

func ReadHttp(basePath, reqName, respName string) ([]Args, *http.Request, *http.Response) {
	path := basePath + reqName
	req, status := httpxtest.NewRequest(httpxtest.ParseRaw(path))
	if !status.OK() {
		return []Args{{Item: fmt.Sprintf("ReadRequest(%v)", path), Got: "", Want: "", Err: status}}, nil, nil
	}
	path = basePath + respName
	resp, status1 := httpxtest.NewResponse(httpxtest.ParseRaw(path))
	if !status1.OK() {
		return []Args{{Item: fmt.Sprintf("ReadResponse(%v)", path), Got: "", Want: "", Err: status1}}, nil, nil
	}
	return nil, req, resp
}

func Headers(got *http.Response, want *http.Response, names ...string) (failures []Args) {
	if names == nil {
		for _, name := range want.Header {
			names = append(names, name[0])
		}
	}
	for _, name := range names {
		wantVal := want.Header.Get(name)
		if wantVal == "" {
			return []Args{{Item: name, Got: "", Want: "", Err: core.NewStatusError(core.StatusInvalidArgument, errors.New(fmt.Sprintf("want header [%v] is missing or empty", name)))}}
		}
		gotVal := got.Header.Get(name)
		if wantVal != gotVal {
			failures = append(failures, Args{Item: name, Got: gotVal, Want: wantVal, Err: nil})
		}
	}
	return failures
}

func ContentT(got *http.Response, want *http.Response) (failures []Args, gotBuf []byte, wantBuf []byte) {
	// validate content type matches
	fails, _ := validateContentType(got, want)
	if fails != nil {
		failures = fails
		return
	}
	var status *core.Status

	// validate body IO
	wantBuf, status = io.ReadAll(want.Body, nil)
	if status.Err != nil {
		failures = []Args{{Item: "want.Body", Got: "", Want: "", Err: status}}
		return
	}
	gotBuf, status = io.ReadAll(got.Body, nil)
	if status.Err != nil {
		failures = []Args{{Item: "got.Body", Got: "", Want: "", Err: status}}
	}
	return
}

// Unmarshal - unmarshal json
func Unmarshal[T any](gotBuf, wantBuf []byte) (failures []Args, gotT []T, wantT []T) {
	err := json.Unmarshal(wantBuf, &wantT)
	if err != nil {
		failures = []Args{{Item: "want.Unmarshal()", Got: "", Want: "", Err: core.NewStatusError(core.StatusJsonDecodeError, err)}}
		return
	}
	err = json.Unmarshal(gotBuf, &gotT)
	if err != nil {
		failures = []Args{{Item: "got.Unmarshal()", Got: "", Want: "", Err: core.NewStatusError(core.StatusInvalidArgument, err)}}
	}
	return
}

func Errorf(t *testing.T, failures []Args) {
	for _, arg := range failures {
		t.Errorf("%v got = %v want = %v", arg.Item, arg.Got, arg.Want)
	}
}

func validateContentType(got *http.Response, want *http.Response) (failures []Args, ct string) {
	ct = want.Header.Get(contentType)
	if ct == "" {
		return []Args{{Item: contentType, Got: "", Want: "", Err: core.NewStatusError(core.StatusInvalidArgument, errors.New("want Response header Content-Type is empty"))}}, ct
	}
	gotCt := got.Header.Get(contentType)
	if gotCt != ct {
		return []Args{{Item: contentType, Got: gotCt, Want: ct, Err: nil}}, ct
	}
	return nil, ct
}

/*
// optional
//if testBytes != nil {
//	failures = testBytes(got, gotBytes, want, wantBytes)
//	if failures != nil {
return
}
//	}


	// if no content is wanted, return
	if len(wantBytes) == 0 {
		return
	}

	// validate content length
	//if len(gotBytes) != len(wantBytes) {
	//	failures = []Args{{Item: "Content-Length", Got: fmt.Sprintf("%v", len(gotBytes)), Want: fmt.Sprintf("%v", len(wantBytes))}}
	//	return
	//}

	// validate content type is application/json
	if ct != contentTypeJson {
		failures = []Args{{Item: "Content-Type", Got: "", Want: "", Err: errors.New(fmt.Sprintf("invalid content type for serialization [%v]", ct))}}
		return
	}

	// unmarshal
	err := json.Unmarshal(wantBytes, &wantT)
	if err != nil {
		failures = []Args{{Item: "want.Unmarshal()", Got: "", Want: "", Err: err}}
		return
	}
	err = json.Unmarshal(gotBytes, &gotT)
	if err != nil {
		failures = []Args{{Item: "got.Unmarshal()", Got: "", Want: "", Err: err}}
	} else {
		content = true
	}

*/
