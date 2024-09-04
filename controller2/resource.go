package controller2

import (
	"github.com/advanced-go/stdlib/core"
	uri2 "github.com/advanced-go/stdlib/uri"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	PrimaryName   = "primary"
	SecondaryName = "secondary"
)

type Resource struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Authority string `json:"authority"`
	Duration  time.Duration
	Handler   core.HttpExchange
}

func newResource(name, host, authority string, duration time.Duration, handler core.HttpExchange) *Resource {
	r := new(Resource)
	r.Name = name
	r.Host = host
	r.Authority = authority
	r.Duration = duration
	r.Handler = handler
	return r
}

func NewPrimaryResource(host, authority string, duration time.Duration, handler core.HttpExchange) *Resource {
	return newResource(PrimaryName, host, authority, duration, handler)
}

func NewSecondaryResource(host, authority string, duration time.Duration, handler core.HttpExchange) *Resource {
	return newResource(SecondaryName, host, authority, duration, handler)
}

func (r *Resource) IsPrimary() bool {
	return r != nil && r.Name == PrimaryName
}

func (r *Resource) BuildURL(uri *url.URL) *url.URL {
	newUrl := strings.Builder{}
	// Scheme and host
	if r.Host != "" {
		scheme := uri2.HttpsScheme
		if strings.Contains(r.Host, uri2.Localhost) {
			scheme = uri2.HttpScheme
		}
		newUrl.WriteString(scheme)
		newUrl.WriteString("://")
		newUrl.WriteString(r.Host)
	} else {
		newUrl.WriteString(uri2.HttpScheme)
		newUrl.WriteString("://")
		newUrl.WriteString(uri2.Internalhost)
	}
	if uri.Path[0] != '/' {
		newUrl.WriteString("/")
	}
	newUrl.WriteString(uri.Path)

	// Query
	if uri.RawQuery != "" {
		newUrl.WriteString("?")
		newUrl.WriteString(uri.Query().Encode())
	}
	u, _ := url.Parse(newUrl.String())
	return u
}

func (r *Resource) timeout(req *http.Request) time.Duration {
	duration := r.Duration
	if r.Duration < 0 {
		duration = 0
	}
	if req == nil || req.Context() == nil {
		return duration
	}
	ct, ok := req.Context().Deadline()
	if !ok {
		return duration
	}
	until := time.Until(ct)
	if until <= duration || duration == 0 {
		return until * -1
	}
	return duration
}
