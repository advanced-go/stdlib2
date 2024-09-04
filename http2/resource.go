package http2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/httpx"
	"github.com/advanced-go/stdlib/json"
	"net/http"
)

type FinalizeFunc func(*http.Response)

type Resource[T any, U any, V any] struct {
	Name             string
	Identity         *http.Response
	MethodNotAllowed *http.Response
	Finalize         FinalizeFunc
	Content          Content2[T, U, V]
}

func NewResource[T any, U any, V any](name string, content Content2[T, U, V], finalize FinalizeFunc) *Resource[T, U, V] {
	r := new(Resource[T, U, V])
	r.Identity = httpx.NewAuthorityResponse(name)
	h2 := make(http.Header)
	h2.Add(httpx.ContentType, httpx.ContentTypeText)
	r.MethodNotAllowed, _ = httpx.NewResponse[core.Log](core.NewStatus(http.StatusMethodNotAllowed).HttpCode(), h2, nil)
	r.Finalize = finalize
	if r.Finalize == nil {
		r.Finalize = defaultFinalize()
	}
	r.Content = content
	return r
}

func (r *Resource[T, U, V]) Count() int {
	return r.Content.Count()
}

func (r *Resource[T, U, V]) Empty() {
	r.Content.Empty()
}

func (r *Resource[T, U, V]) finalize(req *http.Request, status *core.Status) (*http.Response, *core.Status) {
	h2 := make(http.Header)
	if !status.OK() && status.Err != nil {
		h2.Add(httpx.ContentType, httpx.ContentTypeText)
	}
	resp, status1 := httpx.NewResponse[core.Log](status.HttpCode(), h2, status.Err)
	resp.Request = req
	r.Finalize(resp)
	return resp, status1
}

func (r *Resource[T, U, V]) Do(req *http.Request) (*http.Response, *core.Status) {
	switch req.Method {
	case http.MethodGet:
		if req.URL.Path == core.AuthorityRootPath {
			return r.Identity, core.StatusOK()
		}
		items, status := r.Content.Get(req)
		if !status.OK() {
			return r.finalize(req, status)
		}
		reader, bytes, status1 := json.NewReadCloser(items)
		if !status1.OK() {
			return r.finalize(req, status)
		}
		h2 := make(http.Header)
		h2.Add(httpx.ContentType, httpx.ContentTypeJson)
		resp := &http.Response{StatusCode: status1.HttpCode(), Status: status1.String(), Header: h2, ContentLength: bytes, Body: reader}
		resp.Request = req
		r.Finalize(resp)
		return resp, status1
	case http.MethodPut:
		items, status := json.New[[]T](req.Body, req.Header)
		if !status.OK() {
			return r.finalize(req, status)
		}
		return r.finalize(req, r.Content.Put(req, items))
	case http.MethodPatch:
		patch, status := json.New[U](req.Body, req.Header)
		if !status.OK() {
			return r.finalize(req, status)
		}
		return r.finalize(req, r.Content.Patch(req, &patch))
	case http.MethodPost:
		post, status := json.New[V](req.Body, req.Header)
		if !status.OK() {
			return r.finalize(req, status)
		}
		return r.finalize(req, r.Content.Post(req, &post))
	case http.MethodDelete:
		return r.finalize(req, r.Content.Delete(req))
	default:
		status := core.NewStatusError(http.StatusMethodNotAllowed, errors.New(fmt.Sprintf("unsupported method: %v", req.Method)))
		h2 := make(http.Header)
		h2.Add(httpx.ContentType, httpx.ContentTypeText)
		resp, status1 := httpx.NewResponse[core.Log](status.HttpCode(), h2, status.Err)
		return resp, status1
	}
}

func defaultFinalize() func(resp *http.Response) {
	return func(resp *http.Response) {
		if resp.Header == nil {
			resp.Header = make(http.Header)
			if resp.Request != nil {
				resp.Header.Add("X-Method", resp.Request.Method)
			}
		}
	}
}
