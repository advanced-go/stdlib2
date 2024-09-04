package core

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	RegionKey                = "region"
	ZoneKey                  = "zone"
	SubZoneKey               = "sub-zone"
	HostKey                  = "host"
	InstanceIdKey            = "id"
	RouteKey                 = "route"
	RegionZoneHostFmt        = "%v:%v.%v.%v"
	RegionZoneSubZoneHostFmt = "%v:%v.%v.%v.%v"
)

// Origin - location
type Origin struct {
	Region     string `json:"region"`
	Zone       string `json:"zone"`
	SubZone    string `json:"sub-zone"`
	Host       string `json:"host"`
	InstanceId string `json:"instance-id"`
	Route      string `json:"route"`
}

func (o Origin) Tag() string {
	tag := o.Region
	if o.Zone != "" {
		tag += ":" + o.Zone
	}
	if o.SubZone != "" {
		tag += ":" + o.SubZone
	}
	if o.Host != "" {
		tag += ":" + o.Host
	}
	return tag
}

func (o Origin) Uri(class string) string {
	var uri string
	if o.SubZone == "" {
		uri = fmt.Sprintf(RegionZoneHostFmt, class, o.Region, o.Zone, o.Host)
	} else {
		uri = fmt.Sprintf(RegionZoneSubZoneHostFmt, class, o.Region, o.Zone, o.SubZone, o.Host)
	}
	if o.Route != "" {
		uri += "." + o.Route
	}
	return uri
}

func NewValues(o Origin) url.Values {
	values := make(url.Values)
	if o.Region != "" {
		values.Add(RegionKey, o.Region)
	}
	if o.Zone != "" {
		values.Add(ZoneKey, o.Zone)
	}
	if o.SubZone != "" {
		values.Add(SubZoneKey, o.SubZone)
	}
	if o.Host != "" {
		values.Add(HostKey, o.Host)
	}
	return values
}

func NewOrigin(values url.Values) Origin {
	o := Origin{}
	if values != nil {
		o.Region = values.Get(RegionKey)
		o.Zone = values.Get(ZoneKey)
		o.SubZone = values.Get(SubZoneKey)
		o.Host = values.Get(HostKey)
	}
	return o
}

func OriginMatch(target Origin, filter Origin) bool {
	isFilter := false
	if filter.Region != "" {
		if filter.Region == "*" {
			return true
		}
		isFilter = true
		if !StringMatch(target.Region, filter.Region) {
			return false
		}
	}
	if filter.Zone != "" {
		isFilter = true
		if !StringMatch(target.Zone, filter.Zone) {
			return false
		}
	}
	if filter.SubZone != "" {
		isFilter = true
		if !StringMatch(target.SubZone, filter.SubZone) {
			return false
		}
	}
	if filter.Host != "" {
		isFilter = true
		if !StringMatch(target.Host, filter.Host) {
			return false
		}
	}
	return isFilter
}

func StringMatch(target, filter string) bool {
	//if filter == "" {
	//	return true
	//}
	return strings.ToLower(target) == strings.ToLower(filter)
}
