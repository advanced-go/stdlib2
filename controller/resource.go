package controller

import (
	"github.com/advanced-go/stdlib/core"
	uri2 "github.com/advanced-go/stdlib/uri"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Resource struct {
	Name         string `json:"name"`
	Host         string `json:"host"`
	Authority    string `json:"authority"`
	LivenessPath string `json:"liveness"`
	Duration     time.Duration
	Handler      core.HttpExchange
}

func newResource(name, host, authority string, duration time.Duration, livenessPath string, handler core.HttpExchange) *Resource {
	r := new(Resource)
	r.Name = name
	r.Host = host
	r.Authority = authority
	r.LivenessPath = livenessPath
	r.Duration = duration
	if handler != nil {
		r.Handler = handler
	}
	return r
}

func NewPrimaryResource(host, authority string, duration time.Duration, livenessPath string, handler core.HttpExchange) *Resource {
	return newResource(PrimaryName, host, authority, duration, livenessPath, handler)
}

func NewSecondaryResource(host, authority string, duration time.Duration, livenessPath string, handler core.HttpExchange) *Resource {
	return newResource(SecondaryName, host, authority, duration, livenessPath, handler)
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
	/*
		// Authority and path
		if r.Authority != "" {
			if r.Authority[0] != '/' {
				newUrl.WriteString("/")
			}
			newUrl.WriteString(r.Authority)
			newUrl.WriteString(":")

			if uri.Path[0] == '/' {
				newUrl.WriteString(uri.Path[1:])
			} else {
				newUrl.WriteString(uri.Path)
			}
		} else {
			if uri.Path[0] != '/' {
				newUrl.WriteString("/")
			}
			newUrl.WriteString(uri.Path)
		}

	*/
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
	//return uri2.TransformURL(r.Host, uri)
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

/*
	if uri == nil {
		return uri
	}
	scheme := "https"
	host := r.Host
	if host == "" {
		host = localhost
	}
	if strings.Contains(host, localhost) {
		scheme = "http"
	}
	var newUri = scheme + "://" + host
	if r.Authority == "" {
		if len(uri.Path) > 0 {
			newUri += uri.Path
		}
		if len(uri.RawQuery) > 0 {
			newUri += "?" + uri.RawQuery
		}
	} else {
		newUri += "/" + r.Authority
		if len(uri.Path) > 0 {
			newUri += ":" + uri.Path[1:]
		}
		if len(uri.RawQuery) > 0 {
			newUri += "?" + uri.RawQuery
		}
		/*
			uri2, err := url.Parse(r.Authority)
			if err != nil {
				return uri
			}
			newUri = uri2.Scheme + "://"
			if len(uri2.Host) > 0 {
				newUri += uri2.Host
			} else {
				newUri += uri.Host
			}
			if len(uri2.Path) > 0 {
				newUri += uri2.Path
			} else {
				newUri += uri.Path
			}
			if len(uri2.RawQuery) > 0 {
				newUri += "?" + uri2.RawQuery
			} else {
				if len(uri.RawQuery) > 0 {
					newUri += "?" + uri.RawQuery
				}
			}

}
u, err1 := url.Parse(newUri)
if err1 != nil {
return uri
}
return u
*/
