package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	uri2 "github.com/advanced-go/stdlib/uri"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const (
	timeout = time.Second * 3
)

// Ping - function to "ping" an agent
func Ping(ex *Exchange, uri any) *core.Status {
	to, status := createTo(uri)
	if !status.OK() {
		return status
	}
	var response *Message

	result := make(chan *core.Status)
	reply := make(chan *Message, 16)
	msg := NewControlMessage(to, PkgPath, PingEvent)
	msg.ReplyTo = NewReceiverReplyTo(reply)
	err := ex.Send(msg)
	if err != nil {
		return core.NewStatusError(http.StatusInternalServerError, err)
	}
	go Receiver(timeout, reply, result, func(msg *Message) bool {
		response = msg
		return true
	})
	status = <-result
	status.AddLocation()
	if response != nil {
		status.Code = response.Status().Code
		status.Err = response.Status().Err
	}
	close(reply)
	close(result)
	return status
}

func createTo(uri any) (string, *core.Status) {
	if uri == nil {
		return "", core.NewStatusError(http.StatusBadRequest, errors.New("error: Ping() uri is nil"))
	}
	path := ""
	if u, ok := uri.(*url.URL); ok {
		path = u.Path
	} else {
		if u2, ok1 := uri.(string); ok1 {
			path = u2
		} else {
			return "", core.NewStatusError(http.StatusBadRequest, errors.New(fmt.Sprintf("error: Ping() uri is invalid type: %v", reflect.TypeOf(uri).String())))
		}
	}
	p := uri2.Uproot(path)
	if !p.Valid {
		return "", core.NewStatusError(http.StatusBadRequest, errors.New(fmt.Sprintf("error: Ping() uri is not a valid URN %v", path)))
	}
	return p.Authority, core.StatusOK()
}
