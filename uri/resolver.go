package uri

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Resolver struct {
	defaultHost string
}

func NewResolver(defaultHost string) *Resolver {
	r := new(Resolver)
	r.defaultHost = defaultHost
	return r
}

func (r *Resolver) Url(host, path string, query any, h http.Header) string {
	path1 := BuildPath(path, query)
	if h != nil {
		p2 := h.Get(path1)
		if p2 != "" {
			return p2
		}
	}
	if host != "" {
		return Cat(host, path1)
	}
	return Cat(r.defaultHost, path1)
}

func (r *Resolver) UrlWithAuthority(host, authority, version, resource string, query any, h http.Header) string {
	path := BuildPathWithAuthority(authority, version, resource, query)
	if h != nil {
		p2 := h.Get(path)
		if p2 != "" {
			return p2
		}
	}
	if host != "" {
		return Cat(host, path)
	}
	return Cat(r.defaultHost, path)
}

func Cat(host, path string) string {
	origin := BuildHostWithScheme(host)
	if path[0] == '/' {
		return origin + path
	}
	return origin + "/" + path
}

func BuildPath(path string, query any) string {
	return BuildPathWithAuthority("", "", path, query)
}

func BuildPathWithAuthority(authority, version, resource string, query any) string {
	path := strings.Builder{}
	if authority != "" {
		path.WriteString(authority)
		path.WriteString(":")
		path.WriteString(formatVersion2(version))
	}
	path.WriteString(resource)
	path.WriteString(formatQuery(query))
	return path.String()
}

func BuildHostWithScheme(host string) string {
	if host == "" {
		return ""
	}
	origin := strings.Builder{}
	scheme := HttpsScheme
	if strings.Contains(host, Localhost) || strings.Contains(host, Internalhost) {
		scheme = HttpScheme
	}
	origin.WriteString(scheme)
	origin.WriteString("://")
	origin.WriteString(host)
	return origin.String()
}

func formatQuery(query any) string {
	if query == nil {
		return ""
	}
	if v, ok := query.(url.Values); ok {
		encoded := v.Encode()
		if encoded != "" {
			encoded, _ = url.QueryUnescape(encoded)
			return "?" + encoded
		}
		return ""
	}
	if s, ok := query.(string); ok {
		return "?" + s
	}
	return fmt.Sprintf("error: query type is invalid %v", reflect.TypeOf(query))
}

func formatVersion2(version string) string {
	if version == "" {
		return ""
	}
	return version + "/"
}
