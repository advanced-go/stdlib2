package http2

/*
type FinalizeFunc func(*http.Response)

type MatchFunc[T any] func(item *T, r *http.Request) bool
type PatchProcessFunc[T any, U any] func(list *[]T, content *U) *http.Response
type PostProcessFunc[T any, V any] func(list *[]T, content *V) *http.Response

type GetFunc[T any] func(r *http.Request, list []T, match MatchFunc[T], finalize FinalizeFunc) *http.Response
type DeleteFunc[T any] func(r *http.Request, list *[]T, match MatchFunc[T], finalize FinalizeFunc) *http.Response
type PutFunc[T any] func(r *http.Request, list *[]T, finalize FinalizeFunc) *http.Response
type PatchFunc[T any, U any] func(r *http.Request, list *[]T, patch PatchProcessFunc[T, U], finalize FinalizeFunc) *http.Response
type PostFunc[T any, V any] func(r *http.Request, list *[]T, post PostProcessFunc[T, V], finalize FinalizeFunc) *http.Response

func FinalizeResponse(status *core.Status, r *http.Request, finalize FinalizeFunc) *http.Response {
	resp := NewResponse(status, status.Err)
	resp.Request = r
	if finalize != nil {
		finalize(resp)
	}
	return resp
}


*/
