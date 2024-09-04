package uri

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	UrnScheme     = "urn"
	UrnSeparator  = ":"
	VersionPrefix = "v"
)

// Uproot - uproot an embedded uri in a URI or a URI path
func Uproot(in string) *Parsed {
	if in == "" {
		return &Parsed{Valid: false, Err: errors.New("error: invalid input, URI is empty")}
	}
	in = strings.ToLower(in)
	if strings.HasPrefix(in, UrnScheme) {
		return &Parsed{Valid: true, Authority: in, Path: in}
	}
	u, err := url.Parse(in)
	if err != nil {
		return &Parsed{Valid: false, Err: err}
	}
	var str []string
	lower := strings.ToLower(u.Path)
	if lower[0] == '/' {
		str = strings.Split(lower[1:], UrnSeparator)
	} else {
		str = strings.Split(lower, UrnSeparator)
	}
	switch len(str) {
	case 0:
		return &Parsed{Valid: false, Err: errors.New(fmt.Sprintf("error: path has no URN separator [%v]", u.Path))}
	case 1:
		return &Parsed{Valid: true, Authority: str[0], Query: u.RawQuery}
	case 2:
		p := &Parsed{Valid: true, Authority: str[0], Path: str[1], Query: u.RawQuery}
		parseVersion(p)
		index := strings.Index(p.Path, "/")
		if index != -1 {
			p.Resource = p.Path[:index]
		} else {
			p.Resource = p.Path
		}
		return p
	default:
		return &Parsed{Valid: false, Err: errors.New(fmt.Sprintf("error: path has multiple URN separators [%v]", u.Path))}
	}
}

func UprootAuthority(url *url.URL) string {
	if url == nil {
		return ""
	}
	str := strings.Split(url.Path, ":")
	if len(str) != 2 {
		return ""
	}
	if str[0][0] == '/' {
		return str[0][1:]
	}
	return str[0]
}

func parseVersion(p *Parsed) {
	if p == nil {
		return
	}
	if strings.HasPrefix(p.Path, VersionPrefix) {
		i := strings.Index(p.Path, "/")
		if i != -1 {
			p.Version = p.Path[:i]
			p.Path = p.Path[i+1:]
		}
	}
}
