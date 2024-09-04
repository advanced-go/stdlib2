package uri

import (
	"fmt"
	"strings"
	"sync"
)

var (
	localAuthority = "localhost:8080"
)

const (
	httpScheme  = "http"
	httpsScheme = "https"
	localHost   = "localhost"
)

// Attr - key, value pair
type Attr struct {
	Key, Value string
}

// SetLocalAuthority - set the local authority
func SetLocalAuthority(authority string) {
	localAuthority = authority
}

// Resolver2 - type
type Resolver2 struct {
	template *sync.Map
}

// NewResolver - create a resolver
func NewResolver2() *Resolver2 {
	return new(Resolver2)
}

// SetTemplates - configure templates
func (r *Resolver2) SetTemplates(values []Attr) func() {
	if len(values) == 0 {
		r.template = nil
		return func() {}
	}
	m := r.template
	r.template = new(sync.Map)
	for _, attr := range values {
		key, _ := TemplateToken(attr.Key)
		r.template.Store(key, attr.Value)
	}
	return func() {
		r.template = m
	}
}

// Build - perform resolution
func (r *Resolver2) Build(path string, values ...any) string {
	if len(path) == 0 {
		return "resolver error: invalid argument, path is empty"
	}
	return r.BuildWithAuthority(localAuthority, path, values...)
}

// BuildWithAuthority - perform resolution
func (r *Resolver2) BuildWithAuthority(authority, path string, values ...any) string {
	if len(path) == 0 {
		return "resolver error: invalid argument, path is empty"
	}
	if r.template != nil {
		if uri, ok := r.ExpandUrl(path); ok {
			if len(values) > 0 && strings.Index(uri, "%v") != -1 {
				uri = fmt.Sprintf(uri, values...)
			}
			return uri
		}
	}
	if !strings.HasPrefix(path, "/") {
		path += "/"
	}
	if len(values) > 0 {
		path = fmt.Sprintf(path, values...)
	}
	scheme := httpsScheme
	if len(authority) == 0 || strings.HasPrefix(authority, localHost) {
		authority = localAuthority
		scheme = httpScheme
	}
	url2 := scheme + "://" + authority + path
	return url2
}

// ExpandUrl - return the expanded URL
func (r *Resolver2) ExpandUrl(path string) (string, bool) {
	if r.template == nil {
		return "", false
	}
	if v, ok := r.template.Load(path); ok {
		if s, ok2 := v.(string); ok2 {
			return s, true
		}
	}
	return "", false
}
