package uri

import (
	"net/http"
	"net/url"
	"strings"
)

func Resolve(host, authority, version, resource string, values url.Values, h http.Header) string {
	path := BuildPath2(authority, version, resource, values)
	if h != nil {
		p2 := h.Get(path)
		if p2 != "" {
			return p2
		}
	}
	if host == "" {
		return path
	}
	if path[0] == '/' {
		return BuildHostWithScheme(host) + path
	}
	return BuildHostWithScheme(host) + "/" + path
}

func BuildPath2(authority, version, resource string, values url.Values) string {
	path := strings.Builder{}
	path.WriteString(authority)
	path.WriteString(":")
	path.WriteString(formatVersion(version))
	path.WriteString(resource)
	path.WriteString(formatValues(values))
	return path.String()
}

func BuildRsc(version, resource string) string {
	return formatVersion(version) + resource
}

func UNUSED_BuildHostWithScheme(host string) string {
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

/*
newUrl := strings.Builder{}
if host != "" {
scheme := HttpsScheme
if strings.Contains(host, Localhost) {
scheme = HttpScheme
}
newUrl.WriteString(scheme)
newUrl.WriteString("://")
newUrl.WriteString(host)
}
newUrl.WriteString(fmt.Sprintf(path, formatVersion(version)))
newUrl.WriteString(formatValues(values))
return newUrl.String()
newUrl := strings.Builder{}
	if host != "" {
		scheme := httpsScheme
		if strings.Contains(host, localHost) {
			scheme = httpScheme
		}
		newUrl.WriteString(scheme)
		newUrl.WriteString("://")
		newUrl.WriteString(host)
	}
	newUrl.WriteString(authority)
	newUrl.WriteString(":")
	newUrl.WriteString(path)
	newUrl.WriteString(formatValues(values))
	return newUrl.String()
*/
