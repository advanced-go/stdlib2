package uri

import (
	"net/url"
	"strings"
)

const (
	HttpScheme   = "http"
	HttpsScheme  = "https"
	Localhost    = "localhost"
	Internalhost = "internalhost"
)

// BuildURL - build an url with the components provided, escaping the query
// TODO : escaping on path ?? url.PathEscape
func BuildURL(host, version, path string, query any) string {
	newUrl := strings.Builder{}
	if host != "" {
		scheme := HttpsScheme
		if host != "" {
			if strings.Contains(host, Localhost) {
				scheme = HttpScheme
			}
		}
		newUrl.WriteString(scheme)
		newUrl.WriteString("://")
		newUrl.WriteString(host)
	}
	if len(path) > 0 {
		if path[:1] != "/" {
			path += "/"
		}
	}
	if version != "" {
		newUrl.WriteString("/")
		newUrl.WriteString(version)
	}
	newUrl.WriteString(path)
	q := BuildQuery(query)
	if q != "" {
		newUrl.WriteString("?")
		newUrl.WriteString(q)
	}
	return newUrl.String()
}

// BuildURLWithAuthority - build an url with the components provided, escaping the query
// TODO : escaping on path ?? url.PathEscape
func BuildURLWithAuthority2(host, authority, version, path string, query any) string {
	newUrl := strings.Builder{}
	if host != "" {
		scheme := HttpsScheme
		if host != "" {
			if strings.Contains(host, Localhost) {
				scheme = HttpScheme
			}
		}
		newUrl.WriteString(scheme)
		newUrl.WriteString("://")
		newUrl.WriteString(host)
	}
	if authority != "" {
		newUrl.WriteString("/")
		newUrl.WriteString(authority)
	}
	isVersion := false
	newUrl.WriteString(":")
	if version != "" {
		isVersion = true
		newUrl.WriteString(version)
	}
	if len(path) > 0 {
		if !isVersion && path[:1] == "/" {
			path = path[1:]
		}
		if isVersion && path[:1] != "/" {
			path = "/" + path
		}
	}
	newUrl.WriteString(path)
	q := BuildQuery(query)
	if q != "" {
		newUrl.WriteString("?")
		newUrl.WriteString(q)
	}
	return newUrl.String()
}

// BuildQuery - build a query string with escaping
func BuildQuery(query any) string {
	if query == nil {
		return ""
	}
	if s, ok := query.(string); ok {
		return BuildValues(s).Encode()
	}
	if v, ok := query.(url.Values); ok {
		return v.Encode()
	}
	return ""
}

// BuildValues - build an url.Values type
func BuildValues(query string) url.Values {
	values := make(url.Values)
	if query == "" {
		return values
	}
	args := strings.Split(query, "&")
	if len(args) == 0 {
		return values
	}
	value := ""
	for _, v := range args {
		pair := strings.Split(v, "=")
		if len(pair) == 1 || pair[1] == "" {
			value = "invalid"
		} else {
			value = pair[1]
		}
		values.Add(pair[0], value)
	}
	return values
}

// TransformURL - build a new URL by transforming an existing URL
func TransformURL(host string, uri *url.URL) *url.URL {
	if uri == nil {
		return uri
	}
	if host == "" {
		host = uri.Host
		if host == "" {
			host = Localhost
		}
	}
	q, _ := url.QueryUnescape(uri.RawQuery)
	newURL := BuildURL(host, "", uri.Path, q)
	u, err1 := url.Parse(newURL)
	if err1 != nil {
		return uri
	}
	return u
}

/*
	scheme := HttpsScheme
	if host == "" {
		if uri.Host != "" {
			host = uri.Host
		} else {
			host = Localhost
		}
	}
	if strings.Contains(host, Localhost) {
		scheme = HttpScheme
	}
	var newUri = scheme + "://" + host

	if len(uri.Path) > 0 {
		if uri.Path[:1] != "/" {
			newUri += "/"
		}
		newUri += uri.Path
	}
	if len(uri.RawQuery) > 0 {
		newUri += "?" + uri.RawQuery
	}
	/*
		if authority == "" {
			if len(uri.Path) > 0 {
				newUri += uri.Path
			}
			if len(uri.RawQuery) > 0 {
				newUri += "?" + uri.RawQuery
			}
		} else {
			//parsed := Uproot(uri.Path)
			//newUri += "/" + authority
			//if len(parsed.Path) > 0 {
			//	newUri += ":" + parsed.Path //uri.Path[1:]
			//}
			if len(uri.Path) > 0 {
				newUri += uri.Path
			}
			if len(uri.RawQuery) > 0 {
				newUri += "?" + uri.RawQuery
			}
		}

*/
